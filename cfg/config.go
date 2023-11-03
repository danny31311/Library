package cfg

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Cfg struct {
	Port   string
	DbName string
	DbUser string
	DbPass string
	DbHost string
	DbPort string
}

func LoadAndStoreConfig() Cfg {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables %s", err.Error())
	}
	v := viper.New()
	v.SetEnvPrefix("SERV")
	v.SetDefault("PORT", "8080")
	v.SetDefault("DBUSER", "postgres")
	v.SetDefault("DBPASS", os.Getenv("DB_PASSWORD"))
	v.SetDefault("DBHOST", "")
	v.SetDefault("DBPORT", "5432")
	v.Set("DBNAME", "go_library")
	v.AutomaticEnv()
	var cfg Cfg
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Panic(err)
	}
	return cfg
}

func (cfg *Cfg) GetDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
}
