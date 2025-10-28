# PDF Manipulator API

A Go-based REST API for PDF manipulation using Ghostscript. This service provides endpoints to convert images to PDF, extract pages from PDFs, and merge multiple PDFs.

## Features

- **Convert Images to PDF**: Convert single or multiple page images to PDF
- **Extract Pages**: Extract specific pages from a PDF document
- **Merge PDFs**: Combine multiple PDF files into a single document
- **Docker Support**: Containerized application with Ghostscript included
- **RESTful API**: Clean HTTP endpoints with proper error handling

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

## Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd manipulator-go
```

2. Build and run with Docker Compose:
```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`

### Local Development

1. Install Ghostscript on your system:
   - **Ubuntu/Debian**: `sudo apt-get install ghostscript`
   - **macOS**: `brew install ghostscript`
   - **Windows**: Download from [Ghostscript website](https://www.ghostscript.com/download/gsdnld.html)

2. Run the application:
```bash
go mod download
go run main.go
```

## API Endpoints

### Health Check
```
GET /health
```
Returns the service status.

### Convert Image to PDF
```
POST /convert
```

**Parameters:**
- `file` (multipart/form-data): Image file to convert
- `multiple` (optional): Set to `true` for multi-page conversion (returns ZIP file)

**Example:**
```bash
curl -X POST -F "file=@image.jpg" -F "multiple=false" http://localhost:8080/convert --output result.pdf
```

### Extract Pages from PDF
```
POST /extract
```

**Parameters:**
- `file` (multipart/form-data): PDF file
- `pages` (form): Page numbers to extract (e.g., "1,3,5" or "1-5")

**Example:**
```bash
curl -X POST -F "file=@document.pdf" -F "pages=1,3,5" http://localhost:8080/extract --output extracted.pdf
```

### Merge PDFs
```
POST /merge
```

**Parameters:**
- `files` (multipart/form-data): Multiple PDF files

**Example:**
```bash
curl -X POST -F "files=@file1.pdf" -F "files=@file2.pdf" -F "files=@file3.pdf" http://localhost:8080/merge --output merged.pdf
```

## Supported File Formats

- **Input**: PDF, JPEG, PNG, TIFF, BMP
- **Output**: PDF, JPEG (for single page conversion), ZIP (for multi-page conversion)

## Configuration

The service can be configured using environment variables:

- `PORT`: Server port (default: 8080)

## Error Handling

The API returns appropriate HTTP status codes:

- `200`: Success
- `400`: Bad Request (invalid parameters)
- `422`: Validation Error
- `500`: Internal Server Error

Error responses include a JSON object with an error message:
```json
{
  "error": "Error description"
}
```

## Development

### Project Structure
```
manipulator-go/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose configuration
├── internal/
│   ├── handlers/         # HTTP request handlers
│   ├── models/           # Data models
│   └── services/         # Business logic services
└── data/
    └── tmp/              # Temporary files directory
```

### Adding New Features

1. Add new service methods in `internal/services/pdf_service.go`
2. Create corresponding handlers in `internal/handlers/`
3. Update routes in `main.go`
4. Add tests for new functionality

## Troubleshooting

### Common Issues

1. **Ghostscript not found**: Ensure Ghostscript is installed and available in PATH
2. **Permission denied**: Check file permissions for the `data/tmp` directory
3. **Memory issues**: Large files may require increased Docker memory limits

### Logs

View application logs:
```bash
docker-compose logs -f manipulator-go
```

## License

This project is licensed under the MIT License.
