package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	// ad
	r.POST("/ads", createAd(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id", updateAd(a))
	r.GET("/ads/:ad_id", getAd(a))
	r.GET("/ads", getAds(a))

	// user
	r.POST("/user", createUser(a))
	r.GET("/user/:user_id", getUser(a))
	r.PUT("/user/:user_id", updateUser(a))
}
