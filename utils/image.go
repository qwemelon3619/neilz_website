package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

func ImageResize(filepath_ string) error {
	file, err := os.Open(filepath_)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// 1. 크기 조정
	newImg := resize(img, 750, 350)

	// 2. 자르기 (필요한 경우)
	// rect := image.Rect(0, 0, 750, 350)
	// croppedImg := newImg.(interface {
	// 	SubImage(r image.Rectangle) image.Image
	// }).SubImage(rect)

	// 3. 파일 저장 (원래 파일에 덮어쓰기)
	outFile, err := os.Create(filepath_) // 덮어쓰기
	if err != nil {
		return err
	}
	defer outFile.Close()
	ext := filepath.Ext(filepath_)
	log.Print(ext)
	if ext == ".png" {
		err = png.Encode(outFile, newImg) // Encode the cropped image
		if err != nil {
			return err
		}
	} else if ext == ".jpeg" || ext == ".jpg" {
		err = jpeg.Encode(outFile, newImg, nil) // Encode the cropped image
		if err != nil {
			return err
		}
	} else {
		return errors.New("file does not support")
	}

	return nil
}
func resize(img image.Image, width, height int) image.Image {
	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Perform simple nearest-neighbor resizing. (For better quality, consider implementing bilinear or bicubic interpolation.)
	srcBounds := img.Bounds()
	srcWidth := srcBounds.Max.X - srcBounds.Min.X
	srcHeight := srcBounds.Max.Y - srcBounds.Min.Y

	xScale := float64(srcWidth) / float64(width)
	yScale := float64(srcHeight) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * xScale)
			srcY := int(float64(y) * yScale)

			if srcX < 0 {
				srcX = 0
			}
			if srcY < 0 {
				srcY = 0
			}
			if srcX >= srcWidth {
				srcX = srcWidth - 1
			}
			if srcY >= srcHeight {
				srcY = srcHeight - 1
			}

			resizedImg.Set(x, y, img.At(srcX, srcY))
		}
	}
	return resizedImg
}
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
