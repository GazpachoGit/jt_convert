package tc_client

type GetTicketRequest struct {
	Files []FileRequest `json:"files"`
}

type FileRequest struct {
	Uid      string `json:"uid"`
	TypeName string `json:"type"`
}

func GetGetTicketRequestBody(uid, typeName string) ([]byte, error) {
	body := GetTicketRequest{
		Files: []FileRequest{
			{
				Uid:      uid,
				TypeName: typeName,
			},
		},
	}
	return SerializeBody(body)
}
