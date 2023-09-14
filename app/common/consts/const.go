package consts

const (
	UserServiceName = "user"
	MsgServiceName  = "msg"

	FreePortAddress = "localhost:0"
	TCP             = "tcp"

	UserConfigPath    = "app/service/rpc/user/config.yaml"
	HtppApiConfigPath = "app/service/api/http/config.yaml"
	WsApiConfigPath   = "app/service/api/ws/config.yaml"

	ConsulCheckInterval                       = "7s"
	ConsulCheckTimeout                        = "5s"
	ConsulCheckDeregisterCriticalServiceAfter = "15s"

	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"

	MysqlDSN = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	Issuer = "tank_war"

	UserID = "uid"
	User   = "user"
)
