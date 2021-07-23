package core

// Votes 当节点作为leader时，收集来自不同副本的投票消息，其应实现的功能包括添加投票
// 计算某个blockHash下的投票总量，判断是否满足阈值，以及对committed block hash对应高度的区块的投票进行删除
type messageSet struct {

}
