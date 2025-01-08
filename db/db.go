package db

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	l2 "gorm.io/gorm/logger"
	"pack/model"
)

type SqlLiteDB struct {
	Db *gorm.DB
}

func NewSqlLiteDB(fileName string) *SqlLiteDB {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: l2.New(log.New(os.Stdout, "\r\n", log.LstdFlags), l2.Config{
			LogLevel: l2.Error,
			Colorful: true,
		}),
	})
	if err != nil {

		log.Fatalf("Connect db %s error:%v", fileName, err)
	}

	initTable(db)
	return &SqlLiteDB{
		Db: db,
	}
}

func initTable(db *gorm.DB) {
	tables := []interface{}{
		&model.ConcatImageEntity{},
	}
	for _, table := range tables {
		err := db.AutoMigrate(&table)
		if err != nil {

			log.Fatalf("failed to migrate:%v", err)
		}
	}
	log.Println("Table created successfully")
}

func (s *SqlLiteDB) InsertDetail(detail *model.ConcatImageDetail) error {
	entityL := detail.ToEntity()
	err := s.Db.Create(entityL).Error
	if err != nil {
		//log.Printf("insert detail error:%v", err)
		return err
	}
	return nil
}

func (s *SqlLiteDB) SelectByName(name string) (*model.ConcatImageDetail, error) {
	entity := new(model.ConcatImageEntity)
	err := s.Db.Model(entity).Where("out_put_name = ?", name).First(entity).Error
	if err != nil {
		return nil, err
	}
	return entity.ToDetail()
}
