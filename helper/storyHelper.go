package helper

func CurrentUserStoryInteraction(arr []string, username string) bool {
	for _, u := range arr {
		if u == username {
			return true
		}
	}
	return false
}
