package routers

import (
	"net/http"

	"product-service/internal/app/controllers"
	"product-service/internal/domain/services"

	_ "product-service/docs"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	route *gin.Engine,
	tagService *services.TagService,
) {
	/* Controllers */
	tagController := controllers.NewTagController(tagService)

	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	v1 := route.Group("api/v1")
	{
		// health check
		v1.GET("health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "good"}) })
		// tags
		tags := v1.Group("tags")
		{
			tags.GET("", tagController.GetTags)
			tags.GET(":id", tagController.GetTagById)
			tags.POST("", tagController.SaveTag)
			tags.PUT(":id", tagController.UpdateTag)
			tags.DELETE(":id", tagController.DeleteTag)
		}
	}
}
