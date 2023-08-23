package bearer

import "strings"

func ExtractFromHeader(authHeader string) string {
	token := strings.Split(authHeader, "Bearer ")
	return token[1]
}
