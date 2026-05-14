package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	UploadFile(ctx context.Context, fileKey string, r io.Reader) error
	DeleteFile(ctx context.Context, fileKey string) error
	GetFile(ctx context.Context, fileKey string) (io.ReadCloser, error)
}
