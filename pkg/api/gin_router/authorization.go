package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	customer "github.com/nester817/session_authentication_api.git/pkg/storage"
)

func (h *Handler) Login(ctx *gin.Context) {
	var conf customer.Customer
	if err := ctx.ShouldBindJSON(&conf); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(conf.Password) < 6 {
		ctx.IndentedJSON(http.StatusBadRequest, "incorrect password")
		return
	}

	customer, err := h.sqlDb.GetCustomer(
		ctx.Request.Context(),
		conf.Email,
	)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	flag := ComparePasswords([]byte(customer.Password), []byte(conf.Password))
	if !flag {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "incorrect password",
		})
		return
	}

	sessionId := uuid.New().String()
	if err := h.keyValueDb.Set(ctx.Request.Context(), sessionId, fmt.Sprint(customer.Id)); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("session_id", sessionId, 3600, "/", "localhost", false, true)
	ctx.IndentedJSON(http.StatusOK, "login to profile successful")
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	sessionValue, err := h.getSessionId(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	customer, err := h.sqlDb.GetCustomerById(ctx.Request.Context(), sessionValue)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer.Email)
}

func (h *Handler) Logout(ctx *gin.Context) {
	sessionValue, err := h.getSessionId(ctx)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.keyValueDb.Delete(ctx.Request.Context(), sessionValue); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("session_id", "", -1, "/", "localhost", true, true)

	ctx.IndentedJSON(http.StatusOK, "profile deleted successfully")
}

func (h *Handler) getSessionId(ctx *gin.Context) (string, error) {
	cookieValue, err := ctx.Cookie("session_id")
	if err != nil {
		return "", err
	}

	sessionValue, err := h.keyValueDb.Get(ctx.Request.Context(), cookieValue)
	if err != nil {
		return "", err
	}

	return sessionValue, nil
}
