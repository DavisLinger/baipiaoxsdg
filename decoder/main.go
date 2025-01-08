package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"pack/db"
)

func main() {
	inputDir := "." // 当前目录
	outputDir := "output"
	dbPath := `../data.db`
	// 遍历目录获取所有 JPG 文件路径
	dataSource := db.NewSqlLiteDB(dbPath)
	var imagePaths []string
	_ = filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".jpg" {
			imagePaths = append(imagePaths, path)
		}
		return nil
	})
	for in := range imagePaths {
		name := strings.ReplaceAll(imagePaths[in], ".jpg", "")
		detail, err := dataSource.SelectByName(name)
		if err != nil {
			log.Printf("get image detail failed,err:%v", err)
			continue
		}
		srcImg, err := loadImage(imagePaths[in])
		if err != nil {
			log.Printf("load image failed,err:%v", err)
			continue
		}
		log.Printf("load image:%v", detail.Records)
		for _, info := range detail.Records {
			//dst := image.NewRGBA(image.Rect(0, 0, info.X, info.Y))

			// 裁剪图像
			cropped := srcImg.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(info.OffsetX, info.OffsetY, info.X, info.Y))

			//rect := image.Rect(0, 0, info.X, info.Y)

			//draw.Draw(dst, rect, srcImg, image.Point{}, draw.Over)
			filename := filepath.Join(outputDir, info.Name)
			err = saveImage(cropped, filename)
			if err != nil {
				log.Fatalf("保存图片失败: %v", err)
			}
			fmt.Printf("图片已保存: %s\n", filename)
		}
	}
}

// 保存图片
func saveImage(img image.Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}

// 加载图片
func loadImage(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
