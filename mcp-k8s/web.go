package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"mcp-k8s/Service/mcp"
	"mcp-k8s/Util"
	"time"
)

// 1. 配置常量
const (
	JWTSecret = "your-secret-key" // JWT 密钥（生产环境需保密）
	DBDSN     = "root:123456@tcp(127.0.0.1:3306)/gin_user?charset=utf8mb4&parseTime=True&loc=Local"
)

// 2. 模型定义
type User struct {
	gorm.Model
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:100;not null" json:"-"` // json:"-" 避免返回密码
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
}

// 3. 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 4. JWT 工具函数
func GenerateToken(username string) (string, error) {
	// 过期时间：2 小时
	expireTime := time.Now().Add(2 * time.Hour)
	// 构造 claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expireTime.Unix(),
	}
	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// JWT 中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, Response{Code: 401, Message: "未登录"})
			return
		}

		// 解析 token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("未知签名方法: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, Response{Code: 401, Message: "token 无效或过期"})
			return
		}

		// 传递 username 到上下文
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("username", claims["username"].(string))
		c.Next()
	}
}

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello"})
}

// 5. 数据库初始化
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open(DBDSN), &gorm.Config{})
//	if err != nil {
//		panic("数据库连接失败: " + err.Error())
//	}
//	db.AutoMigrate(&User{})
//	return db
//}

var db *gorm.DB

func main() {
	//db = initDB()
	r := gin.Default()
	r.GET("/api/hello", hello)
	r.GET("/api/endpoints/:path", func(c *gin.Context) {
		Util.Output(mcp.McpGetEndpoints(c.Query("path")), "endpoint.xlsx")
		c.JSON(200, gin.H{"message": "success"})
	})
	r.GET("/api/images/:path", func(c *gin.Context) {
		Util.Output(mcp.McpGetImages(c.Query("path")), "image.xlsx")
		c.JSON(200, gin.H{"message": "success"})
	})

	// 6. 公开路由（无需登录）
	//public := r.Group("/api/public")
	//{
	//	// 注册
	//	public.POST("/register", func(c *gin.Context) {
	//		var req struct {
	//			Username string `json:"username" binding:"required"`
	//			Password string `json:"password" binding:"required,min=6"`
	//			Email    string `json:"email" binding:"required,email"`
	//		}
	//		if err := c.ShouldBindJSON(&req); err != nil {
	//			c.JSON(400, Response{Code: 400, Message: err.Error()})
	//			return
	//		}
	//
	//		// 检查用户名是否已存在
	//		var user User
	//		if err := db.Where("username = ?", req.Username).First(&user).Error; err == nil {
	//			c.JSON(400, Response{Code: 400, Message: "用户名已存在"})
	//			return
	//		}
	//
	//		// 创建用户（实际开发中需加密密码，如 bcrypt）
	//		newUser := User{
	//			Username: req.Username,
	//			Password: req.Password, // 注意：生产环境需加密！
	//			Email:    req.Email,
	//		}
	//		if err := db.Create(&newUser).Error; err != nil {
	//			c.JSON(500, Response{Code: 500, Message: "注册失败"})
	//			return
	//		}
	//
	//		c.JSON(200, Response{Code: 0, Message: "注册成功", Data: newUser})
	//	})
	//
	//	// 登录
	//	public.POST("/login", func(c *gin.Context) {
	//		var req struct {
	//			Username string `json:"username" binding:"required"`
	//			Password string `json:"password" binding:"required"`
	//		}
	//		if err := c.ShouldBindJSON(&req); err != nil {
	//			c.JSON(400, Response{Code: 400, Message: err.Error()})
	//			return
	//		}
	//
	//		// 验证用户
	//		var user User
	//		if err := db.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
	//			c.JSON(401, Response{Code: 401, Message: "用户名或密码错误"})
	//			return
	//		}
	//
	//		// 生成 JWT token
	//		token, err := GenerateToken(req.Username)
	//		if err != nil {
	//			c.JSON(500, Response{Code: 500, Message: "生成 token 失败"})
	//			return
	//		}
	//
	//		c.JSON(200, Response{Code: 0, Message: "登录成功", Data: gin.H{"token": token}})
	//	})
	//
	//}

	//// 7. 需登录的路由（JWT 认证）
	//auth := r.Group("/api/auth")
	//auth.Use(JWTMiddleware()) // 添加 JWT 中间件
	//{
	//	// 获取当前用户信息
	//	auth.GET("/profile", func(c *gin.Context) {
	//		username, _ := c.Get("username")
	//		var user User
	//		db.Where("username = ?", username).First(&user)
	//		c.JSON(200, Response{Code: 0, Message: "success", Data: user})
	//	})
	//
	//	// 更新用户信息
	//	auth.PUT("/profile", func(c *gin.Context) {
	//		username, _ := c.Get("username")
	//		var req struct {
	//			Email string `json:"email" binding:"email"`
	//			Age   int    `json:"age"`
	//		}
	//		if err := c.ShouldBindJSON(&req); err != nil {
	//			c.JSON(400, Response{Code: 400, Message: err.Error()})
	//			return
	//		}
	//
	//		db.Model(&User{}).Where("username = ?", username).Updates(req)
	//		c.JSON(200, Response{Code: 0, Message: "更新成功"})
	//	})
	//}

	r.Run(":8080")
}
