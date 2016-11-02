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
	"time"

	"github.com/golang/protobuf/proto"
)

func (coin *Hydruscoin) transfer(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	// parse tx
	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		logger.Errorf("Error decode tx bytes: %v", err)
		return nil, err
	}

	tx, err := ParseTXBytes(txData)
	if err != nil {
		return nil, err
	}

	// coin stat
	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	execResult := &ExecResult{}
	txHash := TxHash(tx)
	if tx.Founder == "" {
		return nil, ErrTxNoFounder
	}

	//founderAccount, err := store.GetAccount(tx.Founder)
	//if err != nil {
	//	return nil, ErrTxNoFounder
	//}

	for _, ti := range tx.Txin {
		prevTxHash := ti.SourceHash
		prevOutputIx := ti.Ix
		ownerAddr := ti.Addr
		keyToPrevOutput := &Key{TxHashAsHex: prevTxHash, TxIndex: prevOutputIx}

		ownerAccount, err := store.GetAccount(ownerAddr)
		if err != nil {
			return nil, err
		}
		txout, ok := ownerAccount.Txouts[keyToPrevOutput.String()]
		if !ok {
			return nil, ErrAccountNoTxOut
		}

		// can spend?
		if txout.Until > 0 {
			untilTime := time.Unix(txout.Until, 0).UTC()
			if untilTime.After(time.Now().UTC()) {
				return nil, ErrTxOutLock
			}
		}

		if ownerAccount.Balance < txout.Value {
			return nil, ErrAccountNotEnoughBalance
		}
		ownerAccount.Balance -= txout.Value
		delete(ownerAccount.Txouts, keyToPrevOutput.String())
		// save owner account
		if err := store.PutAccount(ownerAccount); err != nil {
			return nil, err
		}

		// coin stat
		coinInfo.TxoutTotal -= 1
		execResult.SumPriorOutputs += txout.Value
	}
	// save founder account
	//if err := store.PutAccount(founderAccount); err != nil {
	//	return nil, err
	//}

	for idx, to := range tx.Txout {
		account, err := store.GetAccount(to.Addr)
		if err != nil {
			logger.Warningf("get account[%s] doesnt exist, creating one...", to.Addr)

			account = new(Account)
			account.Txouts = make(map[string]*TX_TXOUT)
			account.Addr = to.Addr

			coinInfo.AccountTotal += 1
		}
		if account.Txouts == nil || len(account.Txouts) == 0 {
			account.Txouts = make(map[string]*TX_TXOUT)
		}

		outKey := &Key{TxHashAsHex: txHash, TxIndex: uint32(idx)}
		if _, ok := account.Txouts[outKey.String()]; ok {
			return nil, ErrCollisionTxOut
		}

		account.Balance += to.Value
		account.Txouts[outKey.String()] = to
		if err := store.PutAccount(account); err != nil {
			return nil, err
		}

		// coin stat
		coinInfo.TxoutTotal += 1
		execResult.SumCurrentOutputs += to.Value
	}

	// current outputs must less than prior outputs
	if execResult.SumCurrentOutputs > execResult.SumPriorOutputs {
		return nil, ErrTxOutMoreThanTxIn
	}

	// one of transfer main point is in == out, no coin mined, no coin lose
	if execResult.SumCurrentOutputs != execResult.SumPriorOutputs {
		return nil, ErrTxInOutNotBalance
	}

	if err := store.PutTx(tx); err != nil {
		logger.Errorf("put tx error: %v", err)
		return nil, err
	}
	logger.Debug("put tx into world state")

	// save coin stat
	coinInfo.TxTotal += 1
	if err := store.PutCoinInfo(coinInfo); err != nil {
		return nil, err
	}

	return proto.Marshal(execResult)
}
