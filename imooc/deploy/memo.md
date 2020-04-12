# 笔记

## 设置工作路径
``` bash
export FABRIC_CFG_PATH=$GOPATH/src/github.com/hyperledger/fabric/imooc/deploy
```

## 环境清理
``` bash
rm -fr config/*
rm -fr crypto-config/*
```

## 生成证书文件
``` bash
cryptogen generate --config=./crypto-config.yaml
```

## 生成创世区块
``` bash
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block
```

## 生成通道的创世交易
``` bash
configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/mychannel.tx -channelID mychannel
configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/assetschannel.tx -channelID assetschannel
```

## 生成组织关于通道的锚节点（主节点）交易
``` bash
configtxgen -profile TwoOrgChannel -outputAnchorPeersUpdate ./config/Org0MSPanchors.tx -channelID mychannel -asOrg Org0MSP
configtxgen -profile TwoOrgChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
```

## 创建通道 
``` bash
peer channel create -o orderer.imocc.com:7050 -c mychannel -f /etc/hyperledger/config/mychannel.tx
peer channel create -o orderer.imocc.com:7050 -c assetschannel -f /etc/hyperledger/config/assetschannel.tx
```

## 加入通道
``` bash
peer channel join -b mychannel.block
peer channel join -b assetschannel.block
```

## 设置主节点
``` bash
peer channel update -o orderer.imocc.com:7050 -c mychannel -f /etc/hyperledger/config/Org1MSPanchors.tx
```

## 链码安装
``` bash
peer chaincode install -n assetsexchange -v 1.0.0 -l golang -p github.com/chaincode/assetsExchange
peer chaincode install -n badexample -v 1.0.0 -l golang -p github.com/chaincode/badexample
最终版本（因为错误导致我们把链码名字换过来） assetsexchange->assets
peer chaincode install -n assets -v 1.0.0 -l golang -p github.com/chaincode/assetsExchange
```

## 链码实例化
``` bash
peer chaincode instantiate -o orderer.imocc.com:7050 -C assetschannel -n assetsexchange -l golang -v 1.0.0 -c '{"Args":["init"]}' -P "OR('org0MSP.member','org1MSP.admin')"
peer chaincode instantiate -o orderer.imocc.com:7050 -C mychannel -n badexample -l golang -v 1.0.0 -c '{"Args":["init"]}'
最终版本assetsexchange->assets, -P 去掉
peer chaincode instantiate -o orderer.imocc.com:7050 -C assetschannel -n assets -l golang -v 1.0.0 -c '{"Args":["init"]}'
```

## 链码交互
``` bash
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["userRegister", "user1", "user1"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["assetEnroll", "asset1", "asset1", "metadata", "user1"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["userRegister", "user2", "user2"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["assetExchange", "user1", "asset1", "user2"]}'
peer chaincode invoke -C assetschannel -n assets -c '{"Args":["userDestroy", "user1"]}'
```

## 链码升级(升级完还是失败)
``` bash
最后将assetsexchange->assets, -P策略去除
后来的二次升级是因为源代码又bug -v 1.0.0 -> -v 1.0.1
peer chaincode install -n assets -v 1.0.1 -l golang -p github.com/chaincode/assetsExchange
peer chaincode upgrade -C assetschannel -n assets -v 1.0.1 -c '{"Args":[""]}'
自己代码还写错了继续升级
peer chaincode install -n assets -v 1.0.2 -l golang -p github.com/chaincode/assetsExchange
peer chaincode upgrade -C assetschannel -n assets -v 1.0.2 -c '{"Args":[""]}'
```



## 链码查询
``` bash
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryUser", "user1"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAsset", "asset1"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryUser", "user2"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAssetHistory", "asset1"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAssetHistory", "asset1", "all"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAssetHistory", "asset1", "enroll"]}'
peer chaincode query -C assetschannel -n assets -c '{"Args":["queryAssetHistory", "asset1", " "]}'
peer chaincode query -C mychannel -n badexample -c '{"Args":[]}'
```

## 命令行模式的背书策略

EXPR(E[,E...])
EXPR = OR AND
E = EXPR(E[,E...])
MSP.ROLE
MSP 组织名 org0MSP org1MSP
ROLE admin member

OR('org0MSP.member','org1MSP.admin')
