package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"hw-3/aggregator"
	"hw-3/proxy"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	queryTimeout = 30 * time.Second
)

type feed struct {
	Tokens  string `yaml:"tokens"`
	Address string `yaml:"address"`
}

type feedConfig struct {
	Feeds []feed `yaml:"feeds"`
}

func parseFeedConfig(configFileName string) (*feedConfig, error) {
	feedConfigFile, err := os.Open(configFileName)
	if err != nil {
		return nil, fmt.Errorf("failed opening config file %s: %w", configFileName, err)
	}
	defer feedConfigFile.Close()

	feedConf := feedConfig{}
	feedConfigDecoder := yaml.NewDecoder(feedConfigFile)
	err = feedConfigDecoder.Decode(&feedConf)
	if err != nil {
		return nil, fmt.Errorf("failed decoding feed info from file %s: %w", configFileName, err)
	}
	return &feedConf, nil
}

func subscribeBlocks(ctx context.Context, client *ethclient.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Errorf("Failed subscribing to block creation: %s", err)
		return
	}

	log.Info("Monitoring blocks")
	for {
		select {
		case err := <-sub.Err():
			log.Errorf("Failed while listnening for new blocks: %s", err)
			return
		case <-ctx.Done():
			return
		case header := <-headers:
			timeoutCtx, cancel := context.WithTimeout(ctx, queryTimeout)
			block, err := client.BlockByHash(timeoutCtx, header.Hash())
			cancel()
			if err != nil {
				log.Errorf("Failed getting block by hash: %s", err)
				return
			}
			log.WithFields(log.Fields{
				"number":       block.Number().Uint64(),
				"transactions": len(block.Transactions()),
			}).Info("New block")
		}
	}
}

func logFeedError(feedHexAddress, tokens, msg string, err error) {
	log.WithFields(log.Fields{
		"feedAddress": feedHexAddress,
		"tokens":      tokens,
	}).Errorf("%s: %s", msg, err)
}

func subscribeEvents(ctx context.Context, client *ethclient.Client, feedHexAddress, tokens string, wg *sync.WaitGroup) {
	defer wg.Done()

	proxyAddress := common.HexToAddress(feedHexAddress)
	proxyInstance, err := proxy.NewProxy(proxyAddress, client)
	if err != nil {
		logFeedError(feedHexAddress, tokens, "Failed acquiring proxy instance", err)
		return
	}

	aggregatorAddress, err := proxyInstance.Aggregator(nil)
	if err != nil {
		logFeedError(feedHexAddress, tokens, "Failed acquiring aggregator address", err)
		return
	}

	aggregatorInstance, err := aggregator.NewAggregator(aggregatorAddress, client)
	if err != nil {
		logFeedError(feedHexAddress, tokens, "Failed acquiring aggregator instance", err)
		return
	}

	aggregatorDecimals, err := aggregatorInstance.Decimals(nil)
	if err != nil {
		logFeedError(feedHexAddress, tokens, "Failed acquiring decimals", err)
		return
	}

	decimalsDivInt := big.NewInt(10)
	decimalsDivInt.Exp(decimalsDivInt, big.NewInt(int64(aggregatorDecimals)), nil)
	decimalsDivFloat := new(big.Float).SetInt(decimalsDivInt)

	ansChan := make(chan *aggregator.AggregatorAnswerUpdated)
	sub, err := aggregatorInstance.WatchAnswerUpdated(nil, ansChan, nil, nil)
	if err != nil {
		logFeedError(feedHexAddress, tokens, "Failed subscribing to price updates", err)
		return
	}

	log.WithFields(log.Fields{
		"feedAddress": feedHexAddress,
		"tokens":  tokens,
	}).Info("Monitoring price")
	for {
		select {
		case err := <-sub.Err():
			logFeedError(feedHexAddress, tokens, "Failed while listening for events", err)
			return
		case <-ctx.Done():
			return
		case ans := <-ansChan:
			newPrice := new(big.Float).Quo(new(big.Float).SetInt(ans.Current), decimalsDivFloat)
			log.WithFields(log.Fields{
				"tokens": tokens,
				"price":      newPrice.Text('f', int(aggregatorDecimals)),
			}).Info("New price")
		}
	}
}

func main() {
	feedConf, err := parseFeedConfig("feed.yaml")
	if err != nil {
		log.Errorf("Failed parsing feed config: %s", err)
		return
	}

	background := context.Background()
	termCtx, termCancel := context.WithCancel(background)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		termCancel()
	}()

	timeoutCtx, timeoutCancel := context.WithTimeout(termCtx, queryTimeout)
	client, err := ethclient.DialContext(timeoutCtx, os.Getenv("ALCHEMY_URL"))
	timeoutCancel()
	if err != nil {
		log.Errorf("Failed dialing node: %s", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(feedConf.Feeds))

	go subscribeBlocks(termCtx, client, &wg)
	go subscribeEvents(termCtx, client, "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419", "ETH / USD", &wg)
	go subscribeEvents(termCtx, client, "0xdc530d9457755926550b59e8eccdae7624181557", "LINK / ETH", &wg)
	go subscribeEvents(termCtx, client, "0xee9f2375b4bdf6387aa8265dd4fb8f16512a1d46", "USDT / ETH", &wg)

	wg.Wait()
}
