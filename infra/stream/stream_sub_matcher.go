package stream

func MatchUserSubscription(userID string, users []StreamSubscribedUser) *StreamSubscribedUser {
	for _, u := range users {
		if u.UserID == userID {
			return &u
		}
	}

	return nil
}
