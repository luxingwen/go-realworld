package articles

import (
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

// articles
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

// GetAll ...
// @Title 创建文章
// @Description 创建文章
// @Param	body	body 	articles.ArticleModelValidator	true		"body for Culture content"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/ [post]
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

// GetAll ...
// @Title 获取文章列表
// @Description 获取文章列表
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/ [get]
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
		common.HandleErr(c, http.StatusUnauthorized, "Require auth!")
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

// GetAll ...
// @Title 获取文章内容
// @Description 获取文章内容
// @Param   slug  path 	string	true		"slug"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug} [get]
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

// GetAll ...
// @Title 更新文章内容
// @Description 更新文章内容
// @Param   slug  path 	string	true		"slug"
// @Param	body		body 	articles.ArticleModelValidator	true		"body for Culture content"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug} [get]
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

// GetAll ...
// @Title 删除文章
// @Description 删除文章
// @Param   slug  path 	string	true		"slug"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug} [delete]
func ArticleDelete(c *gin.Context) {
	slug := c.Param("slug")
	err := DeleteArticleModel(&ArticleModel{Slug: slug})
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	common.HandleOk(c, gin.H{"article": "Delete success"})
}

// GetAll ...
// @Title 喜欢文章
// @Description 喜欢文章
// @Param   slug  path 	string	true		"slug"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug}/favorite [post]
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

// GetAll ...
// @Title 取消喜欢文章
// @Description 取消喜欢文章
// @Param   slug  path 	string	true		"slug"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug}/favorite [delete]
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

// GetAll ...
// @Title 创建评论
// @Description 创建评论
// @Param   slug  path 	string	true		"slug"
// @Param   body  body  articles.CommentModelValidator  true "body"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug}/comments [post]
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

// GetAll ...
// @Title 删除评论
// @Description 删除评论
// @Param   id    path 	string	true		"id"
// @Param   slug  path 	string	true		"slug"
// @Param	body		body 	articles.CommentModelValidator	true		"body for Culture content"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug}/comments/{id} [delete]
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

// GetAll ...
// @Title 获取文章评论列表
// @Description 获取文章评论列表
// @Param   slug  path 	string	true		"slug"
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/articles/{slug}/comments [get]
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

// GetAll ...
// @Title 获取标签列表
// @Description 获取标签（tags）列表
// @Success 200 {string} json "{"code":0,"data": []string,"msg":"ok"}"
// @router /api/tags [get]
func TagList(c *gin.Context) {
	tagModels, err := getAllTags()
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := TagsSerializer{c, tagModels}
	common.HandleOk(c, gin.H{"tags": serializer.Response()})
}

// GetAll ...
// @Title 获取话题类型列表
// @Description 获取话题类型（types）列表
// @Success 200 {string} json "{"code":0,"data": []*TypeResponse,"msg":"ok"}"
// @router /api/types [get]
func TypeList(c *gin.Context) {
	typeList, err := getAllTypes()
	if err != nil {
		common.HandleErr(c, http.StatusNotFound, err.Error())
		return
	}
	serializer := TypesSerializer{c, typeList}
	common.HandleOk(c, gin.H{"types": serializer.Response()})
}
