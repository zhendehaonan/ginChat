package serializer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 返回给前端的数据格式
type Response struct {
	Status int    `json:"status"` //状态码
	Data   any    `json:"data"`   //数据(一般是成功数据)
	Msg    string `json:"msg"`    //信息(一般返回的是错误信息)
	Error  error  `json:"err"`
}

type TokenData struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

// 好友列表返回（以下全是）
type H struct {
	Code  int
	Msg   string
	Data  interface{}
	Rows  interface{}
	Total interface{}
}

func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
