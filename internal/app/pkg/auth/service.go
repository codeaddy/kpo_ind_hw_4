package auth

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"hw4/internal/app/pkg/core"
	"hw4/internal/app/pkg/session"
	"hw4/internal/app/pkg/user"
	"net/http"
	"strings"
	"time"
)

type AuthService struct {
	core *core.CoreService
}

func NewService(core *core.CoreService) *AuthService {
	return &AuthService{core: core}
}

type registerNewUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type registerNewUserResponse struct {
	Id int `json:"id"`
}

// RegisterNewUser godoc
// @Summary		RegisterNewUser
// @Description	Doing new user registration
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		registerNewUserInput	true	"New User"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/auth/register-new-user [post]
func (s *AuthService) RegisterNewUser(c *gin.Context) {
	var input registerNewUserInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error when unmarshalling input"})
		return
	}

	if !strings.Contains(input.Email, "@") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email doesn't contain @"})
		return
	}

	if input.Role != "chef" && input.Role != "manager" {
		input.Role = "customer"
	}

	id, err := s.core.UserService.Create(c, user.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hex.EncodeToString([]byte(input.Password)),
		Role:         input.Role,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(id)
	c.IndentedJSON(http.StatusOK, registerNewUserResponse{Id: id})
}

type authorizationInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authorizationResponse struct {
	Jwt string `json:"jwt"`
}

// Authorization godoc
// @Summary		Authorization
// @Description	Provides authorization
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		authorizationInput	true	"User's email and password"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/auth/authorization [post]
func (s *AuthService) Authorization(c *gin.Context) {
	fmt.Println("asd")
	var input authorizationInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error when unmarshalling input"})
		return
	}

	passHash := hex.EncodeToString([]byte(input.Password))

	user, err := s.core.UserService.GetByEmailPassword(c, input.Email, passHash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := s.core.SessionService.Create(c, session.Session{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session, err := s.core.SessionService.GetById(c, id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, authorizationResponse{Jwt: session.SessionToken})
}

type getUserInfoInput struct {
	Jwt string `json:"jwt"`
}

type getUserInfoResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// GetUserInfo godoc
// @Summary		GetUserInfo
// @Description	Getting user information
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		getUserInfoResponse	true	"User's username, email and role"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/auth/get-user-info [post]
func (s *AuthService) GetUserInfo(c *gin.Context) {
	fmt.Println("GetUserInfo")
	var input getUserInfoInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error when unmarshalling input"})
		return
	}

	session, err := s.core.SessionService.GetByToken(c, input.Jwt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := s.core.UserService.GetById(c, session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, getUserInfoResponse{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})
}
