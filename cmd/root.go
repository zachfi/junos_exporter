package cmd

import (
	"fmt"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xaque208/junos_exporter/exporter"
)

var rootCmd = &cobra.Command{
	Use:   "junos_exporter",
	Short: "Export Junos ARP data Pometheus",
	Long:  "",
	Run:   run,
}

var (
	verbose       bool
	cfgFile       string
	listenAddress string
	interval      int
	junosUsername string
	junosPassword string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.junos_exporter.yaml)")
	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen", "L", ":9100", "The listen address (default is :9100")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 30, "The interval at which to update the data")
	rootCmd.PersistentFlags().StringVarP(&junosUsername, "username", "", "", "The Junos username")
	rootCmd.PersistentFlags().StringVarP(&junosPassword, "password", "", "", "The Junos password")

	rootCmd.MarkPersistentFlagRequired("config")

	viper.SetDefault("interval", 30)
}

// initConfig reads in the config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".junos_exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".junos_exporter")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
		cfgFile = viper.ConfigFileUsed()
	}
}

func run(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	hosts := viper.GetStringSlice("junos.hosts")
	if len(hosts) == 0 {
		log.Fatal("List of hosts is required")
	}

	auth := &junos.AuthMethod{
		Username:   viper.GetString("junos.username"),
		PrivateKey: viper.GetString("junos.keyfile"),
	}

	log.Infof("Starting prometheus HTTP metrics server: %s", listenAddress)
	go exporter.StartMetricsServer(listenAddress)

	interval = viper.GetInt("interval")
	log.Debugf("Tick interval: %d", interval)
	for range time.Tick(time.Duration(interval) * time.Second) {
		log.Debug("Scraping metrics from Junos")
		exporter.ScrapeMetrics(auth, hosts)
	}
}
