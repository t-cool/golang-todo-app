package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Tasks    []Task
}

type Task struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"` // 1: 低, 2: 中, 3: 高
	Done        bool      `json:"done"`
	UserID      uint      `json:"user_id"`
}

var db *gorm.DB

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// データベースの初期化
	var err error
	db, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("データベース接続エラー:", err)
	}

	// マイグレーション
	db.AutoMigrate(&User{}, &Task{})

	// Ginルータの設定
	r := gin.Default()

	// 静的ファイルの提供
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// ルートページの提供
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 認証不要のルート
	r.POST("/api/register", register)
	r.POST("/api/login", login)

	// 認証が必要なルート
	auth := r.Group("/api")
	auth.Use(authMiddleware())
	{
		auth.GET("/tasks", getTasks)
		auth.POST("/tasks", createTask)
		auth.PUT("/tasks/:id", updateTask)
		auth.DELETE("/tasks/:id", deleteTask)
	}

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザー登録に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ユーザー登録が完了しました"})
}

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser User
	if err := db.Where("username = ? AND password = ?", user.Username, user.Password).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証に失敗しました"})
		return
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}

func getTasks(c *gin.Context) {
	userID := c.GetUint("user_id")
	var tasks []Task
	if err := db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UserID = c.GetUint("user_id")
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("id")

	var task Task
	if err := db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "タスクが見つかりません"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&task)
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	userID := c.GetUint("user_id")
	taskID := c.Param("id")

	var task Task
	if err := db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "タスクが見つかりません"})
		return
	}

	db.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "タスクを削除しました"})
}
