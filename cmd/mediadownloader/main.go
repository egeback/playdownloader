package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/egeback/playdownloader/internal/controllers"
	_ "github.com/egeback/playdownloader/internal/docs"
	"github.com/egeback/playdownloader/internal/models"
	"github.com/egeback/playdownloader/internal/utils"
	"github.com/egeback/playdownloader/internal/version"
	"github.com/gin-gonic/gin"
	jsonp "github.com/tomwei7/gin-jsonp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/spf13/viper"
)

var address = ":8081"

// @title Play Media API - Downloader
// @version 1.0
// @description API to download with svt-download

// @contact.name API Support
// @contact.url http://xxxx.xxx.xx
// @contact.email support@egeback.se

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1/
func main() {
	fmt.Printf("%s Running Play Media API - Downloader version: %s (%s)\n", time.Now().Format("2006-01-02 15:04:05"), version.BuildVersion, version.BuildTime)

	svtDlLocation := os.Getenv("SVT_DL_LOCATION")
	if svtDlLocation != "" {
		models.SvtDLLocation = svtDlLocation
	}
	defaultDownloadDir := os.Getenv("DEFAULT_MEDIA_DIRECTORY")
	if defaultDownloadDir != "" {
		models.DefaultDownloadDir = defaultDownloadDir
	}

	viper.SetDefault("users", map[string]string{"user1": "download"})
	viper.SetDefault("basic_auth", false)

	viper.SetEnvPrefix("DOWNLOADER")
	viper.BindEnv("users", "USERS")
	viper.BindEnv("basic_auth", "BASIC_AUTH")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found")
		} else {
			log.Panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	useBasicAuth := false
	if os.Getenv("DOWNLOADER_BASIC_AUTH") != "" {
		useBasicAuth = *utils.GetBoolValueFromString(os.Getenv("DOWNLOADER_BASIC_AUTH"), false)
	} else {
		useBasicAuth = viper.GetBool("basic_auth")
	}
	users := viper.GetStringMapString("users")

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(jsonp.JsonP())
	c := controllers.NewController()
	v1 := r.Group("/api/v1")
	{
		var jobs *gin.RouterGroup
		if useBasicAuth {
			jobs = v1.Group("/jobs", gin.BasicAuth(gin.Accounts(users)))
		} else {
			jobs = v1.Group("/jobs")
		}
		{
			jobs.POST("", c.AddJob)
			jobs.GET("", c.Jobs)
			jobs.POST("/", c.AddJob)
			jobs.GET("/", c.Jobs)

			jobs.GET("/:uuid", c.GetJob)
		}
		common := v1.Group("/")
		{
			common.GET("ping", c.Ping)
			common.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}
	quit := make(chan os.Signal)

	go func() {
		// service connections
		fmt.Printf("Listening and serving HTTP on %s\n", address)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
			if strings.Index(err.Error(), "address already in use") >= 0 {
				quit <- syscall.SIGINT
			}
		}
	}()

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
