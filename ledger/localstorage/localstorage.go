package localstorage

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/storage"
)

// LocalStorage allows a node to save local data.
type LocalStorage struct {
	db *storage.DB
}

// NewLocalStorage create new storage instance.
func NewLocalStorage(db *storage.DB) (*LocalStorage, error) {
	return &LocalStorage{db: db}, nil
}

// SetMessage saves message in storage.
func (s *LocalStorage) SetMessage(ctx context.Context, msg core.SignedMessage) (*core.RecordID, error) {
	buff, err := message.ToBytes(msg)
	if err != nil {
		return nil, err
	}

	return s.db.SetBlob(ctx, msg.Pulse(), buff)
}

// GetMessage retrieves message from storage.
func (s *LocalStorage) GetMessage(ctx context.Context, id core.RecordID) (core.SignedMessage, error) {
	buff, err := s.db.GetBlob(ctx, &id)
	if err != nil {
		return nil, err
	}

	return message.DeserializeSigned(bytes.NewBuffer(buff))
}
