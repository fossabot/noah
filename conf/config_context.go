package conf

import (
	"github.com/Ready-Stock/Noah/db/system"
)

type ConfigContext struct {
	*system.SContext
}

func (ctx *ConfigContext) Cluster() *ClusterConfigContext {
	return &ClusterConfigContext{ctx}
}


