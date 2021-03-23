package control

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index ...
// @Summary	Index
// @Tags	index
// @Accept	json
// @Produce	json
// @Success 200	{string} string "Here is the neuron server"
// @Router / [get]
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Here is the ARK ROBOT server",
	})
}

//NotFinished not implemented
func NotFinished(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "not implemented",
	})
}

//NotSupport ...
func NotSupport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "not support",
	})
}
