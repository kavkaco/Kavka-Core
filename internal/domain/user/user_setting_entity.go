package user

type DeleteAccountTimeType int

const (
	OneMonth DeleteAccountTimeType = iota + 1
	ThreeMonth
	SixMonth
)

type Setting struct {
	ShowLastSeen      bool
	TwoFA             bool
	DeleteAccountTime DeleteAccountTimeType
}
