package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var IsInteractive bool = true

func Initialize() {
	InitDefaults()
	InitConfigPath()
}

func InitConfigPath() {
	viper.SetConfigName("gotype")        // name of config file (without extension)
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.AddConfigPath("$GOPATH/etc/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.gotype") // call multiple times to add many search paths
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}
