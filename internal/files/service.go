package files

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/tangerinefrog/chatter/internal/crypto"
	"github.com/tangerinefrog/chatter/internal/storage"
)

type FileService struct {
	repo    *FilesRepository
	storage storage.FileStorage
	cipher  *crypto.Cipher
}

func NewFileService(repo *FilesRepository, storage storage.FileStorage, cipher *crypto.Cipher) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
		cipher:  cipher,
	}
}

func (s *FileService) UploadFile(
	ctx context.Context,
	chatID uuid.UUID,
	uploaderID uuid.UUID,
	r io.Reader,
	fileName string,
	mimeType string,
	sizeBytes int64,
) (*File, error) {
	id := uuid.New()
	fileKey := fmt.Sprintf("/upload/%s/%s", chatID, id.String())

	file := File{
		ID:         id,
		ChatID:     chatID,
		UploaderID: uploaderID,
		FileKey:    fileKey,
		FileName:   fileName,
		MimeType:   mimeType,
		SizeBytes:  sizeBytes,
	}

	err := s.repo.CreateFile(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file in DB: %w", err)
	}

	fileBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}
	fileEncrypted, err := s.cipher.Encrypt(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt file content: %w", err)
	}

	readerEncrypted := bytes.NewReader(fileEncrypted)

	err = s.storage.UploadFile(ctx, fileKey, readerEncrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	return &file, nil
}

func (s *FileService) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	file, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file from DB: %w", err)
	}

	err = s.storage.DeleteFile(ctx, file.FileKey)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	err = s.repo.DeleteFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file from DB: %w", err)
	}

	return nil
}

func (s *FileService) GetFile(ctx context.Context, fileID uuid.UUID) (*File, io.ReadCloser, error) {
	file, err := s.repo.GetFileByID(ctx, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file from DB: %w", err)
	}

	r, err := s.storage.GetFile(ctx, file.FileKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file URL: %w", err)
	}

	bytesEncrypted, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read encrypted file content: %w", err)
	}

	bytesDecrypted, err := s.cipher.Decrypt(bytesEncrypted)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt file content: %w", err)
	}

	readerDecrypted := io.NopCloser(bytes.NewReader(bytesDecrypted))

	return &file, readerDecrypted, nil
}
