package config

import "github.com/spf13/viper"

func InitDefaults() {
	/*
		viper.SetDefault("ContentDir", "content")
		viper.SetDefault("LayoutDir", "layouts")
		viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
		viper.Set("Verbose", true)
	*/

	viper.SetDefault("tracingonline", true)
	viper.SetDefault("tracingequations", "WARN")
	viper.SetDefault("tracingsyntax", "WARN")
	viper.SetDefault("tracingcommands", "WARN")
	viper.SetDefault("tracinginterpreter", "WARN")

	viper.SetDefault("tracingcapsules", "WARN")
	viper.SetDefault("tracingrestores", "WARN")
}
