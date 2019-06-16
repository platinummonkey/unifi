package cmd

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/platinummonkey/unifi"
	"github.com/platinummonkey/unifi/cmd/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string
var logger *zap.Logger
var client *unifi.Client
var homeDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "unifi",
	Short: "unifi is a client to the UniFi controller APIs",
	Long: `unifi provides a client interface to the UniFI controller APIs.

	This allows you to create custom applications around these APIs,
    hook in configuration management to these devices,
    collect stats to report to your metric collection systems, or
    whatever unique application comes to mind.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if logger != nil {
		defer logger.Sync()
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	viper.SetEnvPrefix("UNIFI")
	cobra.OnInitialize(initConfig)

	// global persistent configuration
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.unifi.yaml)")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "warn", "loglevel")
	rootCmd.PersistentFlags().StringP("baseurl", "b", "", "BaseURL for the controller API. Ex. https://1.2.3.4:8443/")
	rootCmd.PersistentFlags().StringP("username", "u", "", "API username for the controller")
	rootCmd.PersistentFlags().StringP("password", "p", "", "API user password for the controller")
	// http settings
	rootCmd.PersistentFlags().DurationP("timeout", "t", time.Second*30, "API timeout duration")
	rootCmd.PersistentFlags().BoolP("disableTLS", "k", false, "Disable TLS checks on http client")
	rootCmd.PersistentFlags().StringSlice("x509certs", []string{}, "Specify a list of x509 certificates for the http client trust")
	rootCmd.PersistentFlags().String("pemCert", "", "Specify a PEM cert file for the http client trust")

	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("disableTLS", rootCmd.PersistentFlags().Lookup("disableTLS"))
	viper.BindPFlag("x509certs", rootCmd.PersistentFlags().Lookup("x509certs"))
	viper.BindPFlag("pemCert", rootCmd.PersistentFlags().Lookup("pemCert"))
	viper.SetDefault("loglevel", "warn")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	var err error
	homeDir, err = homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".unifi")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// no such config file, create a default empty one.
		defaultFilePath := path.Join(homeDir, ".unifi.yaml")
		if cfgFile == "" && strings.Contains(strings.ToLower(err.Error()), "not found in") {
			err = ioutil.WriteFile(defaultFilePath, defaultConfigData, 0600)
			if err != nil {
				fmt.Printf("unable to create default config file: %s\n", defaultFilePath)
				os.Exit(1)
			} else {
				fmt.Printf("no config file found, creating default in: %s\n", defaultFilePath)
				os.Exit(1)
			}
		}
		fmt.Printf("Unable to read config file %s: %s\n", viper.ConfigFileUsed(), err)
		os.Exit(1)
	}

	// initialize logger
	logger, err = log.InitializeLogger(viper.GetString("loglevel"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// initialize the client
	baseURLStr := viper.GetString("baseurl")
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		logger.Error("invalid baseURL", zap.String("baseurl", baseURLStr), zap.Error(err))
	}
	certConfig := unifi.CertificationConfig{
		DisableCertCheck: viper.GetBool("disableTLS"),
	}
	x509Certs := viper.GetStringSlice("x509certs")
	if len(x509Certs) > 0 {
		for i, certFile := range x509Certs {
			certData, err := ioutil.ReadFile(certFile)
			if err != nil {
				logger.Warn("skipping invalid cert file", zap.String("file", certFile), zap.Error(err))
				continue
			}
			cert, err := x509.ParseCertificate(certData)
			if err != nil {
				logger.Warn("skipping invalid cert file", zap.String("file", certFile), zap.Error(err))
				continue
			}
			if i == 0 {
				certConfig.Certificates = make([]*x509.Certificate, 0)
			}
			certConfig.Certificates = append(certConfig.Certificates, cert)
		}
	}
	if viper.GetString("pemCert") != "" {
		certConfig.PEMCert = viper.GetString("pemCert")
	}
	client, err = unifi.NewClient(baseURL.String(), &certConfig, viper.GetDuration("timeout"))
	if err != nil {
		logger.Error("unable to initialize client", zap.Error(err))
	} else {
		logger.Debug("initialized client")
	}
	err = client.Login(viper.GetString("username"), viper.GetString("password"), false)
	if err != nil {
		logger.Error("unable to authenticate against controller", zap.Error(err))
	} else {
		logger.Debug("successfully authenticated against controller")
	}
}

var defaultConfigData = []byte(`# Autogenerated Default configuration

# loglevel options are: debug, info, warn, error, fatal, none
loglevel: warn
# base URL of the controller API, like https://1.2.3.4:8443/
baseurl: FIXME
timeout: 30s
# username: or specify via the UNIFI_USERNAME env variable 
# password: or specify via the UNIFI_PASSWORD env variable
workers: 1
# state_dir: "" # set the default state directory, by default is ~/.unifi.state/

sites: {}

reporter:
  frequency: 1s
  outputs:
    # Datadog API reporter
	#datadog:
	#  api_key: XXXXX
	# Dogstatsd reporter
	#dogstatsd:
	#  endpoint: # see dogstatsd location, examples are localhost:6379, unix://var/run/datadog.sock
	# STDOUT/Log format
	log: {}

`)
