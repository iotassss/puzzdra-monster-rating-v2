package handler

import "github.com/gin-gonic/gin"

func (h *Handler) HelloHandler(c *gin.Context) {
	c.HTML(200, "hello.html", gin.H{
		"message": "hello, world",
	})
}
