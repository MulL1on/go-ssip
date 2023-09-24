package config

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Key  string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	Name         string       `mapstructure:"name" json:"name"`
	Host         string       `mapstructure:"host" json:"host"`
	Port         int          `mapstructure:"port" json:"port"`
	PasetoInfo   PasetoConfig `mapstructure:"paseto" json:"paseto"`
	OtelInfo     OtelConfig   `mapstructure:"otel" json:"otel"`
	UserSrvInfo  RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	RedisInfo    RedisConfig  `mapstructure:"redis" json:"redis"`
	GroupSrvInfo RPCSrvConfig `mapstructure:"group_srv" json:"group_srv"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type PasetoConfig struct {
	PubKey   string `mapstructure:"pub_key" json:"pub_key"`
	Implicit string `mapstructure:"implicit" json:"implicit"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}
