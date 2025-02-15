package config

var DefaultConfig = Config{
	Database: DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DBName:   "myapp",
	},
	Server: ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
	},
	Log: LogConfig{
		Level: "info",
	},
}
