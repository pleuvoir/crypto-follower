package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pleuvoir/gamine/component/restful"
	"net/http"
)

func NewIndexController() *IndexController {
	return &IndexController{}
}

type IndexController struct {
}

func (t *IndexController) Index(ct *gin.Context) {
	rs := restful.Success("success")
	ct.JSON(http.StatusOK, rs)
}

func (t *IndexController) Welcome(c *gin.Context) {
	success := restful.Success("welcome")
	success.SetData("u find new place")
	c.JSON(http.StatusOK, success)
}
