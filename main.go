package main

import (
	"time"

	"github.com/anirudhgray/mood-harbour-backend/config"
	"github.com/anirudhgray/mood-harbour-backend/infra/database"
	"github.com/anirudhgray/mood-harbour-backend/infra/logger"
	"github.com/anirudhgray/mood-harbour-backend/migrations"
	"github.com/anirudhgray/mood-harbour-backend/routers"
	"github.com/spf13/viper"
)

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Calcutta")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	config.InitialiseOAuthGoogle()

	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	//later separate migration
	migrations.Migrate()

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
