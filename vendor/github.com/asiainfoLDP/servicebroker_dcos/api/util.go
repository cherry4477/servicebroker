package api

func getCredentials(token string) (string, string) {
	return "Authorization", "token=" + token
}
