package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"propmanager/api"
	"propmanager/internal/app/middleware"
	"propmanager/internal/app/model"
	"propmanager/internal/app/repository"
	"propmanager/internal/app/service"
	"propmanager/internal/config"
	"propmanager/internal/db"
)

//go:embed docs/swagger.json
var swaggerJson string

func main() {
	cfg := config.LoadConfig()

	db := db.ConnectDB(cfg)

	err := db.AutoMigrate(&model.Property{}, &model.Image{})
	if err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
	}

	propertyRepository := repository.NewPropertyRepository(db)
	propertyService := service.NewPropertyService(propertyRepository)
	s3Service := service.NewS3Service(&cfg)

	propertyHandler := api.NewPropertyHandler(propertyService, s3Service)

	authConfig := config.LoadAuthConfig()
	authService := service.NewAuthService(&authConfig)
	authHandler := api.NewAuthHandler(authService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.GET("/properties", propertyHandler.GetAllProperties)
	r.GET("/properties/:id", propertyHandler.GetProperty)

	authGroup := r.Group("/")
	authGroup.Use(authConfig.AuthMiddleware())
	{
		authGroup.POST("/properties", propertyHandler.CreateProperty)
		authGroup.PUT("/properties/:id", propertyHandler.UpdateProperty)
		authGroup.DELETE("/properties/:id", propertyHandler.DeleteProperty)
		authGroup.POST("/properties/:id/images", propertyHandler.UploadImage)
		authGroup.DELETE("/properties/:id/images/:image_id", propertyHandler.DeleteImage)
	}

	r.POST("/login", authHandler.Login)

	// Serve Swagger UI with custom swagger.json endpoint
	url := ginSwagger.URL("/swagger.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(swaggerJson))
	})

	log.Fatal(r.Run(":" + cfg.Port))
}
