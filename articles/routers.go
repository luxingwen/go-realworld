package articles

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luxingwen/go-realworld/common"
	"github.com/luxingwen/go-realworld/users"
)

func ArticlesRegister(router *gin.RouterGroup) {
	router.POST("/", ArticleCreate)
	router.PUT("/:slug", ArticleUpdate)
	router.DELETE("/:slug", ArticleDelete)
	router.POST("/:slug/favorite", ArticleFavorite)
	router.DELETE("/:slug/favorite", ArticleUnfavorite)
	router.POST("/:slug/comments", ArticleCommentCreate)
	router.DELETE("/:slug/comments/:id", ArticleCommentDelete)
}

func ArticlesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", ArticleList)
	router.GET("/:slug", ArticleRetrieve)
	router.GET("/:slug/comments", ArticleCommentList)
}

func TagsAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", TagList)
}

func TypesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", TypeList)
}

func ArticleCreate(c *gin.Context) {
	articleModelValidator := NewArticleModelValidator()
	if err := articleModelValidator.Bind(c); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	//fmt.Println(articleModelValidator.articleModel.Author.UserModel)

	if err := SaveOne(&articleModelValidator.articleModel); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serializer := ArticleSerializer{c, articleModelValidator.articleModel}
	common.HandleOk(c, gin.H{"article": serializer.Response()})
}

func ArticleList(c *gin.Context) {
	//condition := ArticleModel{}
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")
	limit := c.Query("limit")
	offset := c.Query("offset")
	typ := c.Query("type")
	articleModels, modelCount, err := FindManyArticle(tag, author, limit, offset, favorited, typ)
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := ArticlesSerializer{c, articleModels}
	common.HandleOk(c, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
}

func ArticleFeed(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	if myUserModel.ID == 0 {
		c.AbortWithError(http.StatusUnauthorized, errors.New("{error : \"Require auth!\"}"))
		return
	}
	articleUserModel := GetArticleUserModel(myUserModel)
	articleModels, modelCount, err := articleUserModel.GetArticleFeed(limit, offset)
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := ArticlesSerializer{c, articleModels}
	common.HandleOk(c, gin.H{"articles": serializer.Response(), "articlesCount": modelCount})
}

func ArticleRetrieve(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "feed" {
		ArticleFeed(c)
		return
	}
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := ArticleSerializer{c, articleModel}
	common.HandleOk(c, gin.H{"article": serializer.Response()})
}

func ArticleUpdate(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	articleModelValidator := NewArticleModelValidatorFillWith(articleModel)
	if err := articleModelValidator.Bind(c); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	articleModelValidator.articleModel.ID = articleModel.ID
	if err := articleModel.Update(articleModelValidator.articleModel); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serializer := ArticleSerializer{c, articleModel}
	common.HandleOk(c, gin.H{"article": serializer.Response()})
}

func ArticleDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := DeleteArticleModel(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	common.HandleOk(c, gin.H{"article": "Delete success"})
}

func ArticleFavorite(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	err = articleModel.favoriteBy(GetArticleUserModel(myUserModel))
	serializer := ArticleSerializer{c, articleModel}
	common.HandleOk(c, gin.H{"article": serializer.Response()})
}

func ArticleUnfavorite(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	myUserModel := c.MustGet("my_user_model").(users.UserModel)
	err = articleModel.unFavoriteBy(GetArticleUserModel(myUserModel))
	serializer := ArticleSerializer{c, articleModel}
	common.HandleOk(c, gin.H{"article": serializer.Response()})
}

func ArticleCommentCreate(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	commentModelValidator := NewCommentModelValidator()
	if err := commentModelValidator.Bind(c); err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	commentModelValidator.commentModel.Article = articleModel

	if err := SaveOne(&commentModelValidator.commentModel); err != nil {
		common.HandleErr(c, http.StatusUnprocessableEntity, err.Error())
		return
	}
	serializer := CommentSerializer{c, commentModelValidator.commentModel}
	common.HandleOk(c, gin.H{"comment": serializer.Response()})
}

func ArticleCommentDelete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	err = DeleteCommentModel([]uint{id})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	common.HandleOk(c, gin.H{"comment": "Delete success"})
}

func ArticleCommentList(c *gin.Context) {
	slug := c.Param("slug")
	articleModel, err := FindOneArticle(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	err = articleModel.getComments()
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := CommentsSerializer{c, articleModel.Comments}
	common.HandleOk(c, gin.H{"comments": serializer.Response()})
}
func TagList(c *gin.Context) {
	tagModels, err := getAllTags()
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := TagsSerializer{c, tagModels}
	common.HandleOk(c, gin.H{"tags": serializer.Response()})
}

func TypeList(c *gin.Context) {
	typeList, err := getAllTypes()
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := TypesSerializer{c, typeList}
	common.HandleOk(c, gin.H{"types": serializer.Response()})
}
