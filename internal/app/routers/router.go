package routers

import (
	"os"
	"strconv"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/app/middlewares"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/services"

	"github.com/gin-gonic/gin"
)

func SetupRoute(
	tagService *services.TagService,
) *gin.Engine {

	// Convert string to bool
	environment, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	if environment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	router := gin.New()
	router.SetTrustedProxies([]string{allowedHosts})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())

	RegisterRoutes(router, tagService) //routes register

	return router
}
