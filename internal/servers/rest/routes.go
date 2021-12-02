package rest

import (
	"github.com/Format-C-eft/middleware/internal/api"
	"github.com/Format-C-eft/middleware/internal/config"
	"github.com/gin-gonic/gin"
)

func initRoutes(router *gin.Engine, cfg *config.Config) {

	router.GET(cfg.Services.Rest.Path+"/auth", api.Login)

	router.GET(cfg.Services.Rest.Path+"/check-login", api.CheckLogin)

	groupJWT := router.Group(cfg.Services.Rest.Path)
	groupJWT.Use(CheckBearerAuth())
	groupJWT.DELETE("auth", api.LogOut)

	groupSession := groupJWT.Group("sessions")
	groupSession.GET("", api.GetSession)
	groupSession.GET(":UUID", api.GetSession)
	groupSession.DELETE(":UUID", api.DropSession)

	for _, v := range cfg.Servers.OneC.Routes {
		groupV := groupJWT.Group(v)
		groupV.GET("", api.OtherMetods)
		groupV.GET("/*actions", api.OtherMetods)
		groupV.PUT("", api.OtherMetods)
		groupV.PUT("/*action", api.OtherMetods)
		groupV.POST("", api.OtherMetods)
		groupV.POST("/*action", api.OtherMetods)
		groupV.DELETE("", api.OtherMetods)
		groupV.DELETE("*action", api.OtherMetods)
	}

	router.NoRoute(api.MetodOrPatchNotFound)
}
