package tests

import (
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/stores"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("package", "tests")
)

func init() {
	if err := config.InitConf("../config.json"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(config.GetDiscoverConf())
	//init stores
	stores.InitStoreKeeper(config.GetStoresConf())
	//init coder
	coder.InitCoder(config.GetCoderConf())
}
