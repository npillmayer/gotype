package config

import "github.com/spf13/viper"

// InitDefaults is usually called by Initialize().
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
