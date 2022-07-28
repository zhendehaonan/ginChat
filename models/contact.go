package models

import (
	"gorm.io/gorm"
)

// 人员关系表
type Contact struct {
	gorm.Model
	OwnerId  int64 //谁的关系信息
	TargetId uint  //对应的谁
	Type     int   //对应的类型     1好友  2群组   3
	Desc     string
}

/*
// 查找所有好友
func SearchFriend(ctx context.Context, userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	userDB := dao.NewUserBasic(ctx)
	userDB.DB.Where("owner_id = ? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	userDB.DB.Where("id in ?", objIds).Find(&users)
	return users
}
*/
