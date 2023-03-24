package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/*

在这个链码中，我们定义了一个ExchangeRateChaincode类型，它有三个方法：

1. `createExchangeRate`：用于创建一个新的兑换比例信息。
2. `getExchangeRate`：用于获取指定兑换比例信息。
3. `listExchangeRates`：用于列出所有可用的兑换比例信息。

我们还定义了一个ExchangeRate类型，用于存储兑换比例信息。每个ExchangeRate对象都包括以下字段：

- `FromCurrency`：要兑换的原始货币。
- `ToCurrency`：要兑换的目标货币。
- `ExchangeRate`：兑换比例。

在`createExchangeRate`方法中，我们首先验证传入的参数数量是否正确。然后，我们根据传入的参数创建一个ExchangeRate对象，并将其序列化为JSON格式。接下来，我们将该对象的JSON表示作为值，使用`PutState`函数将其存储到区块链上。

在`getExchangeRate`方法中，我们首先验证传入的参数数量是否正确。然后，我们使用传入的参数构建一个交易对的键，从区块链中获取该键对应的值。最后，我们将该值作为响应返回给调用方。

在`listExchangeRates`方法中，我们使用`GetStateByPartialCompositeKey`函数获取所有以"exchangeRate"为前缀的键值对。然后，我们迭代这些键值对，并将每个ExchangeRate对象反序列化到一个ExchangeRate数组中。最后，我们将该数组序列化为JSON格式，并将其作为响应返回给调用方。

为了部署和测试这个链码，你需要按照Hyperledger Fabric的指南来构建和安装链码。然后，你可以使用fabric-samples仓库中的示例脚本和配置来启动一个本地的Fabric网络，并使用Fabric命令行工具调用链码。

*/

type ExchangeRateChaincode struct {
}

type ExchangeRate struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	ExchangeRate float64 `json:"exchange_rate"`
}

func (t *ExchangeRateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *ExchangeRateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "createExchangeRate" {
		return t.createExchangeRate(stub, args)
	} else if function == "getExchangeRate" {
		return t.getExchangeRate(stub, args)
	} else if function == "listExchangeRates" {
		return t.listExchangeRates(stub, args)
	}

	return shim.Error(fmt.Sprintf("Invalid function: %s", function))
}

func (t *ExchangeRateChaincode) createExchangeRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fromCurrency := args[0]
	toCurrency := args[1]
	exchangeRate := args[2]

	exchangeRateObj := ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency: toCurrency,
		ExchangeRate: exchangeRate,
	}

	exchangeRateKey := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
	exchangeRateBytes, err := json.Marshal(exchangeRateObj)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(exchangeRateKey, exchangeRateBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ExchangeRateChaincode) getExchangeRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	fromCurrency := args[0]
	toCurrency := args[1]
	exchangeRateKey := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)

	exchangeRateBytes, err := stub.GetState(exchangeRateKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	if exchangeRateBytes == nil {
		return shim.Error("Exchange rate not found")
	}

	return shim.Success(exchangeRateBytes)
}

func (t *ExchangeRateChaincode) listExchangeRates(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	exchangeRateIterator, err := stub.GetStateByPartialCompositeKey("exchangeRate", []string{})

	if err != nil {
		return shim.Error(err.Error())
	}

	defer exchangeRateIterator.Close()

	var exchangeRates []ExchangeRate

	for exchangeRateIterator.HasNext() {
		exchangeRateKV, err := exchangeRateIterator.Next()

		if err != nil {
		
			return shim.Error(err.Error())
	}

	var exchangeRate ExchangeRate

	err = json.Unmarshal(exchangeRateKV.Value, &exchangeRate)

	if err != nil {
		return shim.Error(err.Error())
	}

	exchangeRates = append(exchangeRates, exchangeRate)
}

	exchangeRatesBytes, err := json.Marshal(exchangeRates)

	if err != nil {
		return shim.Error(err.Error())
	}	

	return shim.Success(exchangeRatesBytes)
}

func main() {
	err := shim.Start(new(ExchangeRateChaincode))
	if err != nil {
		fmt.Printf("Error starting ExchangeRateChaincode: %s", err)
	}
」

