package main

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"time"
)

// publish
type PublishBody struct {
	Title   string `json: "title"`
	Cover   string `json: "cover"`
	Content string `json: "content"`
}

func (u *UserAPI) publish(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body PublishBody
	ctx.ReadJSON(&body)
	_uuid := uuid.NewV4().String()

	article := Article{
		Id:        _uuid,
		Title:     body.Title,
		Cover:     body.Cover,
		Content:   body.Content,
		Timestamp: time.Now().Unix(),
	}

	var user User
	u.db.Where("mail = ?", mail).Find(&user)
	u.db.Create(&article)

	_articles := append(user.Articles, _uuid)
	u.db.Model(&user).Where("mail = ?", mail).Update("articles", _articles)

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

// updateUser
type UpdateUserNameBody struct {
	Name string `json:name`
}

func (u *UserAPI) updateUserName(ctx iris.Context) {
	mail := ctx.Params().Get("mail")
	var body UpdateUserNameBody
	ctx.ReadJSON(&body)

	var user User
	if err := u.db.Where("name = ?", body.Name).Find(&user).Error; err != nil {
		u.db.Model(&user).Where("mail = ?", mail).Update("name", body.Name)
		u.db.Where("name = ?", body.Name).Find(&user)

		ctx.JSON(iris.Map{
			"msg": "ok",
			"data": iris.Map{
				"mail": user.Mail,
				"name": user.Name,
			},
		})

		return
	}

	ctx.StatusCode(iris.StatusBadRequest)
}

/// updateArticle
type UpdatePostBody struct {
	Title   string `json:title`
	Cover   string `json:cover`
	Content string `json:content`
}

func (u *UserAPI) updatePost(ctx iris.Context) {
	id := ctx.Params().Get("id")
	mail := ctx.Params().Get("mail")

	var body UpdatePostBody
	ctx.ReadJSON(&body)

	var user User
	article := Article{
		Id: id,
	}

	u.db.Where("mail = ?", mail).Find(&user)
	var _arr []string = user.Articles
	if err := u.db.Where("id IN (?)", _arr).Find(&article).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := u.db.Model(&article).Where("id = ?", id).Updates(map[string]interface{}{
		"title":   body.Title,
		"cover":   body.Cover,
		"content": body.Content,
	}).Error; err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}

// delete article
func (u *UserAPI) deletePost(ctx iris.Context) {
	id := ctx.Params().Get("id")
	mail := ctx.Params().Get("mail")

	// delete post in user
	var user User
	u.db.Where("mail = ?", mail).Find(&user)

	index := 0
	for i, b := range user.Articles {
		if b == id {
			index = i
		}
	}

	if index == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	_arr := user.Articles
	_arr[index] = _arr[len(_arr)-1]
	_arr = _arr[:len(_arr)-1]

	u.db.Model(&user).Where("mail = ?", mail).Update("articles", _arr)

	// delete post
	article := Article{
		Id: id,
	}

	u.db.Delete(&article)

	ctx.JSON(iris.Map{
		"msg": "ok",
	})
}
