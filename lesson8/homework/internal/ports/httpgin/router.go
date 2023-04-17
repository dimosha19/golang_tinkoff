package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	// ad
	r.POST("/ads", myMV, createAd(a))
	r.PUT("/ads/:ad_id/status", myMV, changeAdStatus(a))
	r.PUT("/ads/:ad_id", myMV, updateAd(a))
	r.GET("/ads/:ad_id", myMV, getAd(a))
	r.GET("/ads", myMV, getAds(a))

	// user
	r.POST("/user", myMV, createUser(a))
	r.GET("/user/:user_id", myMV, getUser(a))
	r.PUT("/user/:user_id", myMV, updateUser(a))
}
