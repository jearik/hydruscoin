/*
Copyright Hydrusio Labs Inc. 2016 All Rights Reserved.
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
	"encoding/base64"
	"math"

	"github.com/golang/protobuf/proto"
)

func (coin *Hydruscoin) coinbase(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		logger.Errorf("Decoding base64 error: %v\n", err)
		return nil, err
	}

	tx, err := ParseTXBytes(txData)
	if err != nil {
		logger.Errorf("Unmarshal tx bytes error: %v\n", err)
		return nil, err
	}
	if !tx.Coinbase {
		return nil, ErrMustCoinbase
	}

	txhash := TxHash(tx)
	execResult := &ExecResult{}
	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	// Loop through outputs first
	for index, output := range tx.Txout {
		if output.Addr == "" {
			return nil, ErrInvalidTX
		}

		// change coin info
		coinInfo.CoinTotal += output.Value
		coinInfo.TxoutTotal += 1

		outerAccount, err := store.GetAccount(output.Addr)
		if err != nil {
			logger.Warningf("account[%s] is not existed, creating one...", output.Addr)

			outerAccount = new(Account)
			outerAccount.Addr = output.Addr
			outerAccount.Txouts = make(map[string]*TX_TXOUT)

			coinInfo.AccountTotal += 1
		}
		if outerAccount.Txouts == nil || len(outerAccount.Txouts) == 0 {
			outerAccount.Txouts = make(map[string]*TX_TXOUT)
		}

		currKey := &Key{TxHashAsHex: txhash, TxIndex: uint32(index)}
		if _, ok := outerAccount.Txouts[currKey.String()]; ok {
			return nil, ErrCollisionTxOut
		}

		// store tx out into account
		outerAccount.Txouts[currKey.String()] = output
		outerAccount.Balance += output.Value

		if err := store.PutAccount(outerAccount); err != nil {
			logger.Errorf("Error update account: %v, account info: %+v", err, outerAccount)
			return nil, err
		}
		logger.Debugf("put tx output %s:%v", currKey.String(), output)
		execResult.SumCurrentOutputs += output.Value
	}

	// Now loop over inputs
	for _, input := range tx.Txin {
		if math.MaxUint32 != input.Ix {
			logger.Errorf("coinbase tx can not has other input")
			return nil, ErrMustCoinbase
		}
	}

	if err := store.PutTx(tx); err != nil {
		logger.Errorf("put tx error: %v", err)
		return nil, err
	}
	logger.Debug("put tx into world state")

	// tx total counter
	coinInfo.TxTotal += 1

	// save coin info counter
	if err := store.PutCoinInfo(coinInfo); err != nil {
		logger.Errorf("Error put coin info: %v", err)
		return nil, err
	}

	logger.Debugf("coinbase execute result: %+v", execResult)
	return proto.Marshal(execResult)
}
