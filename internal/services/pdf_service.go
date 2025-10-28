package services

import (
	"archive/zip"
	"fmt"
	"io"
	"manipulator-go/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) ConvertToPDF(file *models.UploadedFile, multiplePages bool) (string, error) {
	if multiplePages {
		return s.convertMultiplePages(file)
	}
	return s.convertSinglePage(file)
}

func (s *PDFService) ExtractPages(file *models.UploadedFile, pages string) (string, error) {
	inputPath := s.createTmpFile()
	outputPath := s.createTmpFile()

	if err := s.saveFile(file, inputPath); err != nil {
		return "", err
	}

	cmd := exec.Command("gs",
		"-q",
		"-dBATCH",
		"-dNOPAUSE",
		"-o", outputPath,
		"-sPageList="+pages,
		"-sDEVICE=pdfwrite",
		inputPath,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ghostscript extraction failed: %v", err)
	}

	return outputPath, nil
}

func (s *PDFService) MergePDFs(files []*models.UploadedFile) (string, error) {
	outputPath := s.createTmpFile()

	cmd := exec.Command("gs",
		"-q",
		"-dNOPAUSE",
		"-o", outputPath,
		"-sDEVICE=pdfwrite",
		"-dBATCH",
	)

	for _, file := range files {
		inputPath := s.createTmpFile()
		if err := s.saveFile(file, inputPath); err != nil {
			return "", err
		}
		cmd.Args = append(cmd.Args, inputPath)
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ghostscript merge failed: %v", err)
	}

	return outputPath, nil
}

func (s *PDFService) convertSinglePage(file *models.UploadedFile) (string, error) {
	inputPath := s.createTmpFile()
	outputPath := s.createTmpFile()

	if err := s.saveFile(file, inputPath); err != nil {
		return "", err
	}

	cmd := exec.Command("gs",
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-sDEVICE=jpeg",
		"-r300",
		"-dJPEGQ=90",
		"-dQUIET",
		"-dFirstPage=1",
		"-dLastPage=1",
		"-o", outputPath,
		inputPath,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ghostscript conversion failed: %v", err)
	}

	return outputPath, nil
}

func (s *PDFService) convertMultiplePages(file *models.UploadedFile) (string, error) {
	path := s.createTmpFile()
	zipPath := path + ".zip"

	if err := s.saveFile(file, path); err != nil {
		return "", err
	}

	outputPattern := path + "_page_%d"

	// Get page count
	pageCountCmd := exec.Command("gs",
		"-q",
		"-dNOSAFER",
		"-dNODISPLAY",
		"-c",
		fmt.Sprintf("(%s) (r) file runpdfbegin pdfpagecount = quit", path),
	)

	output, err := pageCountCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get page count: %v", err)
	}

	pages, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return "", fmt.Errorf("invalid page count: %v", err)
	}

	// Convert all pages to JPEG
	cmd := exec.Command("gs",
		"-dSAFER",
		"-dBATCH",
		"-dNOPAUSE",
		"-sDEVICE=jpeg",
		"-r300",
		"-o", fmt.Sprintf("\"%s\"", outputPattern),
		"-dJPEGQ=90",
		"-q",
		fmt.Sprintf("\"%s\"", path),
		"-c", "quit",
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ghostscript multi-page conversion failed: %v", err)
	}

	// Create ZIP file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for i := 1; i <= pages; i++ {
		pageFile := fmt.Sprintf(outputPattern, i)
		if err := s.addFileToZip(zipWriter, pageFile, fmt.Sprintf("%d.jpg", i)); err != nil {
			return "", fmt.Errorf("failed to add page %d to zip: %v", i, err)
		}
	}

	return zipPath, nil
}

func (s *PDFService) saveFile(file *models.UploadedFile, path string) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file.File); err != nil {
		return err
	}

	return nil
}

func (s *PDFService) addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(zipPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func (s *PDFService) createTmpFile() string {
	return filepath.Join("./data/tmp", uuid.New().String())
}
