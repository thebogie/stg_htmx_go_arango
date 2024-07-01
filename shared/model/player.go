package model

type LoginUserAPI struct {
	LoginUser struct {
		Token    string `json:"token"`
		Userdata struct {
			Id        string `json:"_id"`
			Key       string `json:"_key"`
			Rev       string `json:"_rev"`
			Email     string `json:"email"`
			Firstname string `json:"firstname"`
		} `json:"userdata"`
	} `json:"LoginUser"`
}

type Player struct {
	Id          string `json:"_id"`
	Key         string `json:"_key"`
	Rev         string `json:"_rev"`
	Firstname   string `json:"firstname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"accessToken"`
}
