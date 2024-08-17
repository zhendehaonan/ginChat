package dao

import (
	"context"
	"ginchat/config"
	"ginchat/models"
	"gorm.io/gorm"
)

type UserBasic struct {
	*gorm.DB
}

// 创建该实体类的数据库对象
func NewUserBasic(ctx context.Context) *UserBasic {
	return &UserBasic{config.NewDBClient(ctx)}
}

// 判断用户名是否存在
func (dao *UserBasic) ExistOrNotByUserName(username string) int {
	var count int64
	dao.DB.Model(&models.UserBasic{}).Where("user_name = ?", username).Count(&count)
	return int(count)
}

// 创建用户
func (dao *UserBasic) CreateUser(user *models.UserBasic) error {
	err := dao.DB.Model(&models.UserBasic{}).Create(&user).Error
	return err
}

// 根据用户名查询密码
func (dao *UserBasic) GetPasswordByUserName(username string) (user *models.UserBasic, err error) {
	err = dao.DB.Model(&models.UserBasic{}).Where("user_name = ?", username).First(&user).Error
	return
}

// 根据id删除用户
func (dao *UserBasic) DeleteUserById(id uint) error {
	return dao.DB.Model(&models.UserBasic{}).Where("id = ?", id).Delete(&models.UserBasic{}).Error
}

// 根据id查找用户
func (dao *UserBasic) GetUserById(id uint) (user *models.UserBasic, err error) {
	dao.DB.Model(&models.UserBasic{}).Where("id = ?", id).First(&user)
	return
}

// 根据id修改信息
func (userDao *UserBasic) UpdateUserById(id uint, user *models.UserBasic) error {
	return userDao.DB.Model(&models.UserBasic{}).Where("id=?", id).Updates(&user).Error
}
