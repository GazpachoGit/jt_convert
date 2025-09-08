package tcclient

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const contentType = "application/json"

var (
	ErrBadStatus = fmt.Errorf("bad status from TC")
)

type TCClient struct {
	log            *slog.Logger
	client         *http.Client
	tcURL          string
	tokenForHeader string
	user, password string
}

func NewTCClient(tcURL, user, password string, log *slog.Logger) *TCClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}
	return &TCClient{
		log:      log,
		client:   client,
		tcURL:    tcURL,
		user:     user,
		password: password,
	}
}

func (c *TCClient) processCookies(u *url.URL, log *slog.Logger) {
	for _, cookie := range c.client.Jar.Cookies(u) {
		log.Info("cookie", slog.String("name", cookie.Name), slog.String("value", cookie.Value))
		if cookie.Name == "XSRF-TOKEN" {
			c.tokenForHeader = cookie.Value
		}
	}
}

func (c *TCClient) processErrorStatus(resp *http.Response, log *slog.Logger) error {
	log.Error("login request failed", slog.Int("status", resp.StatusCode))
	statusErr := fmt.Errorf("%w: %s", ErrBadStatus, resp.Status)
	if resp.ContentLength != 0 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("failed to read response body", slog.String("err", err.Error()))
			return errors.Join(err, statusErr)
		} else if len(bodyBytes) > 0 {
			log.Error("response body", slog.String("body", string(bodyBytes)))
			return fmt.Errorf("%w: %s", statusErr, string(bodyBytes))
		}
	}
	return fmt.Errorf("%w: %s", statusErr, "no extra info,body is empty")
}
