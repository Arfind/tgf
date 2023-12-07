package examplechat

import (
	"fmt"
	"github.com/thkhxm/tgf/log"
	"github.com/thkhxm/tgf/rpc"
	"github.com/thkhxm/tgf/service/api/chat"
	"golang.org/x/net/context"
)

//***************************************************
//@Link  https://github.com/thkhxm/tgf
//@Link  https://gitee.com/timgame/tgf
//@QQ群 7400585
//author tim.huang<thkhxm@gmail.com>
//@Description
//2023/3/2
//***************************************************

type ChatService struct {
	rpc.Module
}

func (this *ChatService) RPCSayHello(ctx context.Context, req *string, response *chatapi.SayHelloRes) error {
	var (
		userId = rpc.GetUserId(ctx)
		msg    = *req
	)
	log.DebugTag("example", "RPCSayHello userId=%v ,msg=%v", userId, msg)
	response.Msg = fmt.Sprintf("%v say %v", userId, msg)
	//time.Sleep(time.Second * 10)
	return nil
}

func (this *ChatService) GetName() string {
	return chatapi.ChatService.Name
}

func (this *ChatService) GetVersion() string {
	return chatapi.ChatService.Version
}

func (this *ChatService) Startup() (bool, error) {
	var ()
	return true, nil
}

func (this *ChatService) Destroy(sub rpc.IService) {
	var ()
	log.Info("sub service destroy logic module=%v version %v", this.GetName(), this.GetVersion())
}
