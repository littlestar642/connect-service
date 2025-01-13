package config

import "github.com/spf13/viper"

type Config struct{
	Port string
	RedisAddr string
	KafkaAddr string
}

var appConfig *Config

func Load() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("KAFKA_ADDR", "localhost:9092")

	viper.AutomaticEnv()

	appConfig = &Config{
		Port: viper.GetString("PORT"),
		RedisAddr: viper.GetString("REDIS_ADDR"),
		KafkaAddr: viper.GetString("KAFKA_ADDR"),
	}
}

func New() *Config {
	if(appConfig == nil){
		Load()
	}
	return appConfig
}