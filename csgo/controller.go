package csgo

type csgoController interface {
	ClientChange()
	ModeChange()
	WorldChange()

	RoundStart()
	RoundEnd()
}

//type TourneyMode uint8
//
//const (
//	ModeConnecting TourneyMode = iota
//	ModeKnife
//	ModeVote
//	ModeWarmUp
//	ModePlay
//	ModeDone = ^TourneyMode(0)
//)
