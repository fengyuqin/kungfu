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

package auths

import (
	"sync"
	"time"

	"github.com/fengyuqin/kungfu/v2/logger"
	"github.com/fengyuqin/kungfu/v2/stores"
	"github.com/fengyuqin/kungfu/v2/utils"
)

type IdSourcer func(key string) (int64, error)
type IdStorer func(key string, id int64) error

type IdGenerator struct {
	table     string
	key       string
	idSource  IdSourcer
	idStore   IdStorer
	lock      sync.Mutex
	nextId    int64
	storeTime time.Duration
	storeId   int64
}

type IdOption func(ig *IdGenerator)

func WithIdStoreTime(storeTime time.Duration) IdOption {
	return func(ig *IdGenerator) {
		ig.storeTime = storeTime
	}
}

func NewIdGenerator(table string, key string, idSource IdSourcer, idStore IdStorer, options ...IdOption) *IdGenerator {
	ig := &IdGenerator{
		table:     table,
		key:       key,
		idSource:  idSource,
		idStore:   idStore,
		lock:      sync.Mutex{},
		nextId:    -1,
		storeTime: 1 * time.Minute,
		storeId:   -1,
	}
	for _, option := range options {
		option(ig)
	}
	go utils.SafeRun(func() {
		ig.schedule()
	})
	return ig
}
func (ig *IdGenerator) schedule() {
	ticker := time.NewTicker(ig.storeTime)
	for range ticker.C {
		if err := ig.Store(); err != nil && !stores.IsRedisNull(err) {
			logger.Error(err)
		}
	}
}

func (g *IdGenerator) check() error {
	if !stores.HExists(g.table, g.key) {
		g.lock.Lock()
		defer g.lock.Unlock()
		id, err := g.idSource(g.key)
		if err != nil {
			logger.Reportf("id generator source err: %v", err)
			return err
		}
		if g.nextId != -1 && g.nextId > id {
			id = g.nextId
		}
		if err = stores.HSetNx(g.table, g.key, id); err != nil {
			logger.Reportf("id generator init err: %v", err)
			return err
		}
	}
	return nil
}

func (g *IdGenerator) NextId() (int64, error) {
	if err := g.check(); err != nil {
		return 0, err
	}
	nextId := stores.HIncrBy(g.table, g.key, 1)
	g.nextId = nextId
	return nextId, nil
}

func (g *IdGenerator) Store() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	var id int64
	if err := stores.HGet(g.table, g.key, &id); err != nil {
		return err
	}
	if g.storeId == id {
		return nil
	}
	logger.Infof("id generator store:%d", id)
	err := g.idStore(g.key, id)
	if err == nil {
		g.storeId = id
	}
	return err
}
