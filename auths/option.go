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
	"github.com/farmerx/gorsa"
	"github.com/fengyuqin/kungfu/v2/logger"
	"github.com/fengyuqin/kungfu/v2/serialize"
)

type Option func(e *Encipherer)

func WithRsaPubKey(pubKey string) Option {
	return func(e *Encipherer) {
		if err := gorsa.RSA.SetPublicKey(pubKey); err != nil {
			logger.Fatalf("set rsa public key err:%v", err)
		}
		e.rsaPubKey = pubKey
	}
}

func WithRsaPriKey(priKey string) Option {
	return func(e *Encipherer) {
		if err := gorsa.RSA.SetPrivateKey(priKey); err != nil {
			logger.Fatalf("set rsa pri key err:%v", err)
		}
		e.rsaPriKey = priKey
	}
}

func WithAesKey(aesKey string) Option {
	return func(e *Encipherer) {
		e.aesKey = []byte(aesKey)
	}
}

func WithAesIv(aesIv string) Option {
	return func(e *Encipherer) {
		e.aesIv = []byte(aesIv)
	}
}

func WithUnencrypt(val bool) Option {
	return func(e *Encipherer) {
		e.unencrypt = val
	}
}

func WithSerializer(serializer serialize.Serializer) Option {
	return func(e *Encipherer) {
		e.serializer = serializer
	}
}
