package types

type User struct {
	Id               string `json:"id"`
	Email            string `json:"email"`
	Full_Name        string `json:"fullName"`
	Avatar_Image_Url string `json:"avatarImageUrl"`
	Member_Since     string `json:"memberSince"`
	Password_Hash    string `json:"passwordHash"`
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type UserLogin struct {
	User   User
	Tokens AuthTokens
}
