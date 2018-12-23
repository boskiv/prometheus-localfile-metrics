package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/toorop/gin-logrus"
)

var config *viper.Viper
var log = logrus.New()

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	log.SetLevel(logrus.DebugLevel)
	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())

	config = viper.New()
	if os.Getenv("PLM_STATS_PATH") == "" {
		log.Info("env PLM_STATS_PATH is not set, using default")
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		err = os.Setenv("PLM_STATS_PATH", path.Join(dir, "stats"))
		if err != nil { // Handle errors reading the config file
			log.Error(fmt.Errorf("can not set env var gpe_stats_path: %s", err))
		}
	}

	log.Info("env PLM_STATS_PATH: ", os.Getenv("PLM_STATS_PATH"))

	if os.Getenv("PLM_STATS_PREFIX") == "" {
		log.Info("env PLM_STATS_PREFIX is not set, using default")

		err := os.Setenv("PLM_STATS_PREFIX", "myapp")
		if err != nil { // Handle errors reading the config file
			log.Error(fmt.Errorf("can not set env var gpe_stats_prefix: %s", err))
		}
	}

	log.Info("env PLM_STATS_PREFIX: ", os.Getenv("PLM_STATS_PREFIX"))

	config.SetEnvPrefix("PLM")
	err := config.BindEnv("stats_path")
	if err != nil { // Handle errors reading the config file
		log.Error(fmt.Errorf("Fatal error config file: %s", err))
	}

	err = config.BindEnv("stats_prefix")
	if err != nil { // Handle errors reading the config file
		log.Error(fmt.Errorf("Fatal error config file: %s", err))
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

// GetStats read a directory from config
// find all files and folders
// make stat name from PLM_PREFIX variable and relative dir path
// to PLM_STATS_PATH and filename
// get metric from file content
func GetStats() (string, error) {

	statsPath := config.GetString("stats_path")
	statsPrefix := config.GetString("stats_prefix")

	log.Debug("Stats directory: ", statsPath)

	var sb strings.Builder
	walkErr := filepath.Walk(statsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// get relative path to stats dir
		relPath, err := filepath.Rel(statsPath, path)

		if err != nil {
			return err
		}

		// replace / with _
		var underscorePath = strings.Replace(relPath, "/", "_", -1)
		if !info.IsDir() {
			dat, err := ioutil.ReadFile(path)
			check(err)

			log.Debug("path:", path, "FileInfo.Name:", info.Name(), "Dir:", underscorePath)
			sb.WriteString(statsPrefix)
			sb.WriteString("_")
			sb.WriteString(underscorePath)
			sb.WriteString(" ")
			sb.WriteString(string(dat))
			sb.WriteString("\n")
		}

		return nil
	})

	if walkErr != nil {
		log.Error("walk failed with error:", walkErr)
	}

	return sb.String(), walkErr
}
