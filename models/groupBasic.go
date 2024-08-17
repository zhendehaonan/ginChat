package models

import "gorm.io/gorm"

// 群组表
type GroupBasic struct {
	gorm.Model
	Name    string //群名称
	OwnerId uint   //群主
	Icon    string //群头像
	Desc    string //群描述
	Type    int    //群类型
}
