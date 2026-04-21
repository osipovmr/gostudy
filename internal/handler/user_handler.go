package handler

import (
	"errors"
	"net/http"

	"gostudy/internal/model/dto"
	"gostudy/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(routerGroup *gin.RouterGroup) {
	users := routerGroup.Group("/users")
	{
		users.POST("", h.createUser)
		users.GET("", h.listUsers)
		users.GET("/:id", h.getUser)
		users.PUT("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
	}
}

func (h *UserHandler) createUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid request body"))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			c.JSON(http.StatusConflict, errorResponse("user already exists"))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("internal error"))
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) listUsers(c *gin.Context) {
	users, err := h.userService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("internal error"))
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) getUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user not found"))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("internal error"))
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) updateUser(c *gin.Context) {
	id := c.Param("id")

	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid request body"))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user not found"))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("internal error"))
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user not found"))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("internal error"))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func errorResponse(msg string) gin.H {
	return gin.H{"error": msg}
}
