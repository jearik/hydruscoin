/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package coin

import (
	"github.com/golang/protobuf/proto"
)

func (coin *Hydruscoin) queryAddr(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	addr := args[0]
	queryResult := new(QueryAddrResult)

	account, err := store.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	queryResult.Account = account

	logger.Debugf("query addr[%s] result: %+v", addr, queryResult)
	return proto.Marshal(queryResult)
}

func (coin *Hydruscoin) queryAddrs(store Store, args []string) ([]byte, error) {
	results := &QueryAddrResults{
		Results: make([]*QueryAddrResult, 0),
	}

	for _, arg := range args {
		addr := arg
		queryResult := new(QueryAddrResult)

		account, err := store.GetAccount(addr)
		if err != nil {
			return nil, err
		}
		queryResult.Account = account

		results.Results = append(results.Results, queryResult)
		logger.Debugf("query addr[%s] result: %+v", addr, queryResult)
	}

	return proto.Marshal(results)
}

func (coin *Hydruscoin) queryTx(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	tx, _, err := store.GetTx(args[0])
	if err != nil {
		logger.Errorf("get tx info error: %v", err)
		return nil, err
	}
	logger.Debugf("query tx: %+v", tx)

	return proto.Marshal(tx)
}

func (coin *Hydruscoin) queryCoin(store Store, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, ErrInvalidArgs
	}

	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	logger.Debugf("query lepuscoin info: %+v", coinInfo)
	return proto.Marshal(coinInfo)
}
