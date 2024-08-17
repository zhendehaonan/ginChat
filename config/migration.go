package config

import (
	"ginchat/models"
)

// 对数据表进行迁移
func migration() {
	_ = _db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&models.UserBasic{},
		&models.Message{},
		&models.Contact{},
		&models.GroupBasic{},
	)
}
