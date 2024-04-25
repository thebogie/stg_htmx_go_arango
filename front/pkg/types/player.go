package types

type Player struct {
	Firstname   string `json:"firstname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"accessToken"`
}
