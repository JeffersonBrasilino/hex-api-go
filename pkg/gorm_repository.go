package pkg

import "gorm.io/gorm"

type GormRepository struct {
	Db *gorm.DB
}

func (*GormRepository) AutoMigrateModels(models ...*gorm.Model) {
	
}