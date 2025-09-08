package tcclient

import (
	"bytes"
	"fmt"
	"io"
	model "jt_converter/internal/clients/tc_client/model"
	"log/slog"
	"net/http"
)

func (c *TCClient) GetTicket(uid, typeName string) (string, error) {
	const op = "tc_client.GetTicket"
	log := c.log.With(slog.String("op", op))

	input, err := model.GetGetTicketRequestBody(uid, typeName)
	if err != nil {
		log.Error("failed to create getTicket request body", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("prepared getTicket request body", slog.String("body", string(input)))

	req, err := http.NewRequest("POST", c.tcURL+"/tc/RestServices/Core-2006-03-FileManagement/getFileReadTickets", bytes.NewBuffer(input))
	if err != nil {
		log.Error("failed to create login request", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-XSRF-TOKEN", c.tokenForHeader)
	resp, err := c.client.Do(req)

	if err != nil {
		log.Error("failed to send getTicket request", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = c.processErrorStatus(resp, log)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read getTicket response body", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	ticketResp, err := model.DeserializeTicketResponseBody(data)
	if err != nil {
		log.Error("failed to deserialize getTicket response body", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err = ticketResp.GetPartialErrors(); err != nil {
		log.Error("getTicket response contains partial errors", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	ticket, err := model.GetTicket(ticketResp)
	if err != nil {
		log.Error("failed to find tickets in the response", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully got ticket", slog.String("ticket", ticket))
	return ticket, nil
}
