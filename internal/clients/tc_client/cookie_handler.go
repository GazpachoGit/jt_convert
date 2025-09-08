package tcclient

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (c *TCClient) GetInitialCookies() error {
	const op = "tc_client.GetInitialCookies"
	log := c.log.With(slog.String("op", op))
	log.Info("getting initial cookies from TC")
	resp, err := c.client.Get(c.tcURL + "/")
	if err != nil {
		log.Error("can't reach TC starting page to tc", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = c.processErrorStatus(resp, log)
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully got the initial token values", slog.Int("status", resp.StatusCode))

	c.processCookies(resp.Request.URL, log)
	return nil
}
