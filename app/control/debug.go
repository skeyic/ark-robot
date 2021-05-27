package control

import (
	"github.com/gin-gonic/gin"
	"github.com/skeyic/ark-robot/config"
	"net/http"
)

func Debug(c *gin.Context) {
	c.HTML(http.StatusOK, config.Config.ResourceFolder+"abc.html", gin.H{})
}
