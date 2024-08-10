package configs

type Configs struct {
	DRIVER				string `mapstructure:"DRIVER"`
	DB_USER             string `mapstructure:"DB_USER"`
	DB_PASSWORD         string `mapstructure:"DB_PASSWORD"`
	DB_NAME             string `mapstructure:"DB_NAME"`
	DB_PORT             string `mapstructure:"DB_PORT"`
	HOST_DB             string `mapstructure:"HOST_DB"`
	SSLmode             string `mapstructure:"SSLmode"`
	SERVER_PORT         string `mapstructure:"SERVER_PORT"`
	GOOSE_MIGRATION_DIR string `mapstructure:"GOOSE_MIGRATION_DIR"`
	GOOSE_DRIVER        string `mapstructure:"GOOSE_DRIVER"`
	GOOSE_DBSTRING      string `mapstructure:"GOOSE_DBSTRING"`
}
