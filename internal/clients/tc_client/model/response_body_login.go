package tc_client

import "encoding/json"

type LoginResponseBody struct {
	ServerInfo ServerInfoData `json:"serverInfo"`
}

type ServerInfoData struct {
	LogFile string `json:"logFile"`
}

func DeserializeLoginResponseBody(data []byte) (LoginResponseBody, error) {
	var loginResponseBody LoginResponseBody
	err := json.Unmarshal(data, &loginResponseBody)
	if err != nil {
		return loginResponseBody, err
	}
	return loginResponseBody, nil
}
