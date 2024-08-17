package dao

import (
	"context"
	"ginchat/config"
	"ginchat/models"
	"gorm.io/gorm"
	"log"
)

type GroupBasic struct {
	*gorm.DB
}

// 创建该实体类的数据库对象
func NewGroupBasic(ctx context.Context) *GroupBasic {
	return &GroupBasic{config.NewDBClient(ctx)}
}

// 判断群名是否存在
func (dao *GroupBasic) IsGroupName(name string) bool {
	var count int64
	err := dao.DB.Model(&models.GroupBasic{}).Where("name = ? and type = 3 ", name).Count(&count).Error
	if err != nil {
		log.Printf("Error while checking group: %v", err)
		return false
	}
	return count > 0
}

// 新建群
func (dao *GroupBasic) CreateGroup(group *models.GroupBasic) error {
	return dao.DB.Model(&models.GroupBasic{}).Create(&group).Error

}
