package hotstuff

import "bytes"

type VoteSet []*MsgPrepareVote

// Add return false if duplicate vote msg exist
func (vs VoteSet) Add(vt *MsgPrepareVote) bool {
	for _, v := range vs {
		if v.ViewNum == vt.ViewNum &&
			v.BlockHash == vt.BlockHash &&
			bytes.Equal(v.PartialSig, vt.PartialSig) {
			return false
		}
	}
	vs = append(vs, vt)
	return true
}

// todo
func (vs VoteSet) Marjor() bool {
	return true
}
