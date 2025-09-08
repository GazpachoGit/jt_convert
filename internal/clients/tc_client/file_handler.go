package tcclient

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func (c *TCClient) GetFile(uid, ticket string) ([]byte, error) {
	const op = "tcclient.GetFile"
	log := c.log.With(slog.String("op", op))
	ticketURL := url.QueryEscape(ticket)
	fullURL := c.tcURL + "/fms/fmsdownload/" + uid + ".jt" + "?ticket=" + ticketURL
	log.Info("prepared getFile request", slog.String("url", fullURL))
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Error("failed to create getFile request", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.Error("failed to send getFile request", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = c.processErrorStatus(resp, log)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read getFile response body", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully got the file", slog.String("uid", uid), slog.Int("size", len(data)))
	return data, nil
}
