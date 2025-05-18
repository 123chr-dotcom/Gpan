package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"gpan/models"
	"gpan/utils"
)

// UploadFile 文件上传
func UploadFile(c *gin.Context) {
	// 获取上传文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}

	// 保存文件
	var savedFiles []models.File
	for _, file := range files {
		// 生成唯一文件名
		ext := filepath.Ext(file.Filename)
		filename := utils.GenerateUUID() + ext
		filePath := filepath.Join("uploads", filename)

		// 保存到本地
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
			return
		}

		// 保存到数据库
		// 默认上传到根目录(0)，实际应用中应从请求中获取目录ID
		dbFile := models.File{
			Name:        file.Filename,
			Path:        filePath,
			Size:        file.Size,
			Extension:   ext,
			UserID:      1, // 临时使用用户ID=1，实际应从认证信息中获取
			DirectoryID: 0, // 根目录
		}
		if err := models.CreateFile(&dbFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件信息保存失败"})
			return
		}

		savedFiles = append(savedFiles, dbFile)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"files":   savedFiles,
	})
}

// FileResponse 定义API响应格式
type FileResponse struct {
	ID        uint   `json:"ID"`
	Name      string `json:"Name"`
	Size      int64  `json:"Size"`
	UpdatedAt string `json:"UpdatedAt"`
}

// GetFiles 获取文件列表
func GetFiles(c *gin.Context) {
	files, err := models.GetAllFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件列表失败"})
		return
	}

	// 转换为前端期望的格式
	response := make([]FileResponse, 0) // 确保返回空数组而不是nil
	for _, file := range files {
		response = append(response, FileResponse{
			ID:        file.ID,
			Name:      file.Name,
			Size:      file.Size,
			UpdatedAt: file.UpdatedAt.Format("2006/1/2 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, response)
}

// DownloadFile 文件下载
func DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文件ID",
			"details": gin.H{
				"received_id": idStr,
				"error": err.Error(),
			},
		})
		return
	}

	file, err := models.GetFileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.FileAttachment(file.Path, file.Name)
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的文件ID",
			"details": gin.H{
				"received_id": c.Param("id"),
				"error": err.Error(),
			},
		})
		return
	}

	file, err := models.GetFileByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"message": "文件记录不存在或已被删除"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询文件失败",
			"details": err.Error(),
		})
		return
	}

	// 删除物理文件
	if err := utils.DeleteFile(file.Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "物理文件删除失败",
			"details": err.Error(),
		})
		return
	}

	// 删除数据库记录
	if err := models.DeleteFile(uint(id)); err != nil {
		// 记录详细错误日志
		fmt.Printf("删除文件记录失败(ID:%d): %v\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库记录删除失败",
			"details": err.Error(),
			"file_id": id,
		})
		return
	}

	// 记录成功日志
	fmt.Printf("成功删除文件(ID:%d): %s\n", file.ID, file.Name)

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
		"file": gin.H{
			"id": file.ID,
			"name": file.Name,
		},
	})
}

// SetupFileRoutes 设置文件路由
func SetupFileRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/upload", UploadFile)
		api.GET("/files", GetFiles)
		api.GET("/download/:id", DownloadFile)
		api.DELETE("/files/:id", DeleteFile)
	}
}
