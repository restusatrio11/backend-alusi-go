package media_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"backend-alusi-go/pkg/media"
)

func createMultipartHeader(filename string, content []byte) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	h.Set("Content-Type", "application/octet-stream")

	part, _ := writer.CreatePart(h)
	_, _ = part.Write(content)
	_ = writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(int64(body.Len()))
	return form.File["file"][0]
}

func createTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 128, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func createTestJPEG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 128, B: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	return buf.Bytes()
}

func TestImageOptimizer_PNG_ResizeAndCompress(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "alusi_test_uploads_png")
	defer os.RemoveAll(tempDir)

	optimizer := media.NewImageOptimizer(tempDir, "/uploads/logos")

	// Create 1000x800 PNG image
	rawPNG := createTestImage(1000, 800)
	header := createMultipartHeader("logo-simbatik.png", rawPNG)

	savedURL, err := optimizer.ProcessAndSave(header, "simbatik")
	if err != nil {
		t.Fatalf("Unexpected error processing PNG: %v", err)
	}

	if !strings.HasPrefix(savedURL, "/uploads/logos/logo-simbatik-") || !strings.HasSuffix(savedURL, ".png") {
		t.Fatalf("Unexpected saved URL: %s", savedURL)
	}

	// Verify file exists on disk and dimension is scaled down
	filename := filepath.Base(savedURL)
	diskPath := filepath.Join(tempDir, filename)

	f, err := os.Open(diskPath)
	if err != nil {
		t.Fatalf("Failed to open saved file: %v", err)
	}
	defer f.Close()

	decodedImg, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("Failed to decode saved image: %v", err)
	}

	bounds := decodedImg.Bounds()
	if bounds.Dx() > 512 || bounds.Dy() > 512 {
		t.Fatalf("Expected dimension <= 512px, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestImageOptimizer_JPEG(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "alusi_test_uploads_jpeg")
	defer os.RemoveAll(tempDir)

	optimizer := media.NewImageOptimizer(tempDir, "/uploads/logos")
	rawJPEG := createTestJPEG(600, 600)
	header := createMultipartHeader("portal.jpg", rawJPEG)

	savedURL, err := optimizer.ProcessAndSave(header, "portal")
	if err != nil {
		t.Fatalf("Unexpected error processing JPEG: %v", err)
	}

	if !strings.HasSuffix(savedURL, ".jpg") {
		t.Fatalf("Expected .jpg extension, got %s", savedURL)
	}
}

func TestImageOptimizer_SVG_Sanitization(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "alusi_test_uploads_svg")
	defer os.RemoveAll(tempDir)

	optimizer := media.NewImageOptimizer(tempDir, "/uploads/logos")

	svgMalicious := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
		<script>alert("XSS")</script>
		<circle cx="50" cy="50" r="40" stroke="green" stroke-width="4" fill="yellow" onload="evil()" />
		<a href="javascript:alert(1)"><text x="10" y="20">Click</text></a>
	</svg>`

	header := createMultipartHeader("vector-logo.svg", []byte(svgMalicious))
	savedURL, err := optimizer.ProcessAndSave(header, "vector-app")
	if err != nil {
		t.Fatalf("Unexpected error processing SVG: %v", err)
	}

	if !strings.HasSuffix(savedURL, ".svg") {
		t.Fatalf("Expected .svg URL, got %s", savedURL)
	}

	filename := filepath.Base(savedURL)
	diskPath := filepath.Join(tempDir, filename)
	savedContent, err := os.ReadFile(diskPath)
	if err != nil {
		t.Fatalf("Failed to read saved SVG: %v", err)
	}

	savedStr := string(savedContent)
	if strings.Contains(savedStr, "<script>") || strings.Contains(savedStr, "onload=") || strings.Contains(savedStr, "javascript:") {
		t.Fatalf("SVG was not properly sanitized: %s", savedStr)
	}
	if !strings.Contains(savedStr, "<circle") {
		t.Fatalf("Valid SVG tags were lost: %s", savedStr)
	}
}

func TestImageOptimizer_FileTooLarge(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "alusi_test_uploads_large")
	defer os.RemoveAll(tempDir)

	optimizer := media.NewImageOptimizer(tempDir, "/uploads/logos")

	// 3 MB dummy payload
	largeBytes := make([]byte, 3*1024*1024)
	header := createMultipartHeader("huge.png", largeBytes)

	_, err := optimizer.ProcessAndSave(header, "huge")
	if err == nil || !strings.Contains(err.Error(), "melebihi batas") {
		t.Fatalf("Expected file too large error, got: %v", err)
	}
}

func TestImageOptimizer_UnsupportedType(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "alusi_test_uploads_unsupported")
	defer os.RemoveAll(tempDir)

	optimizer := media.NewImageOptimizer(tempDir, "/uploads/logos")

	header := createMultipartHeader("document.pdf", []byte("%PDF-1.4 dummy content"))
	_, err := optimizer.ProcessAndSave(header, "doc")
	if err == nil || !strings.Contains(err.Error(), "tidak didukung") {
		t.Fatalf("Expected unsupported format error, got: %v", err)
	}
}

func TestImageOptimizer_EmptyFile(t *testing.T) {
	optimizer := media.NewImageOptimizer("", "")
	_, err := optimizer.ProcessAndSave(nil, "empty")
	if err == nil {
		t.Fatalf("Expected error on nil header")
	}
}

func createDummyReader(content []byte) io.Reader {
	return bytes.NewReader(content)
}
