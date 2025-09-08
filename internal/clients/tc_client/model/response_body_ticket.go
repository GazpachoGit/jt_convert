package tc_client

import (
	"encoding/json"
	"fmt"
)

type TicketResponseBody struct {
	GenericResponseBody
	Tickets [][]interface{} `json:"tickets"`
}

func DeserializeTicketResponseBody(data []byte) (TicketResponseBody, error) {
	var ticketResponseBody TicketResponseBody
	err := json.Unmarshal(data, &ticketResponseBody)
	if err != nil {
		return ticketResponseBody, err
	}
	return ticketResponseBody, nil
}

func GetTicket(response TicketResponseBody) (string, error) {
	if response.Tickets != nil {
		for _, arrs := range response.Tickets {
			for _, ticket := range arrs {
				if ticketStr, ok := ticket.(string); ok {
					return ticketStr, nil
				}
			}
		}
	}
	return "", fmt.Errorf("no tickets found in response")
}
