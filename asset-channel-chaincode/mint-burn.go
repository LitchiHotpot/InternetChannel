package main

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 定义代币结构体
type Token struct {
    Owner  string `json:"owner"`
    Amount int    `json:"amount"`
}

// 定义链码结构体
type TokenChaincode struct {
    contractapi.Contract
}

// 铸造代币的函数
func (t *TokenChaincode) Mint(ctx contractapi.TransactionContextInterface, owner string, amount int) error {
    // 创建代币对象
    token := Token{
        Owner:  owner,
        Amount: amount,
    }

    // 序列化代币对象
    tokenAsBytes, err := json.Marshal(token)
    if err != nil {
        return fmt.Errorf("Failed to marshal token. %s", err.Error())
    }

    // 将代币存储到账户对应的状态数据库中
    tokenKey := fmt.Sprintf("token:%s", owner)
    err = ctx.GetStub().PutState(tokenKey, tokenAsBytes)
    if err != nil {
        return fmt.Errorf("Failed to put token. %s", err.Error())
    }

    return nil
}

// 销毁代币的函数
func (t *TokenChaincode) Burn(ctx contractapi.TransactionContextInterface, owner string, amount int) error {
    // 获取账户对应的代币对象
    tokenKey := fmt.Sprintf("token:%s", owner)
    tokenAsBytes, err := ctx.GetStub().GetState(tokenKey)
    if err != nil {
        return fmt.Errorf("Failed to read token. %s", err.Error())
    }

    // 如果账户已经有代币，则将其减少amount
    token := Token{}
    err = json.Unmarshal(tokenAsBytes, &token)
    if err != nil {
        return fmt.Errorf("Failed to unmarshal token. %s", err.Error())
    }
    if token.Amount < amount {
        return fmt.Errorf("Insufficient balance. Owner %s has only %v tokens", owner, token.Amount)
    }
    token.Amount -= amount
    updatedTokenAsBytes, err := json.Marshal(token)
    if err != nil {
        return fmt.Errorf("Failed to marshal token. %s", err.Error())
    }
    err = ctx.GetStub().PutState(tokenKey, updatedTokenAsBytes)
    if err != nil {
        return fmt.Errorf("Failed to put token. %s", err.Error())
    }

    // 调用铸币函数，向另一个账户铸造相应的代币
    err = t.Mint(ctx, "另一个账户的ID", amount)
    if err != nil {
        return fmt.Errorf("Failed to mint token. %s", err.Error())
    }

    return nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(&TokenChaincode{})
    if err != nil {
        fmt.Printf("Error creating token chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting token chaincode: %s", err.Error())
    }
}
