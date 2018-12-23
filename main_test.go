package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func PrepareStats(tmpDir string) {
	log.Debug("Creating stats dir")
	var err error
	err = os.MkdirAll(path.Join(tmpDir, "stats", "some"), os.ModePerm)
	if err != nil {
		log.Error(fmt.Errorf("can not create test path: %s \n", err))
	}

	err = os.MkdirAll(path.Join(tmpDir, "stats", "timers"), os.ModePerm)
	if err != nil {
		log.Error(fmt.Errorf("can not create test path: %s \n", err))
	}

	stat1Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "stat_1"), stat1Data, 0644)
	check(err)

	stat2Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "stat_2"), stat2Data, 0644)
	check(err)

	stat3Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "stat_3"), stat3Data, 0644)
	check(err)

	some1Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "some", "some_1"), some1Data, 0644)
	check(err)

	some2Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "some", "some_2"), some2Data, 0644)
	check(err)

	some3Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "some", "some_3"), some3Data, 0644)
	check(err)

	timer1Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "timers", "timer_1"), timer1Data, 0644)
	check(err)

	timer2Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "timers", "timer_2"), timer2Data, 0644)
	check(err)

	timer3Data := []byte("100")
	err = ioutil.WriteFile(path.Join(tmpDir, "stats", "timers", "timer_3"), timer3Data, 0644)
	check(err)

	log.Debug("Created stats dir:", filepath.Dir(path.Join(tmpDir, "stats")))

}

func TearDownStats(tmpDir string) {
	log.Debug("Clearing up stats dir")

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Error(fmt.Errorf("can not remove test path: %s \n", err))
	}

	err = os.RemoveAll(path.Join(tmpDir, "stats"))
	if err != nil {
		log.Error(fmt.Errorf("can not create test path: %s \n", err))
	}
}

func TestGetStats(t *testing.T) {
	config = viper.New()
	log.SetLevel(log.DebugLevel)

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Error(fmt.Errorf("can not remove test path: %s \n", err))
	}

	PrepareStats(tmpDir)

	err = os.Setenv("GPE_STATS_PATH", path.Join(tmpDir, "stats"))
	if err != nil { // Handle errors reading the config file
		log.Error(fmt.Errorf("can not set env var gpe_stats_path: %s \n", err))
	}

	config.SetEnvPrefix("GPE")
	err = config.BindEnv("stats_path")
	if err != nil { // Handle errors reading the config file
		log.Error(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	result, err := GetStats()
	if err != nil {
		TearDownStats(tmpDir)
		log.Debug(fmt.Errorf("something goes wrong: %s", err))
		t.Error("something goes wrong", err)
	}

	var sb strings.Builder
	sb.WriteString("gateway_some_some_1 100\n")
	sb.WriteString("gateway_some_some_2 100\n")
	sb.WriteString("gateway_some_some_3 100\n")
	sb.WriteString("gateway_stat_1 100\n")
	sb.WriteString("gateway_stat_2 100\n")
	sb.WriteString("gateway_stat_3 100\n")
	sb.WriteString("gateway_timers_timer_1 100\n")
	sb.WriteString("gateway_timers_timer_2 100\n")
	sb.WriteString("gateway_timers_timer_3 100\n")

	log.Debug(result)
	if result != sb.String() {
		t.Error("String should pass template empty")
	}

	TearDownStats(tmpDir)

}
