package pb

type Page struct {
	Size  int   `json:"size"`
	Page  int   `json:"page"`
	Total int64 `json:"total"`
}

// 任务列表
type CronConfigListRequest struct {

}
type CronConfigListReply struct {
	List []*CronConfigListItem	`json:"list"`
	Page *Page              `json:"page"`
}
type CronConfigListItem struct {
	Id int
	Name string
	Protocol int
	Remark string
	Status int
	StatusName string
	UpdateDt string
}

// 任务设置
type CronConfigSetRequest struct {
	Id       int          `json:"id,omitempty"`        // 主键
	Name     string       `json:"name,omitempty"`      // 任务名称
	Protocol int `json:"protocol,omitempty"`  // 协议：1.http、2.grpc、3.系统命令
	Command  string       `json:"command,omitempty"`   // 命令
}
type CronConfigSetResponse struct {
	Id int	`json:"id"`
}