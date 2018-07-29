package users

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luxingwen/go-realworld/common"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UserRegister(router *gin.RouterGroup) {
	router.GET("/", UserRetrieve)
	router.PUT("/", UserUpdate)
}

func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/:username", ProfileRetrieve)
	router.POST("/:username/follow", ProfileFollow)
	router.DELETE("/:username/follow", ProfileUnfollow)
}

func TopUsersAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", GetTopUsers)
}

func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	profileSerializer := ProfileSerializer{c, userModel}
	common.HandleOk(c, gin.H{"profile": profileSerializer.Response()})
}

func ProfileFollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)
	err = myUserModel.following(userModel)
	if err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serializer := ProfileSerializer{c, userModel}
	common.HandleOk(c, gin.H{"profile": serializer.Response()})
}

func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)

	err = myUserModel.unFollowing(userModel)
	if err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serializer := ProfileSerializer{c, userModel}
	common.HandleOk(c, gin.H{"profile": serializer.Response()})
}

func UsersRegistration(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := SaveOne(&userModelValidator.userModel); err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("Duplicate entry '%s'", userModelValidator.userModel.Username)) {
			err = errors.New("该用户名已经被占用")
		}
		if strings.Contains(err.Error(), fmt.Sprintf("Duplicate entry '%s'", userModelValidator.userModel.Email)) {
			err = errors.New("该邮箱已经被注册")
		}
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)
	serializer := UserSerializer{c}
	common.HandleOk(c, gin.H{"user": serializer.Response()})
}

func UsersLogin(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	userModel, err := FindOneUser(&UserModel{Email: loginValidator.userModel.Email})

	if err != nil {
		common.HandleErr(c, http.StatusForbidden, err.Error())
		return
	}

	if userModel.checkPassword(loginValidator.User.Password) != nil {
		common.HandleErr(c, http.StatusForbidden, err.Error())
		return
	}
	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{c}
	common.HandleOk(c, gin.H{"user": serializer.Response()})
}

func UserRetrieve(c *gin.Context) {
	serializer := UserSerializer{c}
	common.HandleOk(c, gin.H{"user": serializer.Response()})
}

func UserUpdate(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(UserModel)
	userModelValidator := NewUserModelValidatorFillWith(myUserModel)
	if err := userModelValidator.Bind(c); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	userModelValidator.userModel.ID = myUserModel.ID
	if err := myUserModel.Update(userModelValidator.userModel); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	UpdateContextUserModel(c, myUserModel.ID)
	serializer := UserSerializer{c}
	common.HandleOk(c, gin.H{"user": serializer.Response()})
}

func GetTopUsers(c *gin.Context) {
	users, err := TopUsers()
	if err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serializer := TopUsersSerializer{c, users}
	common.HandleOk(c, gin.H{"users": serializer.Response()})
}
