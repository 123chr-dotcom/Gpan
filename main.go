package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"gpan/config"
	"gpan/database"
	"gpan/models"
	"gpan/routes"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库 - 优先尝试SQLite
	db, err := database.NewSQLiteDB()
	if err != nil {
		log.Printf("Failed to connect to SQLite, trying PostgreSQL: %v", err)
		// 回退到PostgreSQL
		db, err = database.NewPostgresDB(&cfg.Postgres)
		if err != nil {
			log.Fatalf("Failed to connect to any database: %v", err)
		}
	}
	
	// 确保全局DB变量被正确设置
	database.DB = db
	if database.DB == nil {
		log.Fatal("Database connection is nil after initialization")
	}
	
	// 测试数据库连接
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 数据库迁移
	models := []interface{}{
		&models.User{},
		&models.File{},
		&models.Directory{},
		&models.FileChunk{},
	}
	if err := database.AutoMigrate(db, models...); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化Gin引擎
	r := gin.Default()
	
	// 配置静态文件服务
	r.Static("/static", "./frontend/static")
	r.Static("/uploads", "./uploads")
	
	// 加载HTML模板
	r.LoadHTMLGlob("frontend/templates/*")
	
	// 注册文件路由
	routes.SetupFileRoutes(r)
	
	// 基本路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "GPAN网盘系统",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器
	log.Printf("Server starting on :%s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
