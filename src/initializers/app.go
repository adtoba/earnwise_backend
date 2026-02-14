package initializers

import "github.com/spf13/viper"

type Config struct {
	DBUsername    string `mapstructure:"POSTGRES_USER"`
	DBPassword    string `mapstructure:"POSTGRES_PASSWORD"`
	DBHost        string `mapstructure:"POSTGRES_HOST"`
	DBPort        string `mapstructure:"POSTGRES_PORT"`
	DBName        string `mapstructure:"POSTGRES_DB"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	RedisAddr     string `mapstructure:"REDIS_ADDRESS"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	Port          string `mapstructure:"PORT"`
	JWT_SECRET    string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
