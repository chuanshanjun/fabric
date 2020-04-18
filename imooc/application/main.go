package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"net/http"
	"time"
)

// 此处使用的是gin框架
func main() {
	router := gin.Default()

	// 定义路由
	{
		router.POST("/users", userRegister)
		router.GET("/users/:id", queryUser)
		router.DELETE("/users/:id", deleteUser)
		router.GET("/assets/get/:id", queryAsset)
		router.POST("/assets/enroll", assetsEnroll)
		router.POST("/assets/exchange", assetsExchange)
		router.GET("/assets/exchange/history", assetsExchangeHistory)
	}

	router.Run() // listen and server on 0.0.0.0:8080
}

type UserRegisterRequest struct {
	Id string `form:"id" binding:"required"`
	Name string `form:"name" binding:"required"`
}

// 用户开户
func userRegister(ctx *gin.Context) {
	// 参数处理
	req := new(UserRegisterRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	// 区块链交互
	resp, err := channelExecute("userRegister", [][]byte{
		[]byte(req.Name),
		[]byte(req.Id),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 查询用户
func queryUser(ctx *gin.Context) {
	// Param是指参数从path中获取
	userId := ctx.Param("id")

	resp, err := channelQuery("queryUser", [][]byte{
		[]byte(userId),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 用户销户
func deleteUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	resp, err := channelExecute("userDestroy", [][]byte{
		[]byte(userId),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 查询资产
func queryAsset(ctx *gin.Context) {
	assetId := ctx.Param("id")

	resp, err := channelQuery("queryAsset", [][]byte{
		[]byte(assetId),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

type AssetsEnrollRequest struct {
	// 参数从form表单中获取
	AssetId string `form:"assetsid" binding:"required"`
	AssetName string `form:"assetname" binding:"required"`
	Metadata string `form:"metadata" binding:"required"`
	OwnerId string `form:"ownerid" binding:"required"`
}

// 资产登记
func assetsEnroll(ctx *gin.Context) {
	req := new(AssetsEnrollRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	resp, err := channelExecute("assetEnroll", [][]byte{
		[]byte(req.AssetName),
		[]byte(req.AssetId),
		[]byte(req.Metadata),
		[]byte(req.OwnerId),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

type AssetsExchangeRequest struct {
	AssetId string `form:"assetsid" binding:"required"`
	OriginOwnerId string `form:"originownerid" binding:"required"`
	CurrentOwnerId string `form:"currentownerid" binding:"required"`
}

// 资产转让
func assetsExchange(ctx *gin.Context) {
	req := new(AssetsExchangeRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	resp, err := channelExecute("assetExchange", [][]byte{
		[]byte(req.OriginOwnerId),
		[]byte(req.AssetId),
		[]byte(req.CurrentOwnerId),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 资产历史变更记录
func assetsExchangeHistory(ctx *gin.Context) {
	assetId := ctx.Query("assetid")
	queryType := ctx.Query("querytype")

	resp, err := channelQuery("queryAssetHistory", [][]byte{
		[]byte(assetId),
		[]byte(queryType),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

var (
	sdk *fabsdk.FabricSDK
	channelName = "assetschannel"
	chaincodeName = "assets"
	org = "org1"
	user = "Admin"
	//configPath = "$GOPATH/src/github.com/hyperledger/fabric/imooc/application/config.yaml"
	configPath = "./config.yaml"
)

// sdk通过init()方法去初始化
func init() {
	var err error
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}

// 区块链管理
func manageBlockchain() {
	// 表明自己的身份 17:50
	ctx := sdk.Context(fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := resmgmt.New(ctx)
	if err != nil {
		panic(err)
	}

	// 具体操作
	// 创建channel和order有关，所以还需要加上orderer的信息
	cli.SaveChannel(resmgmt.SaveChannelRequest{}, resmgmt.WithOrdererURL("orderer.imocc.com"))
}

// 区块链数据查询 账本的查询
func queryBlockchain() {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	// 创建账本的相关客户端
	cli, err := ledger.New(ctx)
	if err != nil {
		panic(err)
	}

	// 查询信息
	// 往org1.peer0节点去查询信息，因为区块链每个节点的数据都是一样的，所以往哪个节点查都是一样的
	// 但前提是节点的数据是最新的数据
	resp, err := cli.QueryInfo(ledger.WithTargetURLs("peer0.org1.imocc.com"))
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	// 1 查询当前区块,然后再查询前一个区块，这样一个个往前推
	cli.QueryBlockByHash(resp.BCI.CurrentBlockHash)

	// 2 从0开始查询区块高度，这样可以查出所有的区块
	for i := uint64(0); i <= resp.BCI.Height; i++ {
		cli.QueryBlock(i)
	}

	// 两种方式都可以实现，区块链浏览器读取区块的功能
}

// 区块链交互
func channelExecute(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	// 实例化一个channel的cli
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 状态更新, insert/update/delete
	// 如果你要更新的键值对已经存在了，那它就是update，如果不存在则是insert
	// 如果更新的键值对的值是空的化可以认为他是delete
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetURLs("peer0.org1.imocc.com")) // 往哪个节点去执行
	if err != nil {
		return channel.Response{}, err
	}

	// 链码事件监听,对监听的事件可以做一些业务上的逻辑处理
	go func() {
		// channel
		reg, ccevt, err := cli.RegisterChaincodeEvent(chaincodeName, "eventname")
		if err != nil {
			return
		}
		// 注意注册一个事件及注销一个事件，总是成对出现的
		defer cli.UnregisterChaincodeEvent(reg)

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-ccevt:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
				case <-timeoutctx.Done():
					fmt.Println("event timeout, exit!")
					return
			}
		}

		// 上面通过channel来进行监听事件的处理，下面用event来处理
		// event
		//eventcli, err := event.New(ctx)
		//if err != nil {
		//	return
		//}

		//eventcli.RegisterChaincodeEvent(chaincodeName, "eventname")
	}()

	// 交易状态事件监听
	//go func() {
	//	eventcli, err := event.New(ctx)
	//	if err != nil {
	//		return
	//	}
	//
	//	reg, status, err := eventcli.RegisterTxStatusEvent(string(resp.TransactionID))
	//	defer eventcli.Unregister(reg) // 注册必有注销
	//
	//	timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	//	defer cancel()
	//
	//	for {
	//		select {
	//		case evt := <-status:
	//			fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
	//		case <- timeoutctx.Done():
	//			fmt.Println("event timeout, exist!")
	//			return
	//		}
	//	}
	//}()

	return resp, nil
}

func channelQuery(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	// 实例化一个channel的cli
	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 只是对应状态的查询,select
	return cli.Query(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetURLs("peer0.org1.imocc.com")) // 往哪个节点去执行
}

// 事件监听
//func eventHandle() {
//	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))
//
//	cli, err := event.New(ctx)
//	if err != nil {
//		panic(err)
//	}
//
//	// 交易状态事件
//	// 链码事件 业务事件
//	// 区块事件，当有新的区块产生的时候我就可以监听到
//	reg, blkevent, err := cli.RegisterBlockEvent()
//	if err != nil {
//		panic(err)
//	}
//	defer cli.Unregister(reg)
//
//	timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
//	defer cancel()
//
//	for {
//		select {
//		case evt := <-blkevent:
//		fmt.Printf("received a block", evt)
//		case <-timeoutctx.Done():
//			fmt.Println("event timeout, exist")
//			return
//		}
//	}
//}
