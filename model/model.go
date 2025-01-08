package model

import (
	"encoding/json"
	"image"

	"gorm.io/gorm"
)

type ConcatImageDetail struct {
	OutPutName string `json:"outPutName"`
	Records    []ImageInfo
}

type ImageInfo struct {
	Name    string `json:"name"`
	OffsetX int    `json:"offsetX"` //大图中的X偏移
	OffsetY int    `json:"offsetY"` //大图的Y偏移
	X       int    `json:"x"`       //自己图片的宽
	Y       int    `json:"y"`       //自身图片的高
}

func (c *ConcatImageDetail) ToEntity() *ConcatImageEntity {
	d, _ := json.Marshal(c.Records)
	return &ConcatImageEntity{
		OutPutName: c.OutPutName,
		Detail:     string(d),
	}
}

type ConcatImageEntity struct {
	gorm.Model
	OutPutName string `json:"outPutName" gorm:"uniqueIndex:uk_name"`
	Detail     string `json:"detail"`
}

func (c *ConcatImageEntity) TableName() string {
	return "concat_image"
}

func (c *ConcatImageEntity) ToDetail() (*ConcatImageDetail, error) {
	detail := new(ConcatImageDetail)
	err := json.Unmarshal([]byte(c.Detail), &detail.Records)
	if err != nil {
		return nil, err
	}
	detail.OutPutName = c.OutPutName
	return detail, nil
}

type ImageDetail struct {
	FileName string      `json:"fileName"`
	Img      image.Image `json:"img"`
}
