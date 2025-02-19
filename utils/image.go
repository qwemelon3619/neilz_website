package utils

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
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
	rect := image.Rect(0, 0, 750, 350)
	croppedImg := newImg.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	// 3. 파일 저장 (원래 파일에 덮어쓰기)
	outFile, err := os.Create(filepath_) // 덮어쓰기
	if err != nil {
		return err
	}
	defer outFile.Close()
	ext := filepath.Ext(filepath_)
	log.Print(ext)
	if ext == ".png" {
		err = png.Encode(outFile, croppedImg) // Encode the cropped image
		if err != nil {
			return err
		}
	} else if ext == ".jpeg" || ext == ".jpg" {
		err = jpeg.Encode(outFile, croppedImg, nil) // Encode the cropped image
		if err != nil {
			return err
		}
	} else {
		return errors.New("file does not support")
	}

	return nil
}
func resize(img image.Image, width, height int) image.Image {
	// 새로운 이미지 생성
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// 이미지 비율에 맞춰서 조정
	srcBounds := img.Bounds()
	srcAspect := float64(srcBounds.Dx()) / float64(srcBounds.Dy())
	dstAspect := float64(width) / float64(height)

	var drawRect image.Rectangle
	if srcAspect > dstAspect {
		newHeight := int(float64(width) / srcAspect)
		yOffset := (height - newHeight) / 2
		drawRect = image.Rect(0, yOffset, width, yOffset+newHeight)
	} else {
		newWidth := int(float64(height) * srcAspect)
		xOffset := (width - newWidth) / 2
		drawRect = image.Rect(xOffset, 0, xOffset+newWidth, height)
	}

	// 이미지 그리기 (크기 조정)
	draw.Draw(newImg, drawRect, img, srcBounds.Min, draw.Src)

	return newImg
}
