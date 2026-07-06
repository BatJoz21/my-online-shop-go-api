package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize     = 10 * 1024 * 1024
	UploadRoots     = "uploads/"
	ProductImageDir = "products/image/"
)

var allowedImageExtension = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var allowedImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

func SaveProductImage(file *multipart.FileHeader, context *gin.Context) (*string, error) {
	if file.Size > MaxFileSize {
		return nil, errors.New("File size is too large")
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtension[extension] {
		return nil, errors.New("Invalid file extension!")
	}

	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	n, err := openedFile.Read(buffer)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(buffer[:n])
	if !allowedImageMimeTypes[contentType] {
		return nil, errors.New("Invalid content type: " + contentType)
	}

	os.MkdirAll(UploadRoots+ProductImageDir, os.ModePerm)

	filename := fmt.Sprintf("product_image_%d%s", time.Now().UnixNano(), extension)
	path := filepath.Join(UploadRoots, ProductImageDir, filename)

	err = context.SaveUploadedFile(file, path)
	if err != nil {
		return nil, err
	}

	return &filename, nil
}
