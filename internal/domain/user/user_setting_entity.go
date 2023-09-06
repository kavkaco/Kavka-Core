package user

type DeleteAccountTimeType int

const (
	OneMonth DeleteAccountTimeType = iota + 1
	ThreeMonth
	SixMonth
)

type Setting struct {
	ShowLastSeen      bool                  `json:"show_last_seen"`
	TwoFA             bool                  `json:"two_fa"`
	DeleteAccountTime DeleteAccountTimeType `json:"delete_account_time"`
}
