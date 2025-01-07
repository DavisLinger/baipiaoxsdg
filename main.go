package main

import (
	"image/jpeg"
	"log"
	"os"
	"strings"
)

type imgDetailEntity struct {
	width, height int
	name          string
}

func main() {
	// 一次最多拼接6张
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取文件目录失败,err:%v", err)
	}
	entities, err := os.ReadDir(pwd)
	if err != nil {
		log.Fatalf("读取文件夹失败,err:%v", err)
	}
	var imgList = make([]os.DirEntry, 0)
	for _, entity := range entities {
		entityName := entity.Name()
		if strings.Contains(strings.ToLower(entityName), ".jpg") && !entity.IsDir() {
			imgList = append(imgList, entity)
		}
	}
	readJpgs(imgList)
}

func readJpgs(imgs []os.DirEntry) []*imgDetailEntity {
	res := make([]*imgDetailEntity, 0)
	for _, img := range imgs {
		f, err := os.Open(img.Name())
		if err != nil {
			log.Fatal(err)
		}
		content, err := jpeg.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		wd := content.Bounds().Max.X - content.Bounds().Min.X
		height := content.Bounds().Max.Y - content.Bounds().Min.Y
		log.Printf("name:%v,width:%v,height:%v", img.Name(), wd, height)
		res = append(res, &imgDetailEntity{wd, height, img.Name()})
		f.Close()
	}
	return res
}

func groupByWidth(entity []*imgDetailEntity) (map[int][]*imgDetailEntity, []*imgDetailEntity) {
	mp := make(map[int][]*imgDetailEntity)
	for in := range entity {
		mp[entity[in].width] = append(mp[entity[in].width], entity[in])
	}
	res := make([]*imgDetailEntity, 0)
	for key, entities := range mp {
		if len(entities) > 3 {
			continue
		} else {
			res = append(res, entities...)
			delete(mp, key)
		}
	}
	return mp, res
}

func groupByHeight(entity []*imgDetailEntity) (map[int][]*imgDetailEntity, []*imgDetailEntity) {
	mp := make(map[int][]*imgDetailEntity)
	for in := range entity {
		mp[entity[in].height] = append(mp[entity[in].height], entity[in])
	}
	res := make([]*imgDetailEntity, 0)
	for key, entities := range mp {
		if len(entities) > 3 {
			continue
		} else {
			res = append(res, entities...)
			delete(mp, key)
		}
	}
	return mp, res
}
