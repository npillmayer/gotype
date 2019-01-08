package viperadapter

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config"
	"github.com/spf13/viper"
)

type VConf struct{}

func New() *VConf {
	return &VConf{}
}

func (c *VConf) Init() {
	InitDefaults()
	InitConfigPath()
}

// InitDefaults is usually called by Init().
func InitDefaults() {
	/*
		viper.SetDefault("ContentDir", "content")
		viper.SetDefault("LayoutDir", "layouts")
		viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
		viper.Set("Verbose", true)
	*/

	viper.SetDefault("tracing", "logrus")
	viper.SetDefault("tracingonline", true)
	viper.SetDefault("tracingequations", "Error")
	viper.SetDefault("tracingsyntax", "Error")
	viper.SetDefault("tracingcommands", "Error")
	viper.SetDefault("tracinginterpreter", "Error")
	viper.SetDefault("tracinggraphics", "Error")

	viper.SetDefault("tracingcapsules", "Error")
	viper.SetDefault("tracingrestores", "Error")
	viper.SetDefault("tracingchoices", true)
}

// InitConfigPath is usually called by Init().
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

// IsSet is a predicate wether a configuration flag is set to true.
func (c *VConf) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (c *VConf) GetString(key string) string {
	return viper.GetString(key)
}

func (c *VConf) GetInt(key string) int {
	return viper.GetInt(key)
}

func (c *VConf) GetBool(key string) bool {
	return viper.GetBool(key)
}

var _ config.Configuration = &VConf{}
