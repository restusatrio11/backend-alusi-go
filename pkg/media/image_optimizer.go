package media

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	// MaxFileSize represents 2 MB maximum upload size for logos
	MaxFileSize = 2 * 1024 * 1024
	// MaxDimension is the maximum width or height for raster logos (512px)
	MaxDimension = 512
)

var (
	ErrFileTooLarge    = errors.New("ukuran file melebihi batas maksimum 2 MB")
	ErrUnsupportedType = errors.New("format file tidak didukung (hanya .png, .jpg, .jpeg, .webp, .svg)")
	ErrCorruptedImage  = errors.New("file gambar rusak atau tidak dapat diproses")
	ErrEmptyFile       = errors.New("file upload tidak boleh kosong")

	scriptRegex         = regexp.MustCompile(`(?i)<script[\s\S]*?</script>`)
	eventHandlerRegex   = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	javascriptLinkRegex = regexp.MustCompile(`(?i)href\s*=\s*["']\s*javascript:[^"']*["']`)
	xmlEntityRegex      = regexp.MustCompile(`(?i)<!ENTITY[\s\S]*?>`)
)

// ImageOptimizer handles image validation, compression, resizing, and SVG sanitization
type ImageOptimizer struct {
	uploadDir string
	baseURL   string
}

// NewImageOptimizer creates an optimizer with specified upload destination
func NewImageOptimizer(uploadDir string, baseURL string) *ImageOptimizer {
	if uploadDir == "" {
		uploadDir = "./uploads/logos"
	}
	if baseURL == "" {
		baseURL = "/uploads/logos"
	}
	return &ImageOptimizer{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// EnsureUploadDir creates upload directory if it does not exist
func (o *ImageOptimizer) EnsureUploadDir() error {
	return os.MkdirAll(o.uploadDir, 0755)
}

// ProcessAndSave processes an uploaded multipart file header, optimizes/compresses it, and returns the public URL
func (o *ImageOptimizer) ProcessAndSave(fileHeader *multipart.FileHeader, identifier string) (string, error) {
	if fileHeader == nil {
		return "", ErrEmptyFile
	}

	if fileHeader.Size <= 0 {
		return "", ErrEmptyFile
	}

	if fileHeader.Size > MaxFileSize {
		return "", ErrFileTooLarge
	}

	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file upload: %w", err)
	}
	defer srcFile.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, srcFile); err != nil {
		return "", fmt.Errorf("gagal membaca data file: %w", err)
	}

	fileBytes := buf.Bytes()
	if len(fileBytes) == 0 {
		return "", ErrEmptyFile
	}

	if err := o.EnsureUploadDir(); err != nil {
		return "", fmt.Errorf("gagal menyiapkan direktori penyimpanan: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	mimeType := http.DetectContentType(fileBytes)

	// Clean identifier for safe filename
	cleanID := strings.ToLower(identifier)
	cleanID = regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(cleanID, "-")
	if cleanID == "" {
		cleanID = "logo"
	}

	uniqueSuffix := uuid.New().String()[:8]

	// 1. Handle SVG Vector Files
	if ext == ".svg" || strings.Contains(mimeType, "svg") || (strings.Contains(mimeType, "text/xml") && bytes.Contains(fileBytes, []byte("<svg"))) {
		sanitizedSVG, err := o.SanitizeSVG(fileBytes)
		if err != nil {
			return "", err
		}

		filename := fmt.Sprintf("logo-%s-%s.svg", cleanID, uniqueSuffix)
		targetPath := filepath.Join(o.uploadDir, filename)

		if err := os.WriteFile(targetPath, sanitizedSVG, 0644); err != nil {
			return "", fmt.Errorf("gagal menyimpan file SVG: %w", err)
		}

		return fmt.Sprintf("%s/%s", strings.TrimRight(o.baseURL, "/"), filename), nil
	}

	// 2. Handle Raster Files (PNG, JPEG, WebP)
	isPNG := strings.Contains(mimeType, "png") || ext == ".png"
	isJPEG := strings.Contains(mimeType, "jpeg") || ext == ".jpg" || ext == ".jpeg"
	isWebP := strings.Contains(mimeType, "webp") || ext == ".webp"

	if !isPNG && !isJPEG && !isWebP {
		return "", ErrUnsupportedType
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCorruptedImage, err)
	}

	// Scale down if dimensions exceed MaxDimension
	optimizedImg := o.resizeIfNeeded(img, MaxDimension)

	var outBuf bytes.Buffer
	var outExt string

	if isJPEG {
		outExt = "jpg"
		err = jpeg.Encode(&outBuf, optimizedImg, &jpeg.Options{Quality: 85})
	} else {
		// PNG encoder with BestCompression ensures crisp alpha channel and minimal size (~20-60KB)
		outExt = "png"
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&outBuf, optimizedImg)
	}

	if err != nil {
		return "", fmt.Errorf("gagal mengompresi gambar: %w", err)
	}

	filename := fmt.Sprintf("logo-%s-%s.%s", cleanID, uniqueSuffix, outExt)
	targetPath := filepath.Join(o.uploadDir, filename)

	if err := os.WriteFile(targetPath, outBuf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("gagal menyimpan gambar terkompresi: %w", err)
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(o.baseURL, "/"), filename), nil
}

// SanitizeSVG strips active scripts, malicious event handlers, and entity expansions from SVG
func (o *ImageOptimizer) SanitizeSVG(data []byte) ([]byte, error) {
	if !bytes.Contains(data, []byte("<svg")) && !bytes.Contains(data, []byte("<SVG")) {
		return nil, ErrUnsupportedType
	}

	cleaned := scriptRegex.ReplaceAll(data, []byte(""))
	cleaned = eventHandlerRegex.ReplaceAll(cleaned, []byte(""))
	cleaned = javascriptLinkRegex.ReplaceAll(cleaned, []byte(""))
	cleaned = xmlEntityRegex.ReplaceAll(cleaned, []byte(""))

	return cleaned, nil
}

// resizeIfNeeded scales down an image proportionally if width or height exceeds maxDim
func (o *ImageOptimizer) resizeIfNeeded(src image.Image, maxDim int) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		return src
	}

	var targetW, targetH int
	if w >= h {
		targetW = maxDim
		targetH = (h * maxDim) / w
	} else {
		targetH = maxDim
		targetW = (w * maxDim) / h
	}

	if targetW < 1 {
		targetW = 1
	}
	if targetH < 1 {
		targetH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}
