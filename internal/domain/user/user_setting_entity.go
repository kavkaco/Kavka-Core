package user

type DeleteAccountTimeType int

const (
	OneMonth DeleteAccountTimeType = iota + 1
	ThreeMonth
	SixMonth
)

type Setting struct {
	ShowLastSeen      bool                  `bson:"show_last_seen" json:"showLastSeen"`
	TwoFA             bool                  `bson:"two_fa" json:"twoFa"`
	DeleteAccountTime DeleteAccountTimeType `bson:"delete_account_time" json:"deleteAccountTime"`
}
