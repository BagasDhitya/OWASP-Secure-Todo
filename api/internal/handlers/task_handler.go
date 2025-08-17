package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/BagasDhitya/owasp-secure-todo/internal/models"
	"github.com/BagasDhitya/owasp-secure-todo/internal/repo"
	"github.com/BagasDhitya/owasp-secure-todo/internal/validators"
)

type TaskHandler struct {
	Log   *zap.Logger
	Tasks *repo.TaskRepo
}

func (h *TaskHandler) List(c *gin.Context) {
	uid := c.GetInt64("userID")
	items, err := h.Tasks.ListByUser(c, uid)
	if err != nil {
		h.Log.Error("list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *TaskHandler) Create(c *gin.Context) {
	uid := c.GetInt64("userID")
	var dto validators.TaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil || validators.V.Struct(dto) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	t := &models.Task{Title: dto.Title, Description: dto.Description, Status: dto.Status}
	if err := h.Tasks.Create(c, uid, t); err != nil {
		h.Log.Error("create", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TaskHandler) Update(c *gin.Context) {
	uid := c.GetInt64("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var dto validators.TaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil || validators.V.Struct(dto) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	t := &models.Task{Title: dto.Title, Description: dto.Description, Status: dto.Status}
	if err := h.Tasks.Update(c, uid, id, t); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	uid := c.GetInt64("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Tasks.Delete(c, uid, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
