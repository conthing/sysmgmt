package handlers

import (
	"math/rand"
	"net/http"
	"strings"

	"github.com/conthing/sysmgmt/auth"
	"github.com/conthing/sysmgmt/db"
	"github.com/conthing/sysmgmt/models"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// todo 不同主机token可以通用bug，token时间放长一点

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PasswdRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newpassword"`
}

func refreshTokenRandom() int {
	return 1 + rand.Intn(0x40000000)
}

func Ready(c *gin.Context) {
	users, err := db.GetUserList()
	if err != nil || len(users) == 0 {
		common.Log.Infof("No user in db, need to create")
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusPreconditionFailed,
			Message: "No user in db, need to create",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
	})
}

// SignUp 注册
// 用户角色的变更需要 SQL
func SignUp(c *gin.Context) {
	var info LoginRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		})
		return
	}

	err = db.GetUser(&models.User{Username: info.Username})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "username already in use",
		})
		return
	}

	// Hash Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: "GenerateFromPassword:" + err.Error(),
		})
		return
	}
	hash := string(bytes)

	user := models.User{Username: info.Username, Hash: hash, TokenRandom: refreshTokenRandom()}
	err = db.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: "Create models.User failed: " + err.Error(),
		})
		return
	}

	// 设置参数到middleware中，校验更快
	//auth.SetTokenRandom(user.Username, user.TokenRandom)

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Sign up success",
	})

}

// Login 根据用户名和密码登陆
func Login(c *gin.Context) {
	var info LoginRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		})
		return
	}

	user := models.User{Username: info.Username}
	err = db.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid username or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(info.Password))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid username or password",
		})
		return
	}

	// 签发token
	token, err := auth.GenerateToken(user.Username, user.TokenRandom)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: LoginResponse{
			Token: token,
		},
	})

}

// Passwd 修改密码
func Passwd(c *gin.Context) {
	var info PasswdRequest
	err := c.ShouldBindJSON(&info)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "invalid json body",
		})
		return
	}

	user := models.User{Username: info.Username}
	err = db.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid username or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(info.Password))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid username or password",
		})
		return
	}

	// Hash Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(info.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: "GenerateFromPassword:" + err.Error(),
		})
		return
	}
	hash := string(bytes)

	user.Hash = hash
	user.TokenRandom = refreshTokenRandom()
	err = db.ModifyUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: "ModifyUser failed: " + err.Error(),
		})
		return
	}

	// 设置参数到middleware中，校验更快
	//auth.SetTokenRandom(user.Username, user.TokenRandom)

	// 签发token
	token, err := auth.GenerateToken(user.Username, user.TokenRandom)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: LoginResponse{
			Token: token,
		},
	})

}

// Logout 根据用户名和密码登陆
func Logout(c *gin.Context) {
	var username string
	value, ok := c.Get("username")
	if ok {
		username, ok = value.(string)
	}
	if !ok || username == "" {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	user := models.User{Username: username, TokenRandom: refreshTokenRandom()}
	err := db.ModifyUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: "ModifyUser failed: " + err.Error(),
		})
		return
	}

	// 设置参数到middleware中，校验更快
	//auth.SetTokenRandom(user.Username, user.TokenRandom)

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Logout success",
	})
}
