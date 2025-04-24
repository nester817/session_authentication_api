package router

import "github.com/gin-gonic/gin"

type DB interface {
}

type Handler struct {
	db DB
}

func NewHandler(db DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}
