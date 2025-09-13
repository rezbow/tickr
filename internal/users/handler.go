package users

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
)

func (service *UsersService) CreateUserHandler(c *gin.Context) {
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors := userInput.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	hash, err := hashPassword(userInput.Password)
	if err != nil {
		service.logger.Error("failed to hash password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &entities.User{
		Name:         userInput.Name,
		Email:        userInput.Email,
		Role:         userInput.Role,
		PasswordHash: hash,
	}

	if err := service.createUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (service *UsersService) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.deleteUser(c.Request.Context(), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			service.logger.Error("failed to delete user", "userId", userID.String(), "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	service.logger.Info("user deleted", "userId", userID.String())
	c.JSON(http.StatusNoContent, nil)
}

func (service *UsersService) GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.getUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			service.logger.Error("failed to get user", "userId", userID.String(), "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// getUsers: get all users with pagination
func (service *UsersService) GetUsersHandler(c *gin.Context) {
	var (
		page  int
		limit int
		err   error
	)

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err = strconv.Atoi(pageStr)
	if page <= 0 || err != nil {
		page = 1
	}
	limit, err = strconv.Atoi(limitStr)
	if limit <= 0 || err != nil {
		limit = 10
	}

	users, total, err := service.getUsers(c.Request.Context(), page, limit)
	if err != nil {
		service.logger.Error("failed to get users", "page", page, "limit", limit, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response struct {
		Data      []entities.User `json:"data"`
		Page      int             `json:"page"`
		Limit     int             `json:"limit"`
		Total     int64           `json:"total"`
		TotalPage int             `json:"total_page"`
	}
	response.Data = users
	response.Page = page
	response.Limit = limit
	response.Total = total
	response.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, response)
}

func (service *UsersService) UpdateUserHander(c *gin.Context) {
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedFields UserUpdateInput
	if err := c.ShouldBindJSON(&updatedFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := updatedFields.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	updates, err := updatedFields.ToMap()
	if err != nil {
		service.logger.Error("failed to update user", "userId", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user, err := service.updateUserAtomic(userId, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		service.logger.Error("failed to update user", "userId", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
