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
	"github.com/golang/protobuf/proto"
	"crypto/sha256"
	"encoding/hex"
)

// ParseTXBytes unmarshal txData into TX object
func ParseTXBytes(txData []byte) (*TX, error) {
	tx := new(TX)
	err := proto.Unmarshal(txData, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// TxHash generates the Hash for the transaction.
func TxHash(tx *TX) string {
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		return ""
	}

	fHash := sha256.Sum256(txBytes)
	lHash := sha256.Sum256(fHash[:])
	return hex.EncodeToString(lHash[:])
}

// ParseHydruscoinInfoBytes unmarshal infoBytes into HydruscoinInfo
func ParseHydruscoinInfoBytes(infoBytes []byte) (*HydruscoinInfo, error) {
	info := new(HydruscoinInfo)
	if err := proto.Unmarshal(infoBytes, info); err != nil {
		return nil, err
	}

	return info, nil
}
