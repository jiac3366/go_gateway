package router

import (
	"github.com/e421083458/go_gateway/controller"
	"github.com/e421083458/go_gateway/docs"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//设置cookie存储
	//https://github.com/gin-contrib/sessions
	store, _ := redis.NewStore(10, "tcp", lib.GetStringConf("redis_map.session.server"), "", []byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	})

	//admin_login
	adminLogin := router.Group("/admin_login")
	adminLogin.Use(
		sessions.Sessions("mysession", store),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AdminLoginRegister(adminLogin)
	}

	//admin
	admin := router.Group("/admin")
	admin.Use(
		middleware.TranslationMiddleware(),
		sessions.Sessions("mysession", store),
		middleware.AdminSessionAuthMiddleware(),
	)
	{
		controller.AdminRegister(admin)
	}

	//service
	service := router.Group("/service")
	service.Use(
		middleware.TranslationMiddleware(),
		sessions.Sessions("mysession", store),
		middleware.AdminSessionAuthMiddleware(),
	)
	{
		controller.ServiceRegister(service)
	}

	//app
	app := router.Group("/app")
	app.Use(
		middleware.TranslationMiddleware(),
		sessions.Sessions("mysession", store),
		middleware.AdminSessionAuthMiddleware(),
	)
	{
		controller.APPRegister(app)
	}

	//app
	dash := router.Group("/dashboard")
	dash.Use(
		middleware.TranslationMiddleware(),
		sessions.Sessions("mysession", store),
		middleware.AdminSessionAuthMiddleware(),
	)
	{
		controller.DashBoardRegister(dash)
	}
	return router
}