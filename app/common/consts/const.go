package consts

const (
	UserServiceName  = "user_srv"
	MsgServiceName   = "msg_srv"
	GroupServiceName = "group_srv"
	HttpApiName      = "http_api"
	WsApiName        = "ws_api"
	TransMqName      = "trans_mq"
	PullMqName       = "pull_mq"

	FreePortAddress = "localhost:0"
	TCP             = "tcp"

	UserSrvConfigPath  = "app/service/rpc/user/config.yaml"
	MsgSrvConfigPath   = "app/service/rpc/msg/config.yaml"
	GroupSrvConfigPath = "app/service/rpc/group/config.yaml"
	HttpApiConfigPath  = "app/service/api/http/config.yaml"
	WsApiConfigPath    = "app/service/api/ws/config.yaml"
	TransMqConfigPath  = "app/service/mq/transfer/config.yaml"
	PullMqConfigPath   = "app/service/mq/pull/config.yaml"

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
)

const (
	CommandTypeAckMsg uint32 = iota
	CommandTypeSendMsg
	CommandTypeAckClientId
	CommandTypeGetMsg
)

const (
	MessageTypeSingleChat int8 = iota
	MessageTypeGroupChat
)
