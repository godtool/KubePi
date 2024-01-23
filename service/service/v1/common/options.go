package common

import (
	"github.com/asdine/storm/v3"
	"github.com/godtool/kubeone/service/server"
)

type DBService interface {
	GetDB(options DBOptions) storm.Node
}

var defaultStormNode storm.Node

func SetDefaultStormNode(n storm.Node) {
	defaultStormNode = n
}

type DefaultDBService struct {
}

func (d *DefaultDBService) GetDB(options DBOptions) storm.Node {
	if options.DB != nil {
		return options.DB
	}
	if defaultStormNode != nil {
		return defaultStormNode
	}

	return server.DB()
}

type DBOptions struct {
	DB storm.Node
}
