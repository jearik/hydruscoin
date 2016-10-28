# 简易UTXO与账户模型的数字资产智能合约Hydruscoin
[Hyperledger Fabric](https://github.com/hyperledger/fabric)是Linux基金会开源的区块链平台实现,由IBM主导开发。对广大智能合约开发者来说,我们可以摒弃Fabric底层实现不谈,只需要关注上层智能合约开发包chaincode就可以了。但是由于区块链技术概念较新,更不用说基于区块链技术进行智能合约开发了,能借鉴的文章、demo都很少。我结合自身实践和官方chaincode例子,实现了一个类似于bitcoin的数字货币。[现已开源](https://github.com/hydrusio/hydruscoin), 下面为大家简要讲解一下开发思路:

## chinacode初识
什么样的一个应用程序能被称作是一个chaincode?
```
type SimpleChaincode struct {
}
// Init callback representing the invocation of a chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return nil, nil
}
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return nil, nil
}
// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return nil, nil
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		panic(err)
	}
}
```
只需要`cc`对象实现`Chaincode`接口,并在main函数中通过`shim.Start(cc)`启动,我们就可以说这是一个chaincode(智能合约)应用程序

*注:只要返回错误,接口所进行写入操作都会自动回滚*

## 主要功能
* coinbase交易
* 数字货币转账
* 账户查询
* 交易查询
* 数字货币简单统计

## 实现过程
### Init接口
`Init`主要是用于chaincode注册,在chaincode整个生命周期中,只执行一次。主要目的是初始化一些变量。在我的例子中,hydruscoin只是初始化了一个统计对象:
```
// construct a new store
store := MakeChaincodeStore(stub)

// deploy hydruscoin chaincode only need to set coin stater
if err := store.InitCoinInfo(); err != nil {
	return nil, err
}
```
```
func (s *ChaincodeStore) InitCoinInfo() error {
	coinInfo := &HydruscoinInfo{
		CoinTotal:    0,
		AccountTotal: 0,
		TxoutTotal:   0,
		TxTotal:      0,
		Placeholder:  "placeholder",
	}

	return s.PutCoinInfo(coinInfo)
}

func (s *ChaincodeStore) PutCoinInfo(coinfo *HydruscoinInfo) error {
	coinBytes, err := proto.Marshal(coinfo)
	if err != nil {
		return err
	}

	if err := s.stub.PutState(coinInfoKey, coinBytes); err != nil {
		return err
	}

	return nil
}
```

**注意**,通过REST API调用chaincode `Init`接口时, API是同步返回执行结果。

### Invoke接口
`Invoke`接口是`Chaincode`设计的重中之重,是智能合约与外部进行数据交换的主要通道。通过给`Invoke`接口传递不同的`function`,配以不同的`args`,整个智能合约就活灵活现起来。

**注意**,通过REST API调用chaincode `Invoke`接口时, API是同步返回调用结果,但是至于是否执行成功,需要通过`Query`接口查询。

在我的例子里,`Invoke`主要有两个功能: *coinbase*(产生数字货币)和*transfer*(转移数字货币)。
```
// Invoke function
const (
	IF_COINBASE string = "invoke_coinbase"
	IF_TRANSFER string = "invoke_transfer"
)

// Invoke
func (coin *Hydruscoin) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case IF_COINBASE:
		return coin.coinbase(store, args)
	case IF_TRANSFER:
		return coin.transfer(store, args)
	default:
		return nil, ErrUnsupportedOperation
	}
}
```

**2016/10/28新增**
新增了一个invoke方法,`register`用于注册钱包地址,目的是当用户钱包未有交易产生时也能够查询账户余额


#### Invoke Coinbase
`coinbase`一词源自比特币,意指凭空产生数字货币,即货币发行,央行印钞一个道理。代码看着有些复杂,我简单梳理一下逻辑关系。

1. 解码经过base64编码的交易`TX`
2. 反序列化proto message至交易对象
3. 循环验证交易`TX`的交易输入`TX_IN`,如果不是coinbase交易,则返回错误
4. 循环验证交易`TX`的交易输出`TX_OUT`
    1. 读取交易输出`TX_OUT`相关地址的账号信息,如果错误,默认为账户不存在,直接新建一个
    2. 将交易输出`TX_OUT`的额度转移到账户余额里,并保存该交易输出
    3. 将更新写入区块链中
5. 循环交易输出的同时也将同步更新统计对象
6. 单独存储该交易至链上

#### Invoke Transfer
`transfer`即所谓的转账。业务简单流程如下:

1. 解码经过base64编码的交易`TX`
2. 反序列化proto message至交易对象
3. 验证交易发起者`founder`的身份信息
4. 循环验证交易`TX`的交易输入`TX_IN`
    1. 验证交易发起者`founder`是否拥有此交易输出(*ps:交易输入即之前交易的输出*)
    2. 验证此交易输出`TX_OUT`是否是可花费的
    3. 如果都验证成功,删除之前的交易输出`TX_OUT`,更新交易发起者`founder`的账户信息
5. 循环验证交易`TX`的交易输出`TX_OUT`
    1. 查询交易输出`TX_OUT`的接收者账户信息
    2. 根据交易输出`TX_OUT`,更新接收者账户
6. 循环验证的同时同步更新统计对象
7. 全局验证交易`TX`的正确性
8. 单独存储该交易至链上

#### Invoke Register
`register`方法很简单,接收一个钱包地址,初始化一个账户对象,存入区块链。

### Query接口
`Query`是`Chaincode`的查询接口,通过该接口可以查询存储在区块链上的数据的最终状态。*注意:在`Query`中任何对数据的更改都是不被允许的*

**注意**,通过REST API调用chaincode `Query`接口时, API是同步返回执行结果。

我现在实现了3个简单的查询功能:查询账户集信息,查询交易信息,查询统计信息。
```
// Query function
const (
	QF_ADDRS = "query_addrs"
	QF_TX    = "query_tx"
	QF_COIN  = "query_coin"
)

// Query
func (coin *Hydruscoin) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case QF_ADDRS:
		return coin.queryAddrs(store, args)
	case QF_TX:
		return coin.queryTx(store, args)
	case QF_COIN:
		return coin.queryCoin(store, args)
	default:
		return nil, ErrUnsupportedOperation
	}
}
```

#### Query Addr/Addrs
在现在的设计中,任何人都可以查询任意一个账户的信息,只要你知道账户的地址。这样的设计肯定是有问题的,但是我们不就是个demo么?
```
func (coin *Hydruscoin) queryAddrs(store Store, args []string) ([]byte, error) {
	results := &QueryAddrResults{
		Accounts: make(map[string]*Account),
	}

	for _, addr := range args {
		account, err := store.GetAccount(addr)
		if err != nil {
			logger.Warningf("store.GetAccount() return error: %v", err)
			continue
		}

		results.Accounts[addr] = account
		logger.Debugf("query addr[%s] account: %+v", addr, account)
	}

	return proto.Marshal(results)
}
```
可以看到,仅仅是根据输入的地址信息,去区块链上查找账户信息,未作任何处理。

#### Query TX
在之前的`Coinbase`和`Transfer`中,最后都将交易信息存储在了链上。只要你知道交易`TX`的Hash值,那么你就能通过该接口查询到交易的具体信息
```
// GetTx returns a transaction for the given hash
func (s *ChaincodeStore) GetTx(key string) (*TX, bool, error) {
	data, err := s.stub.GetState(key)
	if err != nil {
		return nil, false, fmt.Errorf("Error getting state from stub:  %s", err)
	}
	if data == nil || len(data) == 0 {
		return nil, false, nil
	}

	tx, err := ParseTXBytes(data)
	if err != nil {
		return nil, false, err
	}

	return tx, true, nil
}
```
#### Query Coin
这个不用多说,就是获取统计对象
```
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
```

## 暂存问题
1. 未对交易进行签名,在demo阶段,暂时还未实现


## Change Log
### 2016/10/28
* 删除查询API`query_addr`,统一使用`query_addrs`
* 增加账号注册API,现阶段只接受唯一参数,即钱包地址