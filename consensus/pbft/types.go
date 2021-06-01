package pbft

type HRS struct {
	Height uint64
	Round  int64
	Step   int64
}

func (h *HRS) Cmp(h1 *HRS) int {
	if h.Height > h1.Height {
		return 1
	} else if h.Height < h1.Height {
		return -1
	} else if h.Round > h1.Round {
		return 1
	} else if h.Round < h1.Round {
		return -1
	} else if h.Step > h1.Step {
		return 1
	} else if h.Step < h1.Step {
		return -1
	} else {
		return 0
	}
}

type RoundStep uint8

const (
	RoundStepNewHeight     = 1
	RoundStepNewRound      = 2
	RoundStepPreVote       = 3
	RoundStepPreVoteWait   = 4
	RoundStepPreCommit     = 5
	RoundStepPreCommitWait = 6
	RoundStepCommit        = 7
)
