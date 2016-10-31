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

import "errors"

var (
	// ErrInvalidArgs is returned if there are some unused args or not enough args in params
	ErrInvalidArgs = errors.New("invalid args")

	// ErrInvalidFunction is returnd if chaincode interface get unsupported function name
	ErrInvalidFunction = errors.New("invalid function")

	// ErrInvalidTxKey returned if given key is invalid
	ErrInvalidTxKey = errors.New("invalid tx key")

	// ErrInvalidTX
	ErrInvalidTX = errors.New("transaction invalid")

	// ErrUnsupportedOperation returned if invoke or query using unsupported function name
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrMustCoinbase
	ErrMustCoinbase = errors.New("tx must be coinbase")

	// ErrCantCoinbase
	ErrCantCoinbase = errors.New("tx must not be coinbase")

	// ErrTxInOutNotBalance returned when txouts + fee != txins
	ErrTxInOutNotBalance = errors.New("tx in & out not balance")

	// ErrTxOutMoreThanTxIn
	ErrTxOutMoreThanTxIn = errors.New("tx out more than tx in")

	// ErrKeyNoData
	ErrKeyNoData = errors.New("state key found, but no data")

	// ErrCollisionTxOut
	ErrCollisionTxOut = errors.New("account has collision tx out")

	// ErrTxNoFounder
	ErrTxNoFounder = errors.New("tx has no founder")

	// ErrAccountNoTxOut
	ErrAccountNoTxOut = errors.New("account has no such tx out")

	// ErrAccountNotEnoughBalance
	ErrAccountNotEnoughBalance = errors.New("account has not enough balance")

	// ErrTxOutLock
	ErrTxOutLock = errors.New("tx out can be spend only after until time")

	// ErrAlreadyRegisterd
	ErrAlreadyRegisterd = errors.New("the addr has been registerd into coin")
)
