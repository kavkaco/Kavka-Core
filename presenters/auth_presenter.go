package presenters

const (
	RefreshTokenHeaderName = "X-Refresh-Token"
	AccessTokenHeaderName  = "X-Access-Token"
)

type AuthRegisterRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
