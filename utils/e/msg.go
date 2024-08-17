package e

var MsgFlags = map[int]string{
	Success:                    "ok",
	Error:                      "fail",
	PwdLenErr:                  "密码长度有误",
	UserNameExist:              "用户名已存在",
	PwdDiffer:                  "密码不一致",
	UserNameNotExist:           "用户名不存在",
	PasswordErr:                "密码错误,请重新输入",
	ErrorAuthToken:             "token认证失败",
	ErrorAuthCheckTokenTimeOut: "token过期",
	ErrorFailEncryption:        "密码加密失败",
	UserNameOrPasswordIsNull:   "用户名或密码为空",
	IsFriend:                   "已经是好友",
	AddFriendFail:              "添加好友失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
