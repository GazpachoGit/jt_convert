package tcclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	model "jt_converter/internal/clients/tc_client/model"
	"log/slog"
	"net/http"
)

func (c *TCClient) Login() error {
	const op = "tc_client.Login"
	log := c.log.With(slog.String("op", op))

	input, err := model.GetLoginRequestBody(c.user, c.password)
	if err != nil {
		log.Error("failed to create login request body", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("prepared login request body", slog.String("body", string(input)))
	req, err := http.NewRequest("POST", c.tcURL+"/tc/RestServices/Core-2011-06-Session/login", bytes.NewBuffer(input))
	if err != nil {
		log.Error("failed to create login request", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-XSRF-TOKEN", c.tokenForHeader)
	resp, err := c.client.Do(req)
	if err != nil {
		log.Error("failed to send login request", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = c.processErrorStatus(resp, log)
		return fmt.Errorf("%s: %w", op, err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read Login response body", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("login response", slog.String(("body"), string(data)))
	loginResp, err := model.DeserializeLoginResponseBody(data)
	if err != nil {
		log.Error("failed to deserialize Login response body", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	if loginResp.ServerInfo.LogFile == "" {
		return errors.New("bad response from login request")
	}

	log.Info("successfully logged in", slog.Int("status", resp.StatusCode))

	c.processCookies(resp.Request.URL, log)

	return nil
}
