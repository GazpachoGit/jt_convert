package service

import (
	"fmt"
	tcclient "jt_converter/internal/clients/tc_client"
	"log/slog"
	"os"
)

type TCService struct {
	c        *tcclient.TCClient
	log      *slog.Logger
	filesDir string
}

func NewTCService(c *tcclient.TCClient, log *slog.Logger, filesDir string) *TCService {
	return &TCService{c: c, log: log, filesDir: filesDir}
}

func (s *TCService) LoadFile(uid, typeName, Name string) error {
	const op = "service.LoadFile"
	log := s.log.With(slog.String("op", op))
	err := s.c.GetInitialCookies()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.c.Login()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	ticket, err := s.c.GetTicket(uid, typeName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	data, err := s.c.GetFile(uid, ticket)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	fileName := s.filesDir + "/" + Name + ".jt"
	if _, err := os.Stat(fileName); err == nil {
		if err := os.Remove(fileName); err != nil {
			log.Error("failed to delete existing file", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	out, err := os.Create(fileName)
	defer out.Close()
	if err != nil {
		log.Error("failed to create a local file", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = out.Write(data)
	if err != nil {
		log.Error("failed to write data to a local file", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully wrote data to a local file", slog.String("file", out.Name()), slog.Int("size", len(data)))
	return nil
}
