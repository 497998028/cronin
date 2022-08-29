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
	List []*CronConfigListItem `json:"list"`
	Page *Page                 `json:"page"`
}
type CronConfigListItem struct {
	Id           int                `json:"id"`
	Name         string             `json:"name"`
	Spec         string             `json:"spec"`
	Protocol     int                `json:"protocol"`
	ProtocolName string             `json:"protocol_name"`
	Remark       string             `json:"remark"`
	Status       int                `json:"status"`
	StatusName   string             `json:"status_name"`
	UpdateDt     string             `json:"update_dt"`
	Command      *CronConfigCommand `json:"command" gorm:"-"`
	CommandStr   string             `json:"-" gorm:"column:command;"` // 这里只能读取字符串后，载入到结构体
}

// 任务设置
type CronConfigSetRequest struct {
	Id       int                `json:"id,omitempty"`       // 主键
	Name     string             `json:"name,omitempty"`     // 任务名称
	Spec     string             `json:"spec"`               // 执行时间表达式
	Protocol int                `json:"protocol,omitempty"` // 协议：1.http、2.grpc、3.系统命令
	Command  *CronConfigCommand `json:"command,omitempty"`  // 命令
	Status   int                `json:"status"`             // 状态
	Remark   string             `json:"remark"`
}
type CronConfigSetResponse struct {
	Id int `json:"id"`
}
type CronConfigCommand struct {
	Http struct {
		Method string `json:"method"`
		Url    string `json:"url"`
		Body   string `json:"body"`
	} `json:"http"`
	Rpc string `json:"rpc"`
	Cmd string `json:"cmd"`
}
