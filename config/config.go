package config

type ConfigRegistry struct {
	App      AppConfig
	Database DatabaseConfig
	Log      LoggerConfig
}
