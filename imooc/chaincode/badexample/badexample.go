package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math/rand"
	"strconv"
	"time"
)

// 类，用于实现chaincode
type BadExampleCC struct {

}

// Init()及Invoke是每个chaincode必须实现的方法
func (c *BadExampleCC) Init(stubInterface shim.ChaincodeStubInterface) pb.Response {
	// 初始化的时候，直接返回初始化成功，并没有做任何初始化的操作
	return shim.Success(nil)
}

func (c *BadExampleCC) Invoke( shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(bytes.NewBufferString(strconv.Itoa(int(rand.Int63n(time.Now().Unix())))).Bytes())
}

func main() {
	// 启动链码
	err := shim.Start(new(BadExampleCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
