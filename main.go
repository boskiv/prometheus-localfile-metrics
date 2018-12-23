package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var config *viper.Viper

func main() {
	log.SetLevel(log.InfoLevel)
	router := gin.Default()
	config = viper.New()
	if os.Getenv("GPE_STATS_PATH") == "" {
		log.Info("env GPE_STATS_PATH is not set, using default")
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		err = os.Setenv("GPE_STATS_PATH", path.Join(dir, "stats"))
		if err != nil { // Handle errors reading the config file
			log.Error(fmt.Errorf("can not set env var gpe_stats_path: %s \n", err))
		}
	}

	log.Info("env GPE_STATS_PATH: ", os.Getenv("GPE_STATS_PATH"))

	config.SetEnvPrefix("GPE")
	err := config.BindEnv("stats_path")
	if err != nil { // Handle errors reading the config file
		log.Error(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.GET("/metrics", metricsHandler)

	srv := &http.Server{
		Addr:    ":9102",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
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

func metricsHandler(c *gin.Context) {
	c.String(http.StatusOK, GetStats())
}

func GetStats() string {

	statsPath := config.GetString("stats_path")


	log.Debug("Stats directory: ", statsPath)

	var sb strings.Builder
	err := filepath.Walk(statsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// get relative path to stats dir
		relPath, err := filepath.Rel(statsPath,path)
		// replace / with _
		var underscorePath = strings.Replace(relPath, "/", "_", -1)
		if ! info.IsDir() {
			log.Debug("path:", path, "FileInfo.Name:", info.Name(), "Dir:", underscorePath)
			sb.WriteString("gateway_")
			sb.WriteString(underscorePath)
			sb.WriteString("\n")
		}

		return nil
	})

	if err != nil {
		result, err := fmt.Fprintf(os.Stderr, "walk failed with error: %v\n", err)
		log.Fatal(result, err)
	}

	return sb.String()
}