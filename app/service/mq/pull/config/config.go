package config

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Key  string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Port         int            `mapstructure:"port" json:"port"`
	OtelInfo     OtelConfig     `mapstructure:"otel" json:"otel"`
	RabbitMQInfo RabbitMQConfig `mapstructure:"rabbitmq" json:"rabbitmq"`
	MysqlInfo    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	MongoDBInfo  MongoDBConfig  `mapstructure:"mongodb" json:"mongodb"`
	LoggerInfo   LoggerConfig   `mapstructure:"logger" json:"logger"`
	RedisInfo    RedisConfig    `mapstructure:"redis" json:"redis"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

type LoggerConfig struct {
	SavePath   string `mapstructure:"savePath" yaml:"savePath"`
	MaxSize    int    `mapstructure:"maxSize" yaml:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge" yaml:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups"`
	IsCompress bool   `mapstructure:"isCompress" yaml:"isCompress"`
	LogLevel   string `mapstructure:"logLevel" yaml:"logLevel"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Salt     string `mapstructure:"salt" json:"salt"`
}

type MongoDBConfig struct {
	Host       string `mapstructure:"host" json:"host"`
	Port       int    `mapstructure:"port" json:"port"`
	Name       string `mapstructure:"db" json:"db"`
	User       string `mapstructure:"user" json:"user"`
	Password   string `mapstructure:"password" json:"password"`
	Collection string `mapstructure:"collection" json:"collection"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}
