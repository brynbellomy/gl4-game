package movesys

type (
	MovementType int
)

const (
	MvmtNone MovementType = iota
	MvmtWalking
	MvmtSprinting
)
