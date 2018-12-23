package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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
	quit := make(chan os.Signal, 1)
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
	stats, err := GetStats()
	if err != nil {
		log.Debug(fmt.Errorf("something goes wrong: %s", err))
	}
	c.String(http.StatusOK, stats)
}

func GetStats() (string, error) {

	statsPath := config.GetString("stats_path")

	log.Debug("Stats directory: ", statsPath)

	var returnErr error

	var sb strings.Builder
	walkErr := filepath.Walk(statsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			returnErr = err
		}
		// get relative path to stats dir
		relPath, relErr := filepath.Rel(statsPath, path)

		if relErr != nil {
			returnErr = relErr
		}

		// replace / with _
		var underscorePath = strings.Replace(relPath, "/", "_", -1)
		if !info.IsDir() {
			dat, err := ioutil.ReadFile(path)
			check(err)

			log.Debug("path:", path, "FileInfo.Name:", info.Name(), "Dir:", underscorePath)
			sb.WriteString("gateway_")
			sb.WriteString(underscorePath)
			sb.WriteString(" ")
			sb.WriteString(string(dat))
			sb.WriteString("\n")
		}

		return nil
	})

	if walkErr != nil {
		returnErr = walkErr
		result, err := fmt.Fprintf(os.Stderr, "walk failed with error: %v\n", walkErr)
		log.Error(result, err)
	}

	return sb.String(), returnErr
}
