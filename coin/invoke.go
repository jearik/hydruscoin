/*
Copyright Hydrusio Inc. 2016 All Rights Reserved.
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

func (coin *Hydruscoin) registerAccount(store Store, args []string) ([]byte, error) {
	addr := args[0]

	if tmpaccount, err := store.GetAccount(addr); err == nil && tmpaccount != nil && tmpaccount.Addr == addr {
		logger.Warningf("Hydruscoin account(%s) already registered.", addr)
		return nil, ErrAlreadyRegisterd
	}

	account := &Account{
		Addr:    addr,
		Balance: 0,
		Txouts:  make(map[string]*TX_TXOUT),
	}
	if err := store.PutAccount(account); err != nil {
		logger.Errorf("store.PutAccount(%#v) return error: %v", account, err)
		return nil, err
	}

	return nil, nil
}
