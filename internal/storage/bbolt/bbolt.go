package bbolt

import (
	"encoding/json"
	"fmt"
	model "jt_converter/internal/storage/model/pmis"
	"log/slog"

	"go.etcd.io/bbolt"
)

const BucketName = "Models"

type Storage struct {
	db  *bbolt.DB
	log *slog.Logger
}

func New(storagePath string, log *slog.Logger) *Storage {
	db, err := bbolt.Open(storagePath+"\\pmi.db", 0666, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &Storage{
		db, log,
	}
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) SavePMIs(key string, m *model.Model) error {
	const op = "storage.bbolt.SavePMIs"
	log := s.log.With(slog.String("op", op), slog.String("key", key))
	log.Debug("start SavePMIs")
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("%w: %s", err, "error during CreateBucketIfNotExists")
		}
		data, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("%w: %s", err, "error during Model Marshaling")
		}
		return bucket.Put([]byte(key), data)
	})
	if err != nil {
		log.Error("error during bbolt SavePMIs transaction", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("successfully complete SavePMIs!")
	return nil
}

func (s *Storage) GetPMIs(keys []string) ([]*model.Model, error) {
	const op = "storage.bbolt.GetPMIs"
	log := s.log.With(slog.String("op", op))
	log.Debug("start GetPMIs")
	data := make([][]byte, 0, 10)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return fmt.Errorf("%s '%s'", "the db bucket not found", BucketName)
		}
		for _, key := range keys {
			val := bucket.Get([]byte(key))
			if val != nil {
				data = append(data, val)
			} else {
				log.Info("object not found in db", slog.String("key", key))
			}
		}
		return nil
	})
	if err != nil {
		log.Error("error during bbolt GetPMIs transaction", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("got some data from bbolt db")

	resp := make([]*model.Model, 0, 10)
	for _, d := range data {
		if d != nil {
			var m model.Model
			err = json.Unmarshal(d, &m)
			if err != nil {
				log.Debug("can't unmarshal db response", slog.String("err", err.Error()))
				continue
			}
			resp = append(resp, &m)
		}
	}
	log.Debug("successfully complete GetPMIs!")
	return resp, nil
}

func (s *Storage) GetKeysList() ([]string, error) {
	const op = "storage.bbolt.GetKeysList"
	log := s.log.With(slog.String("op", op))
	log.Debug("start GetKeysList")

	keys := make([]string, 0, 10)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return fmt.Errorf("%s '%s'", "the db bucket not found", BucketName)
		}

		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})
	if err != nil {
		log.Error("error during bbolt GetKeysList transaction", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("successfully complete GetKeysList!")
	return keys, nil
}

//TODO: delete handler
