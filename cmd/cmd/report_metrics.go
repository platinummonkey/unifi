package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/platinummonkey/unifi/cmd/reporters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// reporterCmd represents the reporter command
var reporterCmd = &cobra.Command{
	Use:   "reporter",
	Short: "Run this as a daemon to collect and report metrics",
	Long: `Run this application as a daemon to periodically collect
and report metrics to the defined output(s).`,
	Run: runDaemon,
}

func init() {
	rootCmd.AddCommand(reporterCmd)
	reporterCmd.Flags().StringSliceP("outputs", "o", []string{}, "Specify a subeset of configured outputs, leave unspecified for all configured outputs")
	reporterCmd.Flags().StringSliceP("sites", "s", []string{}, "Specify a subeset of configured sites, leave unspecified for all configured sites")
	reporterCmd.Flags().IntP("workers", "w", 1, "Specify the number of concurrent collections")
	viper.SetDefault("reporter.frequency", time.Second*1)
	viper.SetDefault("workers", 1)
	viper.SetDefault("state_file", "")
}

func getReporterOutputs(cmd *cobra.Command) (map[string]reporters.Reporter, error) {
	filterOutputs, err := cmd.Flags().GetStringSlice("only-outputs")
	if err != nil {
		filterOutputs = []string{}
	}
	limitOutputs := make(map[string]struct{}, 0)
	for _, o := range filterOutputs {
		limitOutputs[o] = struct{}{}
	}

	outputs := viper.GetStringMap("reporter.outputs")
	if len(outputs) == 0 {
		logger.Error("must specify outputs either via --outputs or in the config file")
		return nil, fmt.Errorf("invalid `outputs`")
	}

	reportOutputs := make(map[string]reporters.Reporter, 0)
	if len(limitOutputs) > 0 {
		for o, v := range outputs {
			if _, ok := limitOutputs[o]; ok {
				r, err := reporters.ReporterFromTypeAndConfig(o, v.(map[string]interface{}))
				if err == nil {
					reportOutputs[o] = r
				}
			}
		}
		return reportOutputs, nil
	} else {
		for o, v := range outputs {
			r, err := reporters.ReporterFromTypeAndConfig(o, v.(map[string]interface{}))
			if err == nil {
				reportOutputs[o] = r
			}
		}
	}

	return reportOutputs, nil
}

type siteConfig struct {
	Tags []string `json:"tags"`
	Name string `json:"name"`
	ID string
	UUID string
}


func getConfiguredSites(cmd *cobra.Command) map[string]siteConfig {
	configuredSites := make(map[string]siteConfig, 0)
	filterSites, err := cmd.Flags().GetStringSlice("sites")
	hasFilteredSites := err == nil && len(filterSites) > 0

	in := func(l []string, t string) bool {
		for _, i := range l {
			if t == i {
				return true
			}
		}
		return false
	}

	// build entirely
	// use them all by default
	sites := viper.GetStringMap("sites")
	if len(sites) == 0 {
		// by default this behavior is all in filtered
		if hasFilteredSites {
			for _, siteID := range filterSites {
				configuredSites[siteID] = siteConfig{
					ID: siteID,
					Name: siteID,
				}
			}
		}
		return configuredSites
	}

	for siteID, siteCfgRaw := range sites {
		var siteCfg siteConfig
		d, err := json.Marshal(siteCfgRaw)
		if err != nil {
			if hasFilteredSites && !in(filterSites, siteID) {
				continue
			}
			configuredSites[siteID] = siteConfig{
				ID: siteID,
				Name: siteID,
			}
			continue
		}
		err = json.Unmarshal(d, &siteCfg)
		if err != nil {
			if hasFilteredSites && !in(filterSites, siteID) {
				continue
			}
			configuredSites[siteID] = siteConfig{
				ID: siteID,
				Name: siteID,
			}
			continue
		}
		if siteCfg.Name == "" {
			siteCfg.Name = siteID
		}
		siteCfg.ID = siteID
		if hasFilteredSites && !in(filterSites, siteCfg.Name) {
			continue
		}
		configuredSites[siteID] = siteCfg
	}

	return configuredSites
}

func getWorkerCount(cmd *cobra.Command) int {
	workerCount, err := cmd.Flags().GetInt("workers")
	if err == nil && workerCount > 0 {
		return workerCount
	}
	workerCount = viper.GetInt("workers")
	if workerCount > 0 {
		return workerCount
	}
	return 1
}

func collectAndReportStats(workerCount int, reporters reporters.Reporters, siteCfg map[string]siteConfig) {
	sites, err := client.AvailableSites()
	if err != nil {
		logger.Error("unable to query sites", zap.Error(err))
		return
	}
	shouldFilterSites := len(siteCfg) != 0
	sitesToCollectStats := make([]siteConfig, 0)
	sitesToQuery := make([]string, 0)
	for _, site := range sites.Data {
		sc, ok := siteCfg[site.Name]
		if shouldFilterSites && !ok {
			logger.Debug("ignoring site as it was not found in configured sites", zap.String("site", site.Name))
		} else if !ok {
			// set the default
			sc = siteConfig{ID: site.Name, Name: site.Name, UUID: site.ID}
		} else {
			sc.UUID = site.ID
		}
		sitesToCollectStats = append(sitesToCollectStats, sc)
		sitesToQuery = append(sitesToQuery, sc.Name)
	}

	logger.Debug("will collect and report stats for", zap.Strings("sites", sitesToQuery))

	workChan := make(chan siteConfig, len(sitesToCollectStats))
	for _, sc := range sitesToCollectStats {
		workChan <- sc
	}

	// fire up workers to collect site stats and report
	wg := &sync.WaitGroup{}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go workerCollectSiteStats(wg, workChan, reporters)
	}
	wg.Wait()
}

func periodicallyReportStats(stopChan chan bool, workerCount int, reporters reporters.Reporters, siteCfg map[string]siteConfig) {
	for {
		select {
		case <-time.After(viper.GetDuration("reporter.frequency")):
			collectAndReportStats(workerCount, reporters, siteCfg)
		case <-stopChan:
			stopChan <- true
		}
	}
}

func runDaemon(cmd *cobra.Command, args []string) {
	outputs, err := getReporterOutputs(cmd)
	if err != nil {
		os.Exit(1)
	}

	reporters := make([]string, 0)
	for output := range outputs {
		reporters = append(reporters, output)
	}
	configuredSites := getConfiguredSites(cmd)
	sitesToQuery := make([]string, 0)
	for _, cfg := range configuredSites {
		sitesToQuery = append(sitesToQuery, cfg.Name)
	}

	initReporterState()

	logger.Debug("Starting daemon with", zap.Strings("reporters", reporters), zap.Strings("sites", sitesToQuery))

	sigChan := make(chan os.Signal, 1)
	exitChan := make(chan bool, 1)
	stopChan := make(chan bool, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			sig := <- sigChan
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				// do a clean stop
				stopChan <- true
				<-stopChan
				// exit
				exitChan <- true
			}
		}
	}()

	go periodicallyReportStats(stopChan, getWorkerCount(cmd), outputs, configuredSites)

	<- exitChan
	db.db.Close()
}
