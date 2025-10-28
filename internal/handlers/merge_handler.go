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

type MergeHandler struct {
	pdfService *services.PDFService
}

func NewMergeHandler(pdfService *services.PDFService) *MergeHandler {
	return &MergeHandler{
		pdfService: pdfService,
	}
}

func (h *MergeHandler) Merge(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Failed to parse multipart form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "No files uploaded"})
		return
	}

	uploadedFiles := make([]*models.UploadedFile, len(files))
	for i, file := range files {
		uploadedFile, err := h.createUploadedFile(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to process file: " + err.Error()})
			return
		}
		uploadedFiles[i] = uploadedFile
	}

	resultPath, err := h.pdfService.MergePDFs(uploadedFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Merge failed: " + err.Error()})
		return
	}

	h.sendFile(c, resultPath)
}

func (h *MergeHandler) createUploadedFile(file *multipart.FileHeader) (*models.UploadedFile, error) {
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

func (h *MergeHandler) sendFile(c *gin.Context, filePath string) {
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

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", string(rune(fileInfo.Size())))
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))

	io.Copy(c.Writer, file)
}
