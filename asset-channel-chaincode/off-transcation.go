package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
1. `InitLedger`：初始化空账本
2. `CreateTransaction`：创建新的交易，包括交易 ID、货币类型、金额、兑换目标、兑换比例和是否已执行等信息
3. `ExecuteTransaction`：执行交易，检查是否有与该交易匹配的待执行交易，如果有则执行两笔交易
4. `TransactionExists`：检查给定 ID 的交易是否存在
5. `ReadTransaction`：读取给定 ID 的交易
6. `UpdateTransaction`：更新给定交易的状态
7. `FindMatchingTransaction`：查找与给定交易匹配的待执行交易
*/

type Transaction struct {
    ID         string `json:"id"`
    Currency   string `json:"currency"`
    Amount     int    `json:"amount"`
    ExchangeTo string `json:"exchangeTo"`
    Rate       int    `json:"rate"`
    IsExecuted bool   `json:"isExecuted"`
}

type Exchange struct {
    contractapi.Contract
}

func (e *Exchange) InitLedger(ctx contractapi.TransactionContextInterface) error {
    return nil
}

func (e *Exchange) CreateTransaction(ctx contractapi.TransactionContextInterface, id string, currency string, amount int, exchangeTo string, rate int) error {
    exists, err := e.TransactionExists(ctx, id)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the transaction %s already exists", id)
    }

    tx := Transaction{
        ID:         id,
        Currency:   currency,
        Amount:     amount,
        ExchangeTo: exchangeTo,
        Rate:       rate,
        IsExecuted: false,
    }

    txJSON, err := json.Marshal(tx)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(tx.ID, txJSON)
    if err != nil {
        return err
    }

    return nil
}

func (e *Exchange) ExecuteTransaction(ctx contractapi.TransactionContextInterface, id string) error {
    tx, err := e.ReadTransaction(ctx, id)
    if err != nil {
        return err
    }

    if tx.IsExecuted {
        return fmt.Errorf("the transaction %s has already been executed", id)
    }

    matchingTx, err := e.FindMatchingTransaction(ctx, id, tx.Currency, tx.ExchangeTo, tx.Rate)
    if err != nil {
        return err
    }

    if matchingTx == nil {
        return fmt.Errorf("no matching transaction found for transaction with id %s", id)
    }

    tx.IsExecuted = true
    matchingTx.IsExecuted = true

    txJSON, err := json.Marshal(tx)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(tx.ID, txJSON)
    if err != nil {
        return err
    }

    matchingTxJSON, err := json.Marshal(matchingTx)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(matchingTx.ID, matchingTxJSON)
    if err != nil {
        return err
    }

    return nil
}

func (e *Exchange) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    txJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, err
    }

    return txJSON != nil, nil
}

func (e *Exchange) ReadTransaction(ctx contractapi.TransactionContextInterface, id string) (*Transaction, error) {
    transactionBytes, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if transactionBytes == nil {
        return nil, fmt.Errorf("the transaction %s does not exist", id)
    }

    var transaction Transaction
    err = json.Unmarshal(transactionBytes, &transaction)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
    }

    return &transaction, nil
}

func (e *Exchange) UpdateTransaction(ctx contractapi.TransactionContextInterface, id string, isExecuted bool) error {
    transaction, err := e.ReadTransaction(ctx, id)
    if err != nil {
        return err
    }

    transaction.IsExecuted = isExecuted

    transactionBytes, err := json.Marshal(transaction)
    if err != nil {
        return fmt.Errorf("failed to marshal transaction JSON: %v", err)
    }

    err = ctx.GetStub().PutState(id, transactionBytes)
    if err != nil {
        return fmt.Errorf("failed to update transaction: %v", err)
    }

    return nil
}


func (e *Exchange) FindMatchingTransaction(ctx contractapi.TransactionContextInterface, id string, currency string, exchangeTo string, rate int) (*Transaction, error) {
    query := fmt.Sprintf(`{"selector":{"$and":[{"id":{"$ne":"%s"}},{"currency":"%s"},{"exchangeTo":"%s"},{"rate":%d},{"isExecuted":false}]}}`, id, currency, exchangeTo, rate)

    resultsIterator, err := ctx.GetStub().GetQueryResult(query)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var tx Transaction
        err = json.Unmarshal(queryResponse.Value, &tx)
        if err != nil {
            return nil, err
        }

        return &tx, nil
    }

    return nil, nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(&Exchange{})
    if err != nil {
        log.Panicf("Error creating exchange chaincode: %v", err)
    }

    if err := chaincode.Start(); err != nil {
        log.Panicf("Error starting exchange chaincode: %v", err)
    }
}

   
