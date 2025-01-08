package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

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

// 保存图片
func saveImage(img image.Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
}

// 拼接图片（按等宽或等高）
func concatImages(images []image.Image, isVertical bool) image.Image {
	var totalWidth, totalHeight int
	for _, img := range images {
		bounds := img.Bounds()
		if isVertical {
			totalWidth = bounds.Dx()
			totalHeight += bounds.Dy()
		} else {
			totalWidth += bounds.Dx()
			totalHeight = bounds.Dy()
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))
	offset := 0
	for _, img := range images {
		bounds := img.Bounds()
		if isVertical {
			rect := image.Rect(0, offset, bounds.Dx(), offset+bounds.Dy())
			draw.Draw(dst, rect, img, image.Point{}, draw.Over)
			offset += bounds.Dy()
		} else {
			rect := image.Rect(offset, 0, offset+bounds.Dx(), bounds.Dy())
			draw.Draw(dst, rect, img, image.Point{}, draw.Over)
			offset += bounds.Dx()
		}
	}
	return dst
}

// 获取图片分组（按等宽优先）
func groupImagesByDimensions(imagePaths []string) ([][]image.Image, error) {
	widthGroups := make(map[int][]image.Image)
	heightGroups := make(map[int][]image.Image)
	processed := make(map[string]bool)
	var groups [][]image.Image

	for _, path := range imagePaths {
		img, err := loadImage(path)
		if err != nil {
			return nil, err
		}
		bounds := img.Bounds()
		width, _ := bounds.Dx(), bounds.Dy()

		// 按宽分组
		if !processed[path] {
			widthGroups[width] = append(widthGroups[width], img)
			processed[path] = true
		}
	}

	// 将按宽分组的图片移动到最终分组
	for _, imgs := range widthGroups {
		if len(imgs) > 0 {
			groups = append(groups, imgs)
		}
	}

	// 按高分组
	for _, path := range imagePaths {
		if !processed[path] {
			img, err := loadImage(path)
			if err != nil {
				return nil, err
			}
			bounds := img.Bounds()
			height := bounds.Dy()
			heightGroups[height] = append(heightGroups[height], img)
			processed[path] = true
		}
	}

	// 将按高分组的图片移动到最终分组
	for _, imgs := range heightGroups {
		if len(imgs) > 0 {
			groups = append(groups, imgs)
		}
	}
	// 宽高都对不齐的
	for _, path := range imagePaths {
		if !processed[path] {
			img, err := loadImage(path)
			if err != nil {
				return nil, err
			}
			groups = append(groups, []image.Image{img})
		}
	}
	return groups, nil
}

func main() {
	inputDir := "." // 当前目录
	outputDir := "output"

	// 遍历目录获取所有 JPG 文件路径
	var imagePaths []string
	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".jpg" {
			imagePaths = append(imagePaths, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("读取目录失败: %v", err)
	}

	// 按等宽或等高分组
	groups, err := groupImagesByDimensions(imagePaths)
	if err != nil {
		log.Fatalf("分组失败: %v", err)
	}

	// 每 6 张图片拼接成一张长图
	for _, group := range groups {
		// 长宽分组都拼不齐的直接输出原图
		if len(group) == 1 {
			filename := filepath.Join(outputDir, fmt.Sprintf("%s.jpg", uuid.New().String()))
			err = saveImage(group[0], filename)
			if err != nil {
				log.Fatalf("保存图片失败: %v", err)
			}
			fmt.Printf("图片已保存: %s\n", filename)
			continue
		}
		for i := 0; i < len(group); i += 6 {
			end := i + 6
			if end > len(group) {
				end = len(group)
			}
			subset := group[i:end]
			// 检测拼接方向（等宽为垂直拼接，等高为水平拼接）
			isVertical := group[0].Bounds().Dx() == subset[0].Bounds().Dx()
			// 拼接图片
			concatImg := concatImages(subset, isVertical)
			// 随机生成文件名
			filename := filepath.Join(outputDir, fmt.Sprintf("%s.jpg", uuid.New().String()))
			// 保存拼接图片
			err = saveImage(concatImg, filename)
			if err != nil {
				log.Fatalf("保存图片失败: %v", err)
			}
			fmt.Printf("图片已保存: %s\n", filename)
		}
	}
	fmt.Println("所有图片处理完成！")
}
