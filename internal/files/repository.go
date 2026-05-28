package files

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/db"
)

type FilesRepository struct {
	q *db.Queries
}

func NewRepository(p *pgxpool.Pool) *FilesRepository {
	return &FilesRepository{
		q: db.New(p),
	}
}

func (r *FilesRepository) CreateFile(ctx context.Context, file File) error {
	err := r.q.CreateFile(ctx, db.CreateFileParams{
		ID:         pgtype.UUID{Bytes: file.ID, Valid: true},
		ChatID:     pgtype.UUID{Bytes: file.ChatID, Valid: true},
		UploaderID: pgtype.UUID{Bytes: file.UploaderID, Valid: true},
		FileKey:    file.FileKey,
		FileName:   file.FileName,
		MimeType:   file.MimeType,
		SizeBytes:  file.SizeBytes,
	})

	return err
}

func (r *FilesRepository) GetFileByID(ctx context.Context, fileID uuid.UUID) (File, error) {
	dbFile, err := r.q.GetFileByID(ctx, pgtype.UUID{Bytes: fileID, Valid: true})
	if err != nil {
		return File{}, err
	}

	return File{
		ID:         dbFile.ID.Bytes,
		ChatID:     dbFile.ChatID.Bytes,
		UploaderID: dbFile.UploaderID.Bytes,
		FileKey:    dbFile.FileKey,
		FileName:   dbFile.FileName,
		MimeType:   dbFile.MimeType,
		SizeBytes:  dbFile.SizeBytes,
	}, nil
}

func (r *FilesRepository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	return r.q.DeleteFile(ctx, pgtype.UUID{Bytes: fileID, Valid: true})
}

func (r *FilesRepository) LinkFileToMessage(ctx context.Context, fileID uuid.UUID, messageID uuid.UUID) error {
	return r.q.LinkFileToMessage(ctx, db.LinkFileToMessageParams{
		ID:        pgtype.UUID{Bytes: fileID, Valid: true},
		MessageID: pgtype.UUID{Bytes: messageID, Valid: true},
	})
}
