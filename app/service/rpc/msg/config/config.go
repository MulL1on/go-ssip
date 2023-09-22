package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Salt     string `mapstructure:"salt" json:"salt"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Key  string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Port         int            `mapstructure:"port" json:"port"`
	MysqlInfo    MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	OtelInfo     OtelConfig     `mapstructure:"otel" json:"otel"`
	PasetoInfo   PasetoConfig   `mapstructure:"paseto" json:"paseto"`
	JwtInfo      JwtConfig      `mapstructure:"jwt" json:"jwt"`
	LoggerInfo   LoggerConfig   `mapstructure:"logger" json:"logger"`
	RedisInfo    RedisConfig    `mapstructure:"redis" json:"redis"`
	RabbitMQInfo RabbitMQConfig `mapstructure:"rabbitmq" json:"rabbitmq"`
	KafkaInfo    KafkaConfig    `mapstructure:"kafka" json:"kafka"`
}

type PasetoConfig struct {
	SecretKey string `mapstructure:"secret_key" json:"secret_key"`
	Implicit  string `mapstructure:"implicit" json:"implicit"`
}

type JwtConfig struct {
	SecretKey   string `mapstructure:"secretKey" yaml:"secretKey"`
	ExpiresTime int64  `mapstructure:"expiresTime" yaml:"expiresTime"`
	Issuer      string `mapstructure:"issuer" yaml:"issuer"`
}

type LoggerConfig struct {
	SavePath   string `mapstructure:"savePath" yaml:"savePath"`
	MaxSize    int    `mapstructure:"maxSize" yaml:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge" yaml:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups"`
	IsCompress bool   `mapstructure:"isCompress" yaml:"isCompress"`
	LogLevel   string `mapstructure:"logLevel" yaml:"logLevel"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

type KafkaConfig struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    int    `mapstructure:"port" json:"port"`
	Timeout int    `mapstructure:"timeout" json:"timeout"`
}
