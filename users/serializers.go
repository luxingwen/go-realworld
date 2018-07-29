package users

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/luxingwen/go-realworld/common"
	"github.com/luxingwen/go-realworld/config"
)

type ProfileSerializer struct {
	C *gin.Context
	UserModel
}

// Declare your response schema here
type ProfileResponse struct {
	ID        uint   `json:"-"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

// Put your response logic including wrap the userModel here.
func (self *ProfileSerializer) Response() ProfileResponse {
	myUserModel := self.C.MustGet("my_user_model").(UserModel)
	profile := ProfileResponse{
		ID:        self.ID,
		Username:  self.Username,
		Bio:       self.Bio,
		Image:     self.Image,
		Following: myUserModel.isFollowing(self.UserModel),
	}
	if profile.Image == "" {
		profile.Image = "http://luxingwen.github.io/images/git01.jpg"
	} else {
		profile.Image = getImgUrl(profile.Image)
	}
	return profile
}

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.c.MustGet("my_user_model").(UserModel)
	user := UserResponse{
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
		Bio:      myUserModel.Bio,
		Image:    myUserModel.Image,
		Token:    common.GenToken(myUserModel.ID),
	}
	if user.Image == "" {
		user.Image = "http://luxingwen.github.io/images/git01.jpg"
	} else {
		user.Image = getImgUrl(user.Image)
	}
	return user
}

type TopUserSerializer struct {
	C    *gin.Context
	User UserModel
}

type TopUserResponse struct {
	Username string `json:"username"`
	Image    string `json:"image"`
}

type TopUsersSerializer struct {
	C     *gin.Context
	Users []*UserModel
}

func (this *TopUserSerializer) Response() *TopUserResponse {
	r := &TopUserResponse{Username: this.User.Username, Image: this.User.Image}
	if r.Image == "" {
		r.Image = "http://luxingwen.github.io/images/git01.jpg"
	} else {
		r.Image = getImgUrl(r.Image)
	}
	return r
}

func (this *TopUsersSerializer) Response() (r []*TopUserResponse) {
	for _, item := range this.Users {
		u := TopUserSerializer{this.C, *item}
		r = append(r, u.Response())
	}
	return
}

func getImgUrl(urlstr string) string {
	if strings.Contains(urlstr, "http://") || strings.Contains(urlstr, "https://") {
		return urlstr
	}
	return config.ServerConf.FilePath + urlstr
}
