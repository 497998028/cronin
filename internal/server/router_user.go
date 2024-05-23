package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 用户列表
func routerUserList(ctx *gin.Context) {
	r := &pb.UserListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewUserService(ctx.Request.Context()).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 用户设置
func routerUserSet(ctx *gin.Context) {
	r := &pb.UserSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}

	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	// 这里有三种情况：add.新增、edit.编辑所有、olay.编辑自己
	authType := ctx.GetString("auth_type")
	if authType == "add" { // 新增
		if r.Id > 0 {
			NewReply(ctx).SetError(pb.ParamError, "权限与操作不匹配").RenderJson()
			return
		}
	} else if authType == "edit" { // 编辑
		if r.Id <= 0 {
			NewReply(ctx).SetError(pb.ParamError, "权限与操作不匹配").RenderJson()
			return
		}
	} else { // 修改自己
		if r.Id != user.UserId {
			NewReply(ctx).SetError(pb.ParamError, "没有操作权限").RenderJson()
			return
		}
	}
	rep, err := biz.NewUserService(ctx.Request.Context()).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 用户详情
func routerUserDetail(ctx *gin.Context) {
	r := &pb.UserDetailRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewUserService(ctx.Request.Context()).Detail(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 修改密码
func routerUserChangePassword(ctx *gin.Context) {
	r := &pb.UserSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	//user, err := GetUser(ctx)
	//if err != nil {
	//	NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
	//	return
	//}
	rep, err := biz.NewUserService(ctx.Request.Context()).ChangePassword(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func routerUserChangeStatus(ctx *gin.Context) {
	r := &pb.UserChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	//user, err := GetUser(ctx)
	//if err != nil {
	//	NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
	//	return
	//}
	rep, err := biz.NewUserService(ctx.Request.Context()).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 设置账号
func routerUserChangeAccount(ctx *gin.Context) {
	r := &pb.UserSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	//user, err := GetUser(ctx)
	//if err != nil {
	//	NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
	//	return
	//}
	rep, err := biz.NewUserService(ctx.Request.Context()).ChangeAccount(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 用户登录
func routerUserLogin(ctx *gin.Context) {
	r := &pb.UserLoginRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewUserService(ctx.Request.Context()).Login(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}
