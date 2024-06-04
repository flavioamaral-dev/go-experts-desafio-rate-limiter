package configs

import "github.com/spf13/viper"

type conf struct {
	LimitByIPMaxRPS         int64  `mapstructure:"LIMIT_BY_IP_MAX_RPS"`
	LimitByIPBlockTimeMs    int64  `mapstructure:"LIMIT_BY_IP_BLOCK_TIME_MS"`
	LimitByTokenMaxRPS      int64  `mapstructure:"LIMIT_BY_TOKEN_MAX_RPS"`
	LimitByTokenBlockTimeMs int64  `mapstructure:"LIMIT_BY_TOKEN_BLOCK_TIME_MS"`
	WebServerPort           string `mapstructure:"WEB_SERVER_PORT"`
	RedisAddr               string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
