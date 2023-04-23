package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"homework9/internal/app"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(myMV)
	api := handler.Group("/api/v1")
	s := &http.Server{Addr: port, Handler: handler}
	AppRouter(api, a)

	return s
}
