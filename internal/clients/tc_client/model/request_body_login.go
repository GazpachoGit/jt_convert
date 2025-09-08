package tc_client

type LoginRequest struct {
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Descrimator string `json:"descrimator"`
	Locale      string `json:"locale"`
	Group       string `json:"group"`
}

func GetLoginRequestBody(user, password string) ([]byte, error) {
	loginBody := LoginRequest{
		Credentials: Credentials{
			User:     user,
			Password: password,
		},
	}
	return SerializeBody(loginBody)
}
