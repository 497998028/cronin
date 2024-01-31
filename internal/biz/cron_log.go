package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/db"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronLogService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronLogService(ctx context.Context, user *auth.UserToken) *CronLogService {
	return &CronLogService{
		ctx:  ctx,
		user: user,
	}
}

// List 列表
func (dm *CronLogService) List(r *pb.CronLogListRequest) (resp *pb.CronLogListResponse, err error) {
	tags := map[string]any{}
	if err := jsoniter.UnmarshalFromString(r.Tags, &tags); err != nil {
		return nil, errs.New(err, "tags传递不规范")
	}

	w := db.NewWhere().In("env", []string{dm.user.Env, ""})
	for k, v := range tags {
		w.JsonIndexEq("tags_key", "tags_val", k, v)
	}

	list := []*models.CronLogSpan{}
	_, err = data.NewCronLogSpanData(dm.ctx).ListPage(w, 1, r.Limit, &list)
	resp = &pb.CronLogListResponse{List: make([]*pb.CronLogSpan, len(list))}
	for i, item := range list {
		resp.List[i] = dm.toOut(item)
	}

	return resp, err
}

// Trace 踪迹
func (dm *CronLogService) Trace(r *pb.CronLogTraceRequest) (resp *pb.CronLogTraceResponse, err error) {
	if r.TraceId == "" {
		return nil, errs.New(nil, "未指定traceId")
	}

	w := db.NewWhere().In("env", []string{dm.user.Env, ""}).Eq("trace_id", r.TraceId)
	list, err := data.NewCronLogSpanData(dm.ctx).List(w, 1000)

	// 树 或 列表；样例为树，那我也树吧。
	resp = &pb.CronLogTraceResponse{
		Data:  []*pb.CronLogTraceItem{},
		Limit: 1000,
		Total: len(list),
	}

	tra := &pb.CronLogTraceItem{
		TraceId: r.TraceId,
		Spans:   []*pb.CronLogSpan{},
	}
	for _, item := range list {
		span := dm.toOut(item)
		tra.Spans = append(tra.Spans, span)
	}
	resp.Data = append(resp.Data, tra)

	return resp, err
}

// Del 删除
func (dm *CronLogService) Del(r *pb.CronLogDelRequest) (resp *pb.CronLogDelResponse, err error) {
	if r.Retention == "" {
		return nil, fmt.Errorf("retention 参数为必须")
	}

	re, err := time.ParseDuration(r.Retention)
	if err != nil {
		return nil, fmt.Errorf("retention 参数有误, %s", err.Error())
	} else if re.Hours() < 24 {
		return nil, fmt.Errorf("retention 参数不得小于24h")
	}
	end := time.Now().Add(-re)
	resp = &pb.CronLogDelResponse{}
	resp.Count, err = data.NewCronLogData(dm.ctx).DelBatch(end)

	return resp, err
}

// 转输出
func (dm *CronLogService) toOut(in *models.CronLogSpan) *pb.CronLogSpan {
	out := &pb.CronLogSpan{
		Timestamp:    in.Timestamp,
		Duration:     in.Duration / 1000,
		Status:       in.Status,
		StatusName:   "",
		StatusDesc:   "",
		TraceId:      in.TraceId,
		SpanId:       in.SpanId,
		ParentSpanId: in.ParentSpanId,
		Service:      in.Service,
		Operation:    in.Operation,
		Tags:         []*pb.CronLogSpanTag{},
		Logs:         []*pb.CronLogSpanLog{},
	}

	jsoniter.Unmarshal(in.Tags, &out.Tags)
	jsoniter.Unmarshal(in.Logs, &out.Logs)

	return out
}
