package config

import (
	"log"

	"github.com/spf13/viper"
)

var Env *EnvConfig

type EnvConfig struct {
	ECRRegistryName   string `mapstructure:"ECR_REGISTRY_NAME"`
	ECRRepositoryName string `mapstructure:"ECR_REPOSITORY_NAME"`
	CDKStackID        string `mapstructure:"CDK_STACK_ID"`
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}

	if err := viper.Unmarshal(&Env); err != nil {
		log.Fatal(err)
	}
}
