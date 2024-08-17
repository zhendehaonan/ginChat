package e

const (
	Success = 200
	Error   = 500
	//user模块的错误
	PwdLenErr                  = 10001 //密码长度有误
	UserNameExist              = 10002 //用户名已存在
	PwdDiffer                  = 10003 //密码和重新输入的密码不一致
	UserNameNotExist           = 10004 //用户名不存在
	PasswordErr                = 10005 //密码错误,请重新输入
	ErrorAuthToken             = 10006 //token认证失败
	ErrorAuthCheckTokenTimeOut = 10007 //token过期
	ErrorFailEncryption        = 10008 //密码加密失败
	UserNameOrPasswordIsNull   = 10009 //用户名或密码为空
	IsFriend                   = 10010 //已经是好友
	AddFriendFail              = 10011 //添加好友失败
)
