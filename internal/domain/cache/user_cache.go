package cache

type UserCacheRepository interface {
	SetToken(token string, username string) error
	GetToken(token string) (string, error) // returns: username & error
}
