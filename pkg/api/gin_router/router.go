package router

import (
	"context"

	"github.com/gin-gonic/gin"
	customer "github.com/nester817/session_authentication_api.git/pkg/storage"
	"golang.org/x/crypto/bcrypt"
)

type SqlDB interface {
	GetCustomer(ctx context.Context, email string) (*customer.Customer, error)
	GetCustomerById(ctx context.Context, id string) (*customer.Customer, error)
	InsertCustomer(ctx context.Context, user *customer.Customer) error
}

type KeyValueDb interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
}

type Handler struct {
	sqlDb      SqlDB
	keyValueDb KeyValueDb
}

func NewHandler(db SqlDB, client KeyValueDb) *Handler {
	return &Handler{
		sqlDb:      db,
		keyValueDb: client,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	profile := r.Group("/profile")
	{
		profile.POST("/registr", h.AddProfile)

		authorization := profile.Group("/authentication")
		{
			authorization.GET("/home", h.GetProfile)
			authorization.POST("/login", h.Login)
			authorization.DELETE("/logout", h.Logout)
		}
	}

	return r
}

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func ComparePasswords(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
