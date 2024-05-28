package config

type OAuth struct {
	UserInfoURL string
}

func newOAuth() *OAuth {
	return &OAuth{
		UserInfoURL: "https://graph.instagram.com/v19.0/me",
	}
}
