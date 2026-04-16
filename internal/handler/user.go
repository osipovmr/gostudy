package handler

import (
	"net/http"

	"gostudy/internal/model/dto"
	"gostudy/internal/model/entity"
	"gostudy/internal/service"

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
		var input dto.CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &entity.User{
			ID:    uuid.New().String(),
			Name:  input.Name,
			Email: input.Email,
		}

		if err := svc.Create(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func listUsers(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := svc.List(c.Request.Context())
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

		user, err := svc.GetByID(c.Request.Context(), id)
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

		var input dto.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &entity.User{
			ID:    id,
			Name:  input.Name,
			Email: input.Email,
		}

		if err := svc.Update(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func deleteUser(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := svc.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
