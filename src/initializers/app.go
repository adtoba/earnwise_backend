package initializers

import "github.com/spf13/viper"

type Config struct {
	DBUsername          string `mapstructure:"POSTGRES_USER"`
	DBPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	DBHost              string `mapstructure:"POSTGRES_HOST"`
	DBPort              string `mapstructure:"POSTGRES_PORT"`
	DBName              string `mapstructure:"POSTGRES_DB"`
	DBEndpointID        string `mapstructure:"POSTGRES_DB_ENDPOINT"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	RedisAddr           string `mapstructure:"REDIS_ADDRESS"`
	RedisUsername       string `mapstructure:"REDIS_USERNAME"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	RedisDB             int    `mapstructure:"REDIS_DB"`
	SSLMode             string `mapstructure:"SSL_MODE"`
	OneSignalAppID      string `mapstructure:"ONESIGNAL_APP_ID"`
	OneSignalAPIKey     string `mapstructure:"ONESIGNAL_APP_API_KEY"`
	Port                string `mapstructure:"PORT"`
	AgoraAppID          string `mapstructure:"AGORA_APP_ID"`
	AgoraAppCertificate string `mapstructure:"AGORA_APP_CERT"`
	GoogleWebClientID   string `mapstructure:"GOOGLE_WEB_CLIENT_ID"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// viper.AutomaticEnv()

	// // Explicitly bind environment variables
	// viper.BindEnv("POSTGRES_USER")
	// viper.BindEnv("POSTGRES_PASSWORD")
	// viper.BindEnv("POSTGRES_HOST")
	// viper.BindEnv("POSTGRES_PORT")
	// viper.BindEnv("POSTGRES_DB")
	// viper.BindEnv("POSTGRES_DB_ENDPOINT")
	// viper.BindEnv("JWT_SECRET")
	// viper.BindEnv("REDIS_ADDRESS")
	// viper.BindEnv("REDIS_USERNAME")
	// viper.BindEnv("REDIS_PASSWORD")
	// viper.BindEnv("REDIS_DB")
	// viper.BindEnv("ONESIGNAL_APP_ID")
	// viper.BindEnv("ONESIGNAL_APP_API_KEY")
	// viper.BindEnv("PORT")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
