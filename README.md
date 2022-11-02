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
INFO[0000] Monitoring blocks
INFO[0001] Monitoring price                              feedAddress=0xdc530d9457755926550b59e8eccdae7624181557 tokens="LINK / ETH"
INFO[0001] Monitoring price                              feedAddress=0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419 tokens="ETH / USD"
INFO[0001] Monitoring price                              feedAddress=0xee9f2375b4bdf6387aa8265dd4fb8f16512a1d46 tokens="USDT / ETH"
INFO[0010] New block                                     number=15884833 transactions=132
INFO[0021] New block                                     number=15884834 transactions=178
INFO[0034] New block                                     number=15884835 transactions=124
...
INFO[0201] New block                                     number=15884849 transactions=154
INFO[0213] New block                                     number=15884850 transactions=101
INFO[0226] New block                                     number=15884851 transactions=221
INFO[0238] New block                                     number=15884852 transactions=120
INFO[0249] New price                                     price=1515.23112562 tokens="ETH / USD"
INFO[0249] New block                                     number=15884853 transactions=123
INFO[0262] New block                                     number=15884854 transactions=156
INFO[0272] New price                                     price=1513.67239791 tokens="ETH / USD"
INFO[0272] New block                                     number=15884855 transactions=131
INFO[0285] New block                                     number=15884856 transactions=216
INFO[0298] New block                                     number=15884857 transactions=128
INFO[0309] New block                                     number=15884858 transactions=169
INFO[0321] New block                                     number=15884859 transactions=147
INFO[0334] New block                                     number=15884860 transactions=83
INFO[0345] New block                                     number=15884861 transactions=141
...
INFO[1821] New block                                     number=15884982 transactions=149
INFO[1832] New price                                     price=1505.82000000 tokens="ETH / USD"
INFO[1832] New block                                     number=15884983 transactions=95
...
```