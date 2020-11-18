package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Version(ctx *gin.Context) {
	ctx.String(http.StatusOK, "/api/v1")
}
