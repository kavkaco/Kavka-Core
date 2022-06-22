package models

type User struct {
	ID                  uint
	Name                string
	Email               string
	Username            string
	Bio                 string
	LastSeen            string
	Chats               []interface{}
	ProfilePhotos       []string
	VerificCode         uint
	VerificTryCount     uint
	VerificCodeExpire   uint
	VerificLimitDate    uint
	CompletedFirstLogin bool
}
