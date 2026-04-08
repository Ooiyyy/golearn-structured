package handler

import (
	"fmt"
	"golearn-structured/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandlers(s *service.TodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

func (h *TodoHandler) Create(c *gin.Context) {
	// Request create todo diterima di handler ini.
	// with jwt
	// userID := c.GetInt("id")
	// without jwt
	userID := 5

	// Input dibaca dari multipart/form-data.
	title := c.PostForm("title")
	note := c.PostForm("note")
	file, err := c.FormFile("image")
	imageURL := ""
	uploadPath := ""

	if err == nil {
		// Simpan file upload ke folder lokal, lalu bentuk URL publiknya.
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		uploadPath = "uploads/" + filename
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal menyimpan gambar",
			})
			return
		}
		imageURL = "http://localhost:8080/" + uploadPath
	}

	err = h.service.CreateTodo(userID, title, note, imageURL)
	if err != nil {
		// Rollback file jika validasi bisnis gagal setelah upload.
		if uploadPath != "" {
			os.Remove(uploadPath)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "todo berhasil dibuat",
	})
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	// Handler membaca query param, lalu delegasi ambil data ke service.
	// with jwt
	// userID := c.GetInt("id")

	// tanpa jwt
	userStr := c.Query("userid")
	userID, _ := strconv.Atoi(userStr)
	todos, err := h.service.GetUserTodos(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data todo",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		// Response akhir selalu dikirim dari handler dalam bentuk JSON.
		"data": todos,
	})
}

// func (h *TodoHandler) Update(c *gin.Context) {
// 	// (Ganti c.Query jadi c.GetInt("id") kalau pakai JWT lagi nanti)
// 	userStr := c.Query("userid")
// 	userID, _ := strconv.Atoi(userStr)

// 	todoIDstr := c.Param("id")
// 	todoID, _ := strconv.Atoi(todoIDstr)

// 	var req struct {
// 		Title    string `json:"title"`
// 		Note     string `json:"note"`
// 		ImageUrl string `json:"image_url"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	err := h.service.EditTodo(todoID, userID, req.Title, req.Note, req.ImageUrl)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate todo"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "todo berhasil diupdate"})
// }

func (h *TodoHandler) Update(c *gin.Context) {
	// Ambil identitas user + todo target dari request URL/query.
	userStr := c.Query("userid")
	userID, _ := strconv.Atoi(userStr)

	todoIDStr := c.Param("id")
	todoID, _ := strconv.Atoi(todoIDStr)

	title := c.PostForm("title")
	note := c.PostForm("note")

	file, err := c.FormFile("image")
	imageURL := ""
	uploadPath := ""

	if err == nil {
		// Jika ada file baru, upload dulu sebelum business update diproses.
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		uploadPath = "uploads/" + filename
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
			return
		}
		imageURL = "http://localhost:8080" + uploadPath
	}
	err = h.service.EditTodo(todoID, userID, title, note, imageURL)

	if err != nil {
		// Hapus file baru jika update gagal agar storage tidak kotor.
		if uploadPath != "" {
			os.Remove(uploadPath)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "todo berhasil diupdate"})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// Endpoint delete: parse param, panggil service, lalu kirim status response.
	userStr := c.Query("userid")
	userID, _ := strconv.Atoi(userStr)

	todoIDstr := c.Param("id")
	todoID, _ := strconv.Atoi(todoIDstr)

	err := h.service.DeleteTodo(todoID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "todo berhasil dihapus"})
}
