package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// 业务实体
type AssertsExchangeCC struct{}

// 在此处定义json tag 主要是我们在存储状态的时候
// 状态的值是以字符数组的形式存储的，那么需要将业务实体序列化之后再存储
// 序列化最好的方式是使用json
// 当然也可以选择 messagepack || protobuf
// 而这里为了方便使用json
type User struct {
	Name string `json:"name"`
	Id string  `json:"id"`
	Assets map[string]string `json:"asserts"`// 资产Id --> 资产Name
}

// 资产
type Asset struct {
	Name string `json:"name"`
	Id string `json:"id"`
	Metadata map[string]string `json:"metadata"` // 特殊属性
}

// 资产变更对象
type AssetHistory struct {
	AssetId string `json:"asset_id"`
	OriginOwnerId string `json:"origin_owner_id"` // 资产的原始拥有者
    CurrentOwnerId string `json:"current_owner_id"` // 资产当前拥有者
}

// 用户开户
func userRegister(args []string) pb.Response {
	return shim.Success(nil)
}

// 用户销户
func userDestroy(args []string) pb.Response {
	return shim.Success(nil)
}

// 资产登记
func assetEnroll(args []string) pb.Response {
	return shim.Success(nil)
}

// 资产转让
func assetExchange(args []string) pb.Response {
	return shim.Success(nil)
}

func (c *AssertsExchangeCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()

	switch funcName {
	case "userRegister":
		return userRegister(args)
	case "userDestroy":
		return userDestroy(args)
	case "assetEnroll":
		return assetEnroll(args)
	case "assetExchange":
		return assetExchange(args)
	}

	return shim.Success(nil)
}

func main() {
	
}
