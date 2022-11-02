# New blocks and price monitor
Monitors new blocks on the Ethereum mainnet and price changes of specified token pairs
## How to specify token pairs
See `feed.yaml`, at the moment it looks like this:
```yaml
feeds:
  - tokens: "ETH / USD"
    address: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419"
  - tokens: "LINK / ETH"
    address: "0xdc530d9457755926550b59e8eccdae7624181557"
  - tokens: "USDT / ETH"
    address: "0xee9f2375b4bdf6387aa8265dd4fb8f16512a1d46"
```

## How to run
```shell
sudo docker build . -t blockchain-monitor
sudo docker run -e ALCHEMY_URL=<YOUR ALCHEMY URL> -it blockchain-monitor:latest
```

## Example of logs

```log

```