package consts

const (
	UserServiceName = "user"
	MsgServiceName  = "msg"
	HttpApiName     = "http_api"
	WsApiName       = "ws_api"
	TransMqName     = "trans_mq"
	PullMqName      = "pull_mq"

	FreePortAddress = "localhost:0"
	TCP             = "tcp"

	UserConfigPath    = "./app/service/rpc/user/config.yaml"
	MsgConfigPath     = "app/service/rpc/msg/config.yaml"
	HttpApiConfigPath = "app/service/api/http/config.yaml"
	WsApiConfigPath   = "app/service/api/ws/config.yaml"
	TransMqConfigPath = "app/service/mq/trans/config.yaml"
	PullMqConfigPath  = "app/service/mq/pull/config.yaml"

	ConsulCheckInterval                       = "7s"
	ConsulCheckTimeout                        = "5s"
	ConsulCheckDeregisterCriticalServiceAfter = "15s"

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	MongoURI    = "mongodb://%s:%s@%s:%d"
	MysqlDSN    = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	RabbitMqUrl = "amqp://%s:%s@%s:%d/"

	Issuer = "go-ssip"

	UserID = "uid"
	User   = "user"

	MessageTypeText  = 0
	MessageTypeImage = 1
)
