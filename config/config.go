package config

type MongoConfig struct {
	Username string
	Password string
	Endpoint string
	Name     string
}

type GeneralConfig struct {
	Mongo      MongoConfig
	ServerHost string
	JwtKey     string
}
