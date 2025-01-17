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
	"github.com/fengyuqin/kungfu/v2/logger"
	"github.com/fengyuqin/kungfu/v2/utils"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"testing"
	"time"

	"github.com/fengyuqin/kungfu/v2/stores"
)

func TestStores(t *testing.T) {
	err := stores.Set("name", "jason", 5*time.Second)
	if err != nil {
		logger.Error(err)
		return
	}
	var res string
	err = stores.Get("name", &res)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Get name res:%+v", res)
	res2 := stores.GetString("name")
	logger.Infof("Get name res2:%+v", res2)

	res3 := stores.GetInt("name")
	logger.Infof("Get name res3:%+v", res3)
}

func TestStoreList(t *testing.T) {
	key := "myList"
	var err error
	err = stores.LPush(key, 1, 2, 3, 4, 5)
	if err != nil {
		logger.Error(err)
		return
	}
	fmt.Println("length:", stores.LLen(key))
	var a string
	if a, err = stores.BRPopString(key); err != nil {
		logger.Error(err)
		return
	}
	fmt.Println("pop:", a)
	fmt.Println("length:", stores.LLen(key))

	if err = stores.BRPop(key, &a); err != nil {
		logger.Error(err)
		return
	}
	fmt.Println("pop:", a)
	fmt.Println("length:", stores.LLen(key))
	select {}
}

func TestZAdd(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	maxNum := 100000
	list := make([]*redis.Z, 0)
	for i := 0; i < maxNum; i++ {
		//score, member := float64(rand.Intn(maxNum)), fmt.Sprintf("member_%d", i)
		score, member := float64(rand.Intn(maxNum)), utils.IntToString(i+1)
		list = append(list, &redis.Z{
			Score:  score,
			Member: member,
		})
	}
	if err := stores.ZAdd("ListRank", list...); err != nil {
		fmt.Println(err)
	}
	//fmt.Println("ZAdd成功")
	//list, err := stores.ZRevRangeWithScores("list_rank", 0, 100)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(list)
	//fmt.Println(len(list))
	//if err := stores.ZRem("list_rank", utils.IntToString(6112)); err != nil {
	//	fmt.Println("111", err)
	//}
	//if rank, err := stores.ZRevRank("list_rank", utils.IntToString(6113)); err != nil {
	//	fmt.Println("222", err)
	//} else {
	//	fmt.Println(rank)
	//}
	//total := stores.ZCard("list_rank")
	//fmt.Println(total)
	//if s, err := stores.ZIncrBy("list_rank", 3, utils.IntToString(6113)); err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(s)
	//}
	//score := stores.ZScore("list_rank", utils.IntToString(6113))
	//fmt.Println(score)
	//if v, err := stores.ZRangeWithScores("list_rank", 0, 0); err != nil {
	//	logger.Error(err)
	//} else {
	//	for _, item := range v {
	//		fmt.Println(item.Member, item.Score)
	//	}
	//}
}
