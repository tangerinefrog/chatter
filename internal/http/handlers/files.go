package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/storage"
	"go.uber.org/zap"
)

type FileHandler struct {
	storage storage.FileStorage
	logger  *zap.Logger
}

func NewFileHandler(storage storage.FileStorage, logger *zap.Logger) *FileHandler {
	return &FileHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *FileHandler) UploadFile(c *echo.Context) error {
	fileBytes, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(400, "file is required")
	}

	src, err := fileBytes.Open()
	if err != nil {
		h.logger.Error("failed to open uploaded file", zap.Error(err))
		return echo.NewHTTPError(500, "failed to process file")
	}
	defer src.Close()

	fileData := make([]byte, fileBytes.Size)
	_, err = src.Read(fileData)
	if err != nil {
		h.logger.Error("failed to read uploaded file", zap.Error(err))
		return echo.NewHTTPError(500, "failed to process file")
	}

	err = h.storage.UploadFile(c.Request().Context(), fileBytes.Filename, fileData)
	if err != nil {
		h.logger.Error("failed to upload file", zap.Error(err))
		return echo.NewHTTPError(500, "failed to upload file")
	}

	return c.NoContent(http.StatusCreated)
}
