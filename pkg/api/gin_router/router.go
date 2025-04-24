package router

import (
	"context"

	"github.com/gin-gonic/gin"
	customer "github.com/nester817/session_authentication_api.git/pkg/storage"
)

type DB interface {
	GetCustomerByEmail(ctx context.Context, email, password string) (*customer.Customer, error)
	InsertCustomer(ctx context.Context, user customer.Customer) error
	DeleteCustomer(ctx context.Context, email, password string) error
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
	profile := r.Group("/profile")
	{
		profile.POST("/", h.AddProfile)
		profile.DELETE("/", h.DeleteProfile)
	}

	return r
}
