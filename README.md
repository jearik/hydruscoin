# Hydruscoin
一个简单支持UTXO模型和Account模型的数字货币,基于Hyperledger Chaincode。主要用于普及智能合约,非商用。

*注:适用于Fabric V0.6*

## 目录
* coin: 货币代码,诸如coinbase,转账等功能
* client: 简单的一些tx生成函数

## 实现功能
* coinbase交易
* 数字货币转账
* 账户查询
* 交易查询
* 数字货币简单统计

## 代码调试
请参考[Hyperledger Chaincode开发调试教程](http://hydrus.io/2016/10/hyperledger-chaincode-debug/)

## 代码解读
请参考[简易UTXO与账户模型的数字资产智能合约Hydruscoin](http://hydrus.io/2016/10/demo-hydruscoin/)

## 代码贡献
欢迎大家提pr丰富hydruscoin的功能

## Change Log
### 2016/10/28
* 删除查询API`query_addr`,统一使用`query_addrs`
* 增加账号注册API,现阶段只接受唯一参数,即钱包地址

### 2016/11/02
* 交易输入`TX_IN`新增地址字段,交易发起者可以使用其他地址的UTXO
* *coinbase*交易默认交易输入为空

## License
Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
```
http://www.apache.org/licenses/LICENSE-2.0
```
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.