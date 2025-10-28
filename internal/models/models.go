package models

import (
	"mime/multipart"
)

type UploadedFile struct {
	File     multipart.File
	Header   *multipart.FileHeader
	Filename string
	Size     int64
}

type ConvertRequest struct {
	MultiplePages bool `form:"multiple"`
}

type ExtractRequest struct {
	Pages string `form:"pages" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
