package core

type ConsensusState uint8

// todo: 更新状态类型，参考libraBFT
const (
	StateAcceptNewView ConsensusState = 1
	StateAcceptNewHeight ConsensusState = 2
	StateAcceptNewRound ConsensusState = 3
	StateAcceptLocked ConsensusState = 4
	StateAcceptCommitted ConsensusState = 5
)

type MessageType uint8

const (
	MsgTypeVote MessageType = 1
	MsgTypeProposal MessageType = 2
)

type voteMessage struct {

}

type proposalMessage struct {

}