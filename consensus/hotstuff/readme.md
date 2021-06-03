## hotstuff protocol
this package implemented event-drive-hotstuff protocol.

#### 消息类型
* 广播消息: 
MsgProposal: 应该包含block
MsgVote: 应该包含block hash, block number, signature

* 内部消息:
HotstuffRequest: 从miner/worker seal时，将block打包到request中发送给共识引擎, 用pubsub去做
NextRound: 由paceMaker触发

#### valset
* 存储每个block对应的valset(后面再持久化)
* 对valset排序
* roundrobin算法(后面再改进)

#### 状态机
![steps 流程图](steps.png)

1.miner worker
* newWork 以太坊矿工worker通过taskLoop内闭包commit及定时器实现挖矿的定时处理
* prepare 根据worker当前高度及parent hash构造新的区块头，因为私钥在共识内保存，所以调用共识接口对区块头进行初始签名
* 
#### paceMaker

#### 空块的处理, 需不需要changeView

#### 信号触发的问题
nextView

#### 