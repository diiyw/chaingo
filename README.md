# Chaingo

这是使用go语言开发的区块链。A simple Blockchain in Golang
# 依赖

GO > 1.11

# 操作
编译
```
$ export GO111MODULE=on
$ go build
```
创建钱包
```
$ chaingo account -new
address: UMuthrbWw7AEZ9k5vgqnSMajCfJAHuAb7
```
查看钱包余额
```
$ chaingo account -fund UMuthrbWw7AEZ9k5vgqnSMajCfJAHuAb7
balance: 0
```
创建创世块(请先删除data目录)
```
$ chaingo chain create
Mining to: UMuthrbWw7AEZ9k5vgqnSMajCfJAHuAb7
```
查看钱包余额
```
$ chaingo account -fund UMuthrbWw7AEZ9k5vgqnSMajCfJAHuAb7
balance: 25
```
查看区块
```
$ chaingo chain print
============ Block 00000c7cf58978148957fff32e4d217e5666167fd0c1b02eeb3a04cd1ad17513 ============
Height: 0
Prev. block:
--- Transaction 20f05ca079400566cf6379f0d5239455d7eee5f9bb16b2ea66126002701b8d36:
     Input 0:
       TXID:
       Out:       -1
       Signature:
       PubKey:    5468652054696d65732031362f4a616e2f32303138204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f7220776f726c64
     Output 0:
       Value:  25
       Script: 2c1fddf12934347f5c7f3897e2d530373decb418
```
# 说明

- 区块链存储采用leveldb

- 挖矿难度目前固定较低，出块不会太长

- 节点的发现目前只是硬写入了测试的一些端口，只是简单的服务

- 自带区块文件，钱包，通过chain_test.go可以创建创世块

- 挖矿奖励固定，可以无限挖（未实现减半）

- 可以自行编译，查看命令帮助

# 参考
[https://github.com/Jeiwan/blockchain_go](https://github.com/Jeiwan/blockchain_go)

[https://github.com/liuchengxu/blockchain-tutorial/blob/master/SUMMARY.md](https://github.com/liuchengxu/blockchain-tutorial/blob/master/SUMMARY.md)