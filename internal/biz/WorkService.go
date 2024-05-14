package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/pb"
)

type WorkService struct {
	db   *db.MyDB
	ctx  context.Context
	user *auth.UserToken
}

func NewWorkService(ctx context.Context, user *auth.UserToken) *WorkService {
	return &WorkService{
		ctx:  ctx,
		user: user,
		db:   db.New(ctx),
	}
}

// 工作表格
func (dm *WorkService) Table(r *pb.WorkTableRequest) (resp *pb.WorkTableReply, err error) {
	/*
		查询所有任务，
			目前仅做待审核的。
	*/
	resp = &pb.WorkTableReply{
		List: []*pb.WorkTableItem{},
	}

	sql := `SELECT COUNT(*) total, 'config' join_type, env FROM cron_config  GROUP BY env
UNION ALL
SELECT COUNT(*) total, 'pipeline' join_type, env FROM cron_pipeline  GROUP BY env`

	list := []*pb.WorkTableItem{}
	dm.db.Raw(sql).Scan(&list)
	if len(list) == 0 {
		return resp, nil
	}

	envs, err := NewDicService(dm.ctx, dm.user).getDb(enum.DicEnv)
	if err != nil {
		return nil, err
	}

	listMap := map[string]*pb.WorkTableItem{}
	for _, val := range list {
		listMap[val.Env+"|"+val.JoinType] = val
	}
	for _, env := range envs {
		if item, ok := listMap[env.Key+"|config"]; ok {
			item.EnvTitle = env.Name
			resp.List = append(resp.List, item)
		}
		if item, ok := listMap[env.Key+"|pipeline"]; ok {
			item.EnvTitle = env.Name
			resp.List = append(resp.List, item)
		}
	}
	return resp, nil
}
