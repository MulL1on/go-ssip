package initialize

import (
	"flag"
	"go-ssip/app/common/consts"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/group/global"
	"go.uber.org/zap"
)

func InitFlag() (string, int) {
	IP := flag.String(consts.IPFlagName, consts.IPFlagValue, consts.IPFlagUsage)
	Port := flag.Int(consts.PortFlagName, 0, consts.PortFlagUsage)

	//parsing flag, and if Port is 0, will automatically get an empty Port
	flag.Parse()
	if *Port == 0 {
		*Port, _ = tools.GetFreePort()
	}
	g.Logger.Info("parse flag successfully", zap.String("ip", *IP), zap.Int("port", *Port))
	return *IP, *Port
}
