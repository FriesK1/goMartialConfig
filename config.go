package goMartialConfig

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("config.filename", "config")
	viper.SetDefault("config.format", "yaml")
}

func SetAppVersion(verStr string) {
	var appVersion *semver.Version
	var err error

	appVersion, err = semver.NewVersion(verStr)
	if err != nil {
		log.Fatalln("Error setting app version")
	}

	viper.SetDefault("app.version", appVersion)
}

func SetAppName(appName, appPrefix string) {
	viper.SetDefault("app.name", appName)
	viper.SetDefault("app.prefix", appPrefix)
}

func GetConfigs() {
	var helpRequested bool
	var versionRequested bool

	viper.SetConfigName(viper.GetString("config.filename"))
	viper.SetConfigType(viper.GetString("config.format"))

	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", viper.GetString("app.prefix")))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", viper.GetString("app.prefix")))
	viper.AddConfigPath("./etc/")
	viper.WatchConfig()

	viper.SetEnvPrefix(viper.GetString("app.prefix"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	pflag.BoolVarP(&helpRequested, "help", "h", false, "This message")
	pflag.BoolVarP(&versionRequested, "version", "v", false, "Get Program Version")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			log.Fatal("invalid config file format")
		}
	}

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	if helpRequested {
		showHelp()
	}
	if versionRequested {
		showVersion()
	}
}

func showHelp() {
	fmt.Printf("Usage: %s:\n", os.Args[0])
	pflag.PrintDefaults()

	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%s: %s\n", os.Args[0], viper.Get("app.version"))

	os.Exit(0)
}
