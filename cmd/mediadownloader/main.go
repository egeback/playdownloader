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

	//Set viper config
	viper.SetDefault("users", map[string]string{"user1": "download"})
	viper.SetDefault("basic_auth", false)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath(".")

	//Load viper configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found")
		} else {
			log.Panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	//Set viper environment parameters
	viper.SetEnvPrefix("DOWNLOADER")
	viper.BindEnv("users", "USERS")
	viper.BindEnv("basicAuth", "BASIC_AUTH")
	viper.AutomaticEnv()

	//Load parameters
	useBasicAuth := viper.GetBool("basic_auth")
	users := viper.GetStringMapString("users")

	svtDlLocation := viper.GetString("SVT_DL_LOCATION")
	if svtDlLocation != "" {
		models.SvtDLLocation = svtDlLocation
	}

	defaultDownloadDir := viper.GetString("DEFAULT_MEDIA_DIRECTORY")
	if defaultDownloadDir != "" {
		models.DefaultDownloadDir = defaultDownloadDir
	}

	//Configure gin routers and middleware
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
			jobs.DELETE("/:uuid", c.StopJob)
		}
		common := v1.Group("/")
		{
			common.GET("ping", c.Ping)
			common.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	//start server and add callbacks for SIGINT
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

	//Wait for SIGINT
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	c.Stop()
	log.Println("Server exiting")
}
