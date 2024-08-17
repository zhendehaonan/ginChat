package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// 用户表
type UserBasic struct {
	gorm.Model
	UserName      string
	Password      string
	Phone         string
	email         string
	Identity      string //   唯一标识
	ClientIp      string
	ClientPort    string
	LoginTime     *time.Time
	HeartbeatTime *time.Time
	LoginOutTime  *time.Time
	IsLogout      bool
	DeviceInfo    string
}

const (
	PasswordCost        = 12       //密码加密难度
	Active       string = "active" //激活用户
)

// 密码加密
func (user *UserBasic) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// 比对密码是否一致
func (user UserBasic) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
