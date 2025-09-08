package tc_client

import "encoding/json"

type Policy struct{}
type State struct{}

type GenericBody struct {
	Body   interface{} `json:"body"`
	Header bodyHeader  `json:"header"`
}

type bodyHeader struct {
	Policy Policy `json:"policy"`
	State  State  `json:"state"`
}

func SerializeBody(body interface{}) ([]byte, error) {
	commonBody := GenericBody{
		Body: body,
		Header: bodyHeader{
			Policy: Policy{},
			State:  State{},
		},
	}
	json, err := json.Marshal(commonBody)
	return json, err
}
