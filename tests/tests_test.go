/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package tests

import (
	"fmt"
	"github.com/fengyuqin/kungfu/v2/config"
	"github.com/fengyuqin/kungfu/v2/discover"
	"github.com/fengyuqin/kungfu/v2/logger"
	"github.com/fengyuqin/kungfu/v2/rpc"
	"github.com/fengyuqin/kungfu/v2/stores"
	"github.com/fengyuqin/kungfu/v2/treaty"
	"testing"
)

func init() {
	if err := config.InitConf("../examples/nano_demo/config.json"); err != nil {
		logger.Fatal(err)
	}
	//init discover
	discover.InitDiscoverer(config.GetDiscoverConf())
	//init stores
	stores.InitStoreKeeper(config.GetStoresConf())
}

func TestSub(t *testing.T) {
	s := rpc.NewRssBuilder(&treaty.Server{ServerId: "test_001"})
	fmt.Println(s)
	s1 := s.Build()
	s2 := s1.SetSuffix("wel").Build()
	fmt.Println(s, s1, s2)
}
