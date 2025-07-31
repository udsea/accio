package image

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/mattn/go-isatty"
	"github.com/nfnt/resize"
)

// Constants for image processing
const (
	DefaultImageSize = 250 // Default size for profile images (250x250)
	CacheDir         = "./cache/images"
)

// ImageProcessor handles image processing and display
type ImageProcessor struct {
	CacheDir string
}

// NewImageProcessor creates a new image processor
func NewImageProcessor() *ImageProcessor {
	cacheDir := os.Getenv("PROFILE_IMAGE_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = CacheDir
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create cache directory: %v\n", err)
	}

	return &ImageProcessor{
		CacheDir: cacheDir,
	}
}

// ProcessImage processes an image from a reader
func (p *ImageProcessor) ProcessImage(imageData io.Reader, size int) (image.Image, error) {
	// Decode image
	img, _, err := image.Decode(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image
	if size <= 0 {
		size = DefaultImageSize
	}

	// Resize image while preserving aspect ratio
	resized := resize.Thumbnail(uint(size), uint(size), img, resize.Lanczos3)

	// Crop to square if needed
	width := resized.Bounds().Dx()
	height := resized.Bounds().Dy()

	if width != height {
		// Crop to square
		resized = imaging.CropCenter(resized, min(width, height), min(width, height))
	}

	return resized, nil
}

// SaveImageToCache saves an image to the cache
func (p *ImageProcessor) SaveImageToCache(img image.Image, username, platform string) (string, error) {
	// Create filename
	filename := fmt.Sprintf("%s_%s.png", platform, username)
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	// Create full path
	fullPath := filepath.Join(p.CacheDir, filename)

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode image as PNG
	if err := png.Encode(file, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return fullPath, nil
}

// GetCachedImagePath returns the path to a cached image
func (p *ImageProcessor) GetCachedImagePath(username, platform string) string {
	filename := fmt.Sprintf("%s_%s.png", platform, username)
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	return filepath.Join(p.CacheDir, filename)
}

// IsCachedImageAvailable checks if a cached image is available
func (p *ImageProcessor) IsCachedImageAvailable(username, platform string) bool {
	path := p.GetCachedImagePath(username, platform)
	_, err := os.Stat(path)
	return err == nil
}

// DisplayImage displays an image in the terminal
func (p *ImageProcessor) DisplayImage(img image.Image) (string, error) {
	// Check if terminal supports images
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return "", fmt.Errorf("terminal does not support images")
	}

	// Convert image to ASCII art for terminals that don't support images
	var buf bytes.Buffer

	// Encode image as PNG
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// For terminals that support iTerm2 image protocol
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Use iTerm2 image protocol
	imgData := buf.Bytes()
	encoded := encodeBase64(imgData)

	return fmt.Sprintf("\033]1337;File=inline=1;width=%dpx;height=%dpx:%s\a\n",
		width, height, encoded), nil
}

// encodeBase64 encodes data as base64
func encodeBase64(data []byte) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	var result strings.Builder
	result.Grow((len(data)*8 + 5) / 6)

	for i := 0; i < len(data); i += 3 {
		// Process up to 3 bytes at a time
		var chunk uint32
		chunkLen := 0

		for j := 0; j < 3 && i+j < len(data); j++ {
			chunk = (chunk << 8) | uint32(data[i+j])
			chunkLen++
		}

		// Pad if necessary
		for j := chunkLen; j < 3; j++ {
			chunk = chunk << 8
		}

		// Convert to base64
		for j := 0; j < 4; j++ {
			if j <= chunkLen {
				idx := (chunk >> (18 - j*6)) & 0x3F
				result.WriteByte(alphabet[idx])
			} else {
				result.WriteByte('=')
			}
		}
	}

	return result.String()
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
