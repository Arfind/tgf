package rpc

import (
	"github.com/Arfind/tgf"
	"github.com/Arfind/tgf/log"
	"golang.org/x/net/context"
)

//***************************************************
//@Link  https://github.com/thkhxm/tgf
//@Link  https://gitee.com/timgame/tgf
//@QQ群 7400585
//author tim.huang<thkhxm@gmail.com>
//@Description
//2023/2/26
//***************************************************

// GateService
// @Description: 默认网关
type GateService struct {
	Module
	tcpBuilder ITCPBuilder
	tcpService ITCPService
}

func (g *GateService) GetName() string {
	return tgf.GatewayServiceModuleName
}

func (g *GateService) GetVersion() string {
	return "1.0"
}

func (g *GateService) Startup() (bool, error) {
	var ()
	g.tcpService = newDefaultTCPServer(g.tcpBuilder)
	g.tcpService.Run()
	return true, nil
}

func (g *GateService) UploadUserNodeInfo(ctx context.Context, args *UploadUserNodeInfoReq, reply *UploadUserNodeInfoRes) error {
	var ()
	if ok := g.tcpService.UpdateUserNodeInfo(args.UserId, args.ServicePath, args.NodeId); !ok {
		reply.ErrorCode = -1
	}
	log.DebugTag("gate", "修改用户节点信息 userId=%v servicePath=%v nodeId=%v res=%v", args.UserId, args.ServicePath, args.NodeId, reply)
	return nil
}

func (g *GateService) Login(ctx context.Context, args *LoginReq, reply *LoginRes) error {
	var ()
	//踢人,重复登录的
	if !g.tcpService.Offline(args.UserId, true) {
		BorderRPCMessage(ctx, Offline.New(&OfflineReq{UserId: args.UserId}, new(OfflineRes)))
	}
	err := g.tcpService.DoLogin(args.UserId, args.TemplateUserId)
	return err
}

func (g *GateService) Offline(ctx context.Context, args *OfflineReq, reply *OfflineRes) error {
	var ()
	if GetNodeId(ctx) == tgf.NodeId {
		return nil
	}
	g.tcpService.Offline(args.UserId, args.Replace)
	return nil
}

func (g *GateService) ToUser(ctx context.Context, args *ToUserReq, reply *ToUserRes) error {
	var ()
	for _, userId := range args.UserId {
		done := g.tcpService.ToUser(userId, args.MessageType, args.Data)
		if done {
			log.DebugTag("gate", "主动推送 userId=%v msgLen=%v", args.UserId, len(args.Data))
		}
	}

	return nil
}

func GatewayService(tcpBuilder ITCPBuilder) IService {
	service := &GateService{}
	service.tcpBuilder = tcpBuilder
	return service
}

var Gate = &Module{Name: "gate", Version: "1.0"}

var (
	UploadUserNodeInfo = &ServiceAPI[*UploadUserNodeInfoReq, *UploadUserNodeInfoRes]{
		ModuleName:  Gate.Name,
		Name:        "UploadUserNodeInfo",
		MessageType: Gate.Name + "." + "UploadUserNodeInfo",
	}

	ToUser = &ServiceAPI[*ToUserReq, *ToUserRes]{
		ModuleName:  Gate.Name,
		Name:        "ToUser",
		MessageType: Gate.Name + "." + "ToUser",
	}

	Login = &ServiceAPI[*LoginReq, *LoginRes]{
		ModuleName:  Gate.Name,
		Name:        "Login",
		MessageType: Gate.Name + "." + "Login",
	}

	Offline = &ServiceAPI[*OfflineReq, *OfflineRes]{
		ModuleName:  Gate.Name,
		Name:        "Offline",
		MessageType: Gate.Name + "." + "Offline",
	}
)

type UploadUserNodeInfoReq struct {
	UserId      string
	NodeId      string
	ServicePath string
}

type UploadUserNodeInfoRes struct {
	ErrorCode int32
}

type ToUserReq struct {
	Data        []byte
	UserId      []string
	MessageType string
}

type ToUserRes struct {
	ErrorCode int32
}

type LoginReq struct {
	UserId         string
	TemplateUserId string
}

type LoginRes struct {
	ErrorCode int32
}

type OfflineReq struct {
	UserId string
	//是否重复登录踢人行为
	Replace bool
}

type OfflineRes struct {
	ErrorCode int32
}

type DefaultArgs struct {
	C string
}

type DefaultReply struct {
	C int32
}

type DefaultBool struct {
	C bool
}

type EmptyReply struct {
}
