package router

import (
	"ginchat/api/v1"
	"ginchat/docs"
	"ginchat/middleware"
	"ginchat/models"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"html/template"
	"net/http"
	"strconv"
)

//添加swagger，验证swagger案例，。下面四个注解是配置swagger（1.方法名  3.返回状态码和返回值 4.请求路径和请求方法）

// HelloWorld
// @Tags swagger案例
// @Success 200 {string} helloworld
// @Router /helloworld [get]
func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "helloworld")
}

//首页

// Index
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func Index(ctx *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(ctx.Writer, "index")
}

// ToRegister
// @Tags 注册页面
// @Success 200 {string} ToRegister
// @Router /toRegister [get]
func ToRegister(ctx *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(ctx.Writer, "index")
}

// ToChat
// @Tags 访问主页面（登录成功后的页面）
// @Success 200 {string} ToChat
// @Router /toChat [get]
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	//fmt.Println("ToChat>>>>>>>>", user)
	ind.Execute(c.Writer, user)
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	//跨域
	r.Use(middleware.Cors())
	//加载静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")
	//添加swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//验证swagger案例
	r.POST("/helloworld", HelloWorld)
	//访问首页(登陆页面)
	r.GET("/index", Index)
	//访问注册页面
	r.GET("/toRegister", ToRegister)
	//访问主页面（登录成功后的页面）
	r.GET("/toChat", ToChat)
	//用户操作（user）路由模块
	v1 := r.Group("/v1/user")
	{
		//创建用户
		v1.POST("/createUser", api.CreateUser)
		//登录
		v1.POST("/login", api.Login)
		//后续请求加入登陆保护
		v1.Use(middleware.JWT())
		//删除用户
		v1.DELETE("/delete", api.Delete)
		//用户修改密码
		v1.PUT("/update", api.Update)
		//查找好友
		v1.POST("/searchFriend", api.SearchFriend)
		//上传图片
		v1.POST("/upload", api.Upload)
		//添加好友
		v1.POST("/addFriend", api.AddFriend)
		//新建群组
		v1.POST("/createGroup", api.CreateGroup)
	}
	//发送消息
	r.GET("/sendMsg", service.SendMsg)
	r.GET("/sendUserMsg", api.SendUserMsg)
	r.Run(":8080")
	return r
}
