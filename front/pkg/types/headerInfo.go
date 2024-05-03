package types

type HeaderUser struct {
	Email     string
	Firstname string
}

type HeaderInfo struct {
	Title string
	User  HeaderUser
}
