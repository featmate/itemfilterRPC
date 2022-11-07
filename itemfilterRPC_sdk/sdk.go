package itemfilterRPC_test

import (
	"github.com/Golang-Tools/grpcsdk"
	"github.com/featmate/itemfilterRPC/itemfilterRPC_pb"
)

type Conf grpcsdk.SDKConfig

var SDK = grpcsdk.New(itemfilterRPC_pb.NewITEMFILTERRPCClient, &itemfilterRPC_pb.ITEMFILTERRPC_ServiceDesc)
