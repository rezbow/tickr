package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (service *UsersService) CreateUserHandler(c *gin.Context) {
	var userInput UserCreateDTO
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
	c.JSON(http.StatusCreated, UserEntityToUserResponse(user))
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
	c.Status(http.StatusNoContent)
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
	c.JSON(http.StatusOK, UserEntityToUserResponse(user))
}

// getUsers: get all users with pagination
func (service *UsersService) GetUsersHandler(c *gin.Context) {
	var p utils.Pagination

	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})
		return
	}

	users, total, err := service.getUsers(c.Request.Context(), &p)
	if err != nil {
		service.logger.Error("failed to get users", "page", p.Page, "page_size", p.PageSize, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":      UserEntitiesToUserResponse(users),
		"total":     total,
		"page":      p.Page,
		"page_size": p.PageSize,
	})
}

func (service *UsersService) UpdateUserHander(c *gin.Context) {
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedFields UserUpdateDTO
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
	c.JSON(http.StatusOK, UserEntityToUserResponse(user))
}

func (service *UsersService) LoginHandler(c *gin.Context) {
	var loginInput LoginDTO
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := loginInput.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Find user by email
	var user entities.User
	if err := service.db.Where("email = ?", loginInput.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			service.logger.Error("failed to find user", "email", loginInput.Email, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT access token
	accessToken, err := service.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		service.logger.Error("failed to generate access token", "userId", user.ID.String(), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate refresh token
	refreshTokenString, err := service.jwtService.GenerateRefreshToken()
	if err != nil {
		service.logger.Error("failed to generate refresh token", "userId", user.ID.String(), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Store refresh token in database
	_, err = service.refreshTokenService.CreateRefreshToken(user.ID, refreshTokenString)
	if err != nil {
		service.logger.Error("failed to store refresh token", "userId", user.ID.String(), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	response := LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		User:         UserEntityToUserResponse(&user),
	}

	service.logger.Info("user logged in", "userId", user.ID.String(), "email", user.Email)
	c.JSON(http.StatusOK, response)
}

func (service *UsersService) RefreshTokenHandler(c *gin.Context) {
	var refreshInput RefreshTokenDTO
	if err := c.ShouldBindJSON(&refreshInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors := refreshInput.Validate(); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// Get refresh token from database
	refreshToken, err := service.refreshTokenService.GetRefreshToken(refreshInput.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get user information
	var user entities.User
	if err := service.db.Where("id = ?", refreshToken.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			service.logger.Error("failed to find user", "userId", refreshToken.UserID.String(), "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Generate new access token
	accessToken, err := service.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		service.logger.Error("failed to generate access token", "userId", user.ID.String(), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	response := RefreshTokenResponseDTO{
		AccessToken: accessToken,
	}

	service.logger.Info("access token refreshed", "userId", user.ID.String())
	c.JSON(http.StatusOK, response)
}

func (service *UsersService) LogoutHandler(c *gin.Context) {
	// Get refresh token from request body if provided
	var logoutInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&logoutInput); err == nil && logoutInput.RefreshToken != "" {
		// Invalidate the refresh token
		if err := service.refreshTokenService.DeleteRefreshToken(logoutInput.RefreshToken); err != nil {
			service.logger.Error("failed to delete refresh token", "error", err)
		}
	}

	userID, exists := c.Get("user_id")
	if exists {
		// Delete all refresh tokens for this user
		if userUUID, ok := userID.(uuid.UUID); ok {
			if err := service.refreshTokenService.DeleteAllUserRefreshTokens(userUUID); err != nil {
				service.logger.Error("failed to delete user refresh tokens", "userId", userUUID.String(), "error", err)
			}
		}
		service.logger.Info("user logged out", "userId", userID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (service *UsersService) GetProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := service.getUser(c.Request.Context(), userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			service.logger.Error("failed to get user profile", "userId", userUUID.String(), "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, UserEntityToUserResponse(user))
}
