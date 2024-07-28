package stream

import (
	"fmt"
	"testing"

	"github.com/kavkaco/Kavka-Core/utils/random"
)

var usersList = func() []StreamSubscribedUser {
	var users []StreamSubscribedUser

	for i := 0; i < 1_000_000; i++ {
		users = append(users, StreamSubscribedUser{UserID: fmt.Sprintf("%v", random.GenerateUserID())})
	}

	// Append a specific user id to find later
	users = append(users, StreamSubscribedUser{UserID: "12345678"})

	return users
}()

func BenchmarkMatchUserSub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, u := range usersList {
			if u.UserID == "12345678" {
				return
			}
		}
	}
}
