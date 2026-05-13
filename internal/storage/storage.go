package storage

import "context"

type FileStorage interface {
	UploadFile(ctx context.Context, fileName string, data []byte) error
}
