package handlers

import (
	"io"
	"manipulator-go/internal/models"
	"manipulator-go/internal/services"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ExtractHandler struct {
	pdfService *services.PDFService
}

func NewExtractHandler(pdfService *services.PDFService) *ExtractHandler {
	return &ExtractHandler{
		pdfService: pdfService,
	}
}

func (h *ExtractHandler) ExtractPages(c *gin.Context) {
	var req models.ExtractRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request: " + err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "No file uploaded"})
		return
	}

	uploadedFile, err := h.createUploadedFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to process file: " + err.Error()})
		return
	}

	resultPath, err := h.pdfService.ExtractPages(uploadedFile, req.Pages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Extraction failed: " + err.Error()})
		return
	}

	h.sendFile(c, resultPath)
}

func (h *ExtractHandler) createUploadedFile(file *multipart.FileHeader) (*models.UploadedFile, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	return &models.UploadedFile{
		File:     src,
		Header:   file,
		Filename: file.Filename,
		Size:     file.Size,
	}, nil
}

func (h *ExtractHandler) sendFile(c *gin.Context, filePath string) {
	defer os.Remove(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to read result file"})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get file info"})
		return
	}

	ext := filepath.Ext(filePath)
	contentType := "application/octet-stream"
	if ext == ".pdf" {
		contentType = "application/pdf"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", string(rune(fileInfo.Size())))
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))

	io.Copy(c.Writer, file)
}
