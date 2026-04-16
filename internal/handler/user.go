package handler

import (
	"gostudy/internal/service"
	"net/http"

	"gostudy/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUserRoutes(router *gin.Engine, userSvc service.UserService) {
	users := router.Group("/api/v1/users")
	{
		users.POST("", createUser(userSvc))
		users.GET("", listUsers(userSvc))
		users.GET("/:id", getUser(userSvc))
		users.PUT("/:id", updateUser(userSvc))
		users.DELETE("/:id", deleteUser(userSvc))
	}
}

func createUser(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "body is required"})
			return
		}
		var input model.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &model.User{
			ID:    uuid.New().String(),
			Name:  input.Name,
			Email: input.Email,
		}
		if err := svc.Create(user); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func listUsers(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := svc.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func getUser(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := svc.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func updateUser(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var input model.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &model.User{
			Name:  input.Name,
			Email: input.Email,
		}
		if err := svc.Update(id, user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func deleteUser(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := svc.Delete(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
