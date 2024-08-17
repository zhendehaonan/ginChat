package dao

import (
	"context"
	"ginchat/config"
	"ginchat/models"
	"gorm.io/gorm"
	"log"
)

type Contact struct {
	*gorm.DB
}

// 创建该实体类的数据库对象
func NewContact(ctx context.Context) *Contact {
	return &Contact{config.NewDBClient(ctx)}
}

// 查找好友(自己写的 没按博主的写)
// 根据owner_id查找所关联的好友
func (dao *Contact) GetFriendByOwnerId(id uint) ([]*models.Contact, error) {
	contact := make([]*models.Contact, 0)
	err := dao.DB.Model(&models.Contact{}).Where("owner_id = ? and type=1", id).Find(&contact).Error
	return contact, err
}

// 根据target_id查找所关联的好友
func (dao *Contact) GetFriendByTargetId(objIds []uint64) ([]models.UserBasic, error) {
	userBasic := make([]models.UserBasic, 0)
	err := dao.DB.Model(&models.UserBasic{}).Where("id in ?", objIds).Find(&userBasic).Error
	return userBasic, err
}

// 查询是否已经是好友
func (dao *Contact) IsFriend(id uint, targetId uint) bool {
	var count int64
	err := dao.DB.Model(&models.Contact{}).Where("owner_id = ? and target_id = ? and type = 1 ", id, targetId).Count(&count).Error
	if err != nil {
		log.Printf("Error while checking friendship: %v", err)
		return false
	}
	return count > 0
}

// 添加好友
func (dao *Contact) AddFriend(user models.Contact) error {
	return dao.DB.Model(&models.Contact{}).Create(&user).Error
}
