package api

import (
	"strings"

	"github.com/chenchenyu/gomysqldemo/dao"
	"github.com/chenchenyu/gomysqldemo/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func InitRouter(r *gin.Engine) {
	r.GET("/ping", ping)

	r.POST("/user", UserCreate)
	r.DELETE("/user", UserDeleteByID)
	r.PUT("/user", UserUpdateByID)
	r.GET("/user", UserSelectByName)

}

func UserCreate(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(401, gin.H{"code": -1, "message": err.Error()})
		return
	}
	err = dao.NewMysqlDao(c).UserInsert(c, &user)
	if err != nil {
		c.JSON(500, gin.H{"code": -2, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": user})
}

func UserUpdateByID(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(401, gin.H{"code": -1, "message": err.Error()})
		return
	}
	err = dao.NewMysqlDao(c).UserUpdateByID(c, &user)
	if err != nil {
		c.JSON(401, gin.H{"code": -2, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": user})
}

func UserDeleteByID(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(401, gin.H{"code": -1, "message": err.Error()})
		return
	}
	err = dao.NewMysqlDao(c).UserDeleteByID(c, &user)
	if err != nil {
		c.JSON(401, gin.H{"code": -2, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": user})
}

func UserSelectByName(c *gin.Context) {
	nameQuery := c.Query("names")
	names := strings.Split(nameQuery, ",")
	spew.Dump(names)
	if len(names) == 0 {
		c.JSON(401, gin.H{"code": -1, "message": "query is empty"})
		return
	}
	users, err := dao.NewMysqlDao(c).UsersSelectByName(c, names)
	if err != nil {
		c.JSON(500, gin.H{"code": -2, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": users})
}
