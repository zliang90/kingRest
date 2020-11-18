package router

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/zliang90/kingRest/internal/restful/api/v1"
	apiV2 "github.com/zliang90/kingRest/internal/restful/api/v2"
	"net/http"
)

func InitRouter(env string) http.Handler {
	// decide to disable debug
	if env == "prod" {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	/*------------------------------------ api v1 -------------------------------------*/
	v1 := engine.Group("/api/v1")

	{
		v1.GET("/version", apiV1.Version)

		// configuration
	}

	/*------------------------------------ api v2 -------------------------------------*/
	v2 := engine.Group("/api/v2")
	{
		v2.GET("/version", apiV2.Version)
	}

	return engine
}
