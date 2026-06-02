package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/files"
	"github.com/tangerinefrog/chatter/internal/http/dto"
	"go.uber.org/zap"
)

type FileHandler struct {
	logger      *zap.Logger
	fileService *files.FileService
}

func NewFileHandler(fileService *files.FileService, logger *zap.Logger) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		logger:      logger,
	}
}

func (h *FileHandler) UploadFile(c *echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	chatID, ok := c.Get("chat_id").(uuid.UUID)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error("failed to open uploaded file", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}
	defer file.Close()

	createdFile, err := h.fileService.UploadFile(
		c.Request().Context(),
		chatID,
		userID,
		file,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		fileHeader.Size,
	)
	if err != nil {
		h.logger.Error("failed to upload file", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	fileUrl := fmt.Sprintf("/api/chats/%s/files/%s", chatID.String(), createdFile.ID.String())

	resp := dto.NewFileResponse{
		ID:        createdFile.ID.String(),
		Name:      createdFile.FileName,
		MimeType:  createdFile.MimeType,
		SizeBytes: createdFile.SizeBytes,
		Url:       fileUrl,
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *FileHandler) DownloadFile(c *echo.Context) error {
	fileIDParam := c.Param("fileID")
	fileID, err := uuid.Parse(fileIDParam)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	file, r, err := h.fileService.GetFile(c.Request().Context(), fileID)
	if err != nil {
		h.logger.Error("failed to get file URL", zap.String("FileID", fileID.String()), zap.Error(err))
		return c.NoContent(http.StatusNotFound)
	}
	defer r.Close()

	return c.Stream(http.StatusOK, file.MimeType, r)
}

func (h *FileHandler) DeleteFile(c *echo.Context) error {
	fileIDParam := c.Param("fileID")
	fileID, err := uuid.Parse(fileIDParam)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	err = h.fileService.DeleteFile(c.Request().Context(), fileID)

	if err != nil {
		h.logger.Error("failed to delete file", zap.String("FileID", fileID.String()), zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
