package app

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"

	"github.com/matiasronny13/go-note/config"
	"github.com/matiasronny13/go-note/internal/pkg/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type webServer struct {
	appConfig *config.AppConfig
}

func NewWebServer(config *config.AppConfig) *webServer {
	return &webServer{config}
}

func createMyRender(appConfig *config.AppConfig) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index",
		path.Join(appConfig.Web.BasePath, "web/template/shared/base.html"),
		path.Join(appConfig.Web.BasePath, "web/template/home/index.html"),
		path.Join(appConfig.Web.BasePath, "web/template/shared/error_list.html"))
	r.AddFromFiles("index_partial", path.Join(appConfig.Web.BasePath, "web/template/home/index.html"), path.Join(appConfig.Web.BasePath, "web/template/shared/error_list.html"))
	r.AddFromFiles("error_list", path.Join(appConfig.Web.BasePath, "web/template/shared/error_list.html"))
	r.AddFromFiles("error", path.Join(appConfig.Web.BasePath, "web/template/shared/base.html"), path.Join(appConfig.Web.BasePath, "web/template/shared/error.html"))
	return r
}

func ConfigureWebRouter(appConfig *config.AppConfig, dbClient *supabase.Client) *gin.Engine {
	var crypto CryptoService = NewCryptoService(os.Getenv("CRYPTO_KEY"), os.Getenv("CRYPTO_IV_PAD"))
	var service AppService = NewAppService(appConfig, dbClient, crypto)
	var handler WebHandler = NewWebHandler(appConfig, service)

	router := gin.Default()
	router.HTMLRender = createMyRender(appConfig)

	// Middlewares
	router.Use(cors.Default())
	router.Use(gin.CustomRecovery(handler.UnexpectedError))
	router.Use(handler.AuthenticateUser())

	// Static Assets
	router.Static("/assets", path.Join(appConfig.Web.BasePath, "web/assets"))
	router.StaticFile("/favicon.ico", path.Join(appConfig.Web.BasePath, "web/favicon.ico"))

	// Routers
	router.GET("", handler.Default)
	router.GET("/:id", handler.GetById)
	router.POST("/:id/delete", handler.DeleteById)
	router.POST("/encrypt", handler.Encrypt)
	router.POST("/:id/encrypt", handler.Encrypt)
	router.POST("/decrypt", handler.Decrypt)
	router.POST("/:id/decrypt", handler.Decrypt)
	router.POST("", handler.Create)
	router.POST("/:id", handler.Save)

	return router
}

func (r *webServer) Run() {
	if dbClient, err := supabase.NewClient(r.appConfig.SupabaseUrl, r.appConfig.SupabaseKey, nil); err != nil {
		panic(err)
	} else {
		slog.Info("Database connected")
		if data, _, err := dbClient.From("test").Select("id, value", "exact", false).Execute(); err == nil {
			var result []*model.Test
			err := json.Unmarshal(data, &result)
			if err == nil {
				slog.Info(result[0].Value)
			}
		}
		ConfigureWebRouter(r.appConfig, dbClient).Run(r.appConfig.Web.Host)
	}
}
