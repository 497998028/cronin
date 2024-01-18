package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 设置 sql 连接源 列表
func routerUserList(ctx *gin.Context) {
	r := &pb.UserListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewUserService(ctx.Request.Context()).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerUserSet(ctx *gin.Context) {
	r := &pb.UserSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewUserService(ctx.Request.Context()).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}