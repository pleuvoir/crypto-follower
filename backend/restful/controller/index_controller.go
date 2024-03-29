package controller

import (
	"crypto-follower/restful/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewIndexController() *IndexController {
	return &IndexController{}
}

type IndexController struct {
}

func (t *IndexController) Index(ct *gin.Context) {
	rs := model.Success("success")
	ct.JSON(http.StatusOK, rs)
}

func (t *IndexController) Welcome(c *gin.Context) {
	success := model.Success("welcome")
	success.SetData("u find new place")
	c.JSON(http.StatusOK, success)
}
