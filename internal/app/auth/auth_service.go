package auth

import "Kavka/internal/domain/user"

type Service struct {
	userRepository user.Repository
}
