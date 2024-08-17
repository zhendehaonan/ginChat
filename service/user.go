package service

import (
	"context"
	"fmt"
	"ginchat/dao"
	"ginchat/models"
	"ginchat/serializer"
	"ginchat/utils"
	"ginchat/utils/e"
	utils "ginchat/utils/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"mime/multipart"
	"net/http"
	"path"
	"time"
)

// 前端传过来的数据格式（json）
type UserService struct {
	UserName   string `json:"user_name" form:"user_name"`
	Password   string `json:"password" form:"password"`
	RePassword string `json:"re_password" form:"re_password"`
}

// 前端传过来的数据格式（json）  添加好友
type ContactService struct {
	TargetId uint `json:"target_id" form:"target_id"`
}

// 前端传过来的数据格式（json） 群组
type GroupService struct {
	GroupName string                `json:"group_name" form:"group_name"`
	Icon      *multipart.FileHeader `json:"icon" form:"icon"`
	Desc      string                `json:"desc" form:"desc"`
	Type      int                   `json:"type" form:"type"`
}

// 用户注册
func (userService *UserService) UserCreate(ctx context.Context) serializer.Response {
	var user models.UserBasic
	code := e.Success
	var err error
	//获取数据库操作对象
	userDao := dao.NewUserBasic(ctx)
	//判断用户名或密码是否为空
	if userService.UserName == "" || userService.Password == "" || userService.RePassword == "" {
		code = e.UserNameOrPasswordIsNull
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//判断用户名是否存在
	count := userDao.ExistOrNotByUserName(userService.UserName)
	if count != 0 {
		code = e.UserNameExist
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if len(userService.Password) > 16 || len(userService.Password) < 8 {
		code = e.PwdLenErr
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if userService.Password == userService.RePassword {
		//密码加密
		if err = user.SetPassword(userService.Password); err != nil {
			code = e.ErrorFailEncryption
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		user.UserName = userService.UserName
		//创建用户
		err = userDao.CreateUser(&user)
		if err != nil {
			code = e.Error
		}
	} else {
		code = e.PwdDiffer
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 用户登录
func (userService *UserService) Login(ctx context.Context) serializer.Response {
	var user *models.UserBasic
	code := e.Success
	var err error
	userDao := dao.NewUserBasic(ctx)
	//查询用户名是否存在
	count := userDao.ExistOrNotByUserName(userService.UserName)
	//不存在
	if count == 0 {
		code = e.UserNameNotExist
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//根据用户名查询密码
	user, err = userDao.GetPasswordByUserName(userService.UserName)
	//比对密码是否一致
	//不一致
	if user.CheckPassword(userService.Password) == false {
		code = e.PasswordErr
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//token发放（token整体流程 ::: token发放---token解析----token验证）
	token, err := util.GenerateToken(user.ID, user.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "token生成失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.TokenData{User: user, Token: token},
	}
}

// 通过用户id删除用户
func (userService *UserService) Delete(ctx context.Context, id uint) serializer.Response {
	code := e.Success
	var err error
	userDao := dao.NewUserBasic(ctx)
	//通过id删除
	err = userDao.DeleteUserById(id)
	//删除失败
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "删除失败",
			Error:  err,
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "删除成功",
	}
}

// 用户修改密码
func (userService *UserService) Update(ctx context.Context, id uint) serializer.Response {
	var user *models.UserBasic
	code := e.Success
	var err error
	userDao := dao.NewUserBasic(ctx)
	//通过id查找用户
	user, _ = userDao.GetUserById(id)
	if equal := user.CheckPassword(userService.Password); equal {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "密码与上一次密码一致,请重新输入",
		}
	}
	//密码加密
	if len(userService.Password) > 8 && len(userService.Password) < 16 {
		if err = user.SetPassword(userService.Password); err != nil {
			code = e.ErrorFailEncryption
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	//更新用户信息
	err = userDao.UpdateUserById(id, user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   user,
		Msg:    "修改成功",
	}
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 发送消息
func SendMsg(ctx *gin.Context) {
	ws, err := upGrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	MsgHandler(ws, ctx)
}
func MsgHandler(ws *websocket.Conn, ctx *gin.Context) {
	for {
		msg, err := utils.Subscribe(ctx, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		now := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", now, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 查找好友(自己写的 没按博主的写)
func SearchFriends(ctx context.Context, id uint) serializer.Response {
	objIds := make([]uint64, 0)
	contactDB := dao.NewContact(ctx)
	var err error
	code := e.Success
	friend, err := contactDB.GetFriendByOwnerId(id)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Error:  err,
		}
	}
	for _, v := range friend {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users, err := contactDB.GetFriendByTargetId(objIds)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Error:  err,
		}
	}
	return serializer.Response{
		Data:   users,
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 上传图片
func Upload(ctx *gin.Context) serializer.Response {
	file, _ := ctx.FormFile("file") // file为字段名
	//获取当前日期参数,当做文件名
	name := NowTimeString()
	extstring := path.Ext(file.Filename)
	dst := "./upload/" + name + extstring
	//直接用原来的文件名
	//dst := "upload_files/" + file.Filename
	ctx.SaveUploadedFile(file, dst)

	// 都是gin.context作为入参
	return serializer.Response{
		Status: e.Success,
		Msg:    e.GetMsg(e.Success),
	}
}

// 获取当前日期的参数当做文件名
func NowTimeString() string {
	now := time.Now()
	str := fmt.Sprintf("%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond())
	return str
}

// 添加好友
func (contactService *ContactService) AddFriend(ctx context.Context, id uint) serializer.Response {
	code := e.Success
	var err error
	var contactUser models.Contact
	contactDao := dao.NewContact(ctx)
	if id == contactService.TargetId {
		code = e.AddFriendFail
		return serializer.Response{
			Status: code,
			Msg:    "不能添加自己为好友",
		}
	}
	//查询是否已经是好友
	bool := contactDao.IsFriend(id, contactService.TargetId)
	if bool {
		//已经是好友
		code = e.IsFriend
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//添加好友
	contactUser.OwnerId = int64(id)
	contactUser.TargetId = contactService.TargetId
	contactUser.Type = 1
	err = contactDao.AddFriend(contactUser)
	if err != nil {
		code = e.AddFriendFail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err,
		}
	}
	//查询所有好友信息
	friend, _ := contactDao.GetFriendByOwnerId(id)
	return serializer.Response{
		Status: code,
		Msg:    "添加好友成功",
		Data:   friend,
	}
}

// 新建群聊
func (groupService *GroupService) CreateGroup(ctx context.Context, id uint) serializer.Response {
	var user *models.GroupBasic
	code := e.Success
	var err error
	userDao := dao.NewGroupBasic(ctx)
	//判断群名是否存在
	bool := userDao.IsGroupName(groupService.GroupName)
	//群名已存在
	if bool {
		return serializer.Response{
			Status: e.Error,
			Msg:    "群名已存在",
		}
	}
	//判断群名是否为空
	if groupService.GroupName == "" {
		return serializer.Response{
			Status: e.Error,
			Msg:    "群名不能为空",
		}
	}
	if len(groupService.GroupName) > 10 {
		return serializer.Response{
			Status: e.Error,
			Msg:    "群名过长",
		}
	}
	//新建群
	user = &models.GroupBasic{
		Name:    groupService.GroupName,
		OwnerId: id,
		Desc:    groupService.Desc,
		Type:    groupService.Type,
	}
	//处理图片
	name := NowTimeString()
	extstring := path.Ext(groupService.Icon.Filename)
	dst := name + extstring
	g := gin.Context{}
	g.SaveUploadedFile(groupService.Icon, "./upload/"+dst)
	user.Icon = dst
	err = userDao.CreateGroup(user)
	if err != nil {
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    "创建失败",
			Error:  err,
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "创建成功",
	}
}
