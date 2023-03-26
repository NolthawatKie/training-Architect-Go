package todo

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text" binding:"required"`
	gorm.Model
}

// to create storer interface
type storer interface {
	New(*Todo) error
}

func (Todo) TableName() string {
	return "todos"
}

type TodoHandler struct {
	//db *gorm.DB
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

// create type to change behavior
type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	TransactionID() string
	Audience() string
}

func (t *TodoHandler) NewTask(c Context) {
	var todo Todo
	// if err := c.ShouldBindJSON(&todo); err != nil {
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		// transactionID := c.Request.Header.Get("TransactionID")
		transactionID := c.TransactionID()
		// aud, _ := c.Get("aud")
		aud := c.Audience()
		log.Println(transactionID, aud, "not allowed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not allowed",
		})
		return
	}

	//r := t.db.Create(&todo)
	err := t.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}

// func (t *TodoHandler) List(c *gin.Context) {
// 	var todos []Todo
// 	r := t.db.Find(&todos)
// 	if err := r.Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, todos)
// }

// func (t *TodoHandler) Remove(c *gin.Context) {
// 	idParam := c.Param("id")

// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	r := t.db.Delete(&Todo{}, id)
// 	if err := r.Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status": "success",
// 	})
// }
