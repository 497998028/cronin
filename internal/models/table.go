package models

import (
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// mysql最低版本
// 5.7.8 开始支持json函数，低于程序会报错
var mysqlLower = []int{5, 7, 7}

// 注册表结构
func AutoMigrate(db *db.MyDB) {
	if err := mysqlLowerCheck(db); err != nil {
		panic(err.Error())
	}

	// 迁移表结构
	err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").
		AutoMigrate(&CronSetting{}, &CronConfig{}, &CronPipeline{}, &CronLogSpan{}, &CronUser{}, &CronAuthRole{})
	if err != nil {
		panic(fmt.Sprintf("mysql 表初始化失败，%s", err.Error()))
	}
	// 初始化数据
	err = db.Where("scene=? and status=?", SceneEnv, enum.StatusActive).FirstOrCreate(&CronSetting{
		Scene:    "env",
		Name:     "public",
		Title:    "public",
		Content:  `{"default":2}`,
		Status:   enum.StatusActive,
		CreateDt: time.Now().Format(time.DateTime),
		UpdateDt: time.Now().Format(time.DateTime),
	}).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
	msg := &CronSetting{}
	err = db.Where("scene=?", SceneMsg).Find(msg).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
	if msg.Id == 0 { // 后期会有多条默认消息模板
		db.CreateInBatches([]*CronSetting{
			{
				Scene: "msg",
				Name:  "",
				Title: "企微xx群",
				Content: `{
	"http":{
		"method":"POST",
		"url":"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xx",
		"body":"{
			\"msgtype\": \"text\",
			\"text\": {
				\"content\": \"时间：[[log.create_dt]]\\n任务 [[config.name]]执行[[log.status_name]]了 \\n耗时[[log.duration]]秒\\n响应：[[log.body]]\",
				\"mentioned_mobile_list\": [[user.mobile]]
			}
		}",
		"header":[{"key":"","value":""}]
	}
}`,
				Status:   enum.StatusActive,
				CreateDt: time.Now().Format(time.DateTime),
				UpdateDt: time.Now().Format(time.DateTime),
			},
		}, 10)
	}

	// 历史 数据修正
	historyDataRevise(db)
}

// mysql 最低版本检测
func mysqlLowerCheck(db *db.MyDB) error {
	version := ""
	err := db.Raw("SELECT VERSION()").Scan(&version).Error
	if err != nil {
		return fmt.Errorf("mysql 版本获取失败，%s", err.Error())
	}

	temp1 := strings.Split(version, "-")
	temp2 := strings.Split(temp1[0], ".")
	isLower := true
	for i, n := range temp2 {
		val, _ := strconv.Atoi(n)
		if mysqlLower[i] < val {
			isLower = false
			break
		}
	}
	if isLower {
		return fmt.Errorf("mysql最低要求版本 5.7.8 当前为 %s", version)
	}
	return nil
}

// 历史源修正
// 解决 0.6.1 之前的版本格式不一致问题
func historyDataRevise(db *db.MyDB) {
	set := &CronSetting{}
	db.Where("scene='sys_tag_history_update'").Find(set)
	if set.Scene != "" {
		return // 存在表示已经修复过了
	}

	// sql_source 历史数据修正
	if err := db.Exec(`UPDATE cron_setting set content=concat('{"sql":',content,'}') WHERE scene='sql_source' and content->'$.sql' is null;`).Error; err != nil {
		log.Println("历史 sql_source 数据修正错误", err.Error())
	}

	// config cmd 历史数据修正
	cmdType := "sh"
	if runtime.GOOS == "windows" {
		cmdType = "cmd"
	}
	err := db.Exec(fmt.Sprintf(`UPDATE cron_config SET command=JSON_REPLACE(command,'$.cmd', CAST(concat('{"type":"%s","origin":"local","statement":{"type":"local","git":{},"local":',command->'$.cmd','}}') as JSON)) WHERE JSON_TYPE(command->'$.cmd') = 'STRING';`, cmdType)).Error
	if err != nil {
		log.Println("历史 config cmd 数据修正错误", err.Error())
	}
	// config sql 历史数据修正
	list := []*CronConfig{}
	db.Where("JSON_TYPE(command->'$.sql') = 'OBJECT' and command->'$.sql.origin' is null").Select("id", "command").Find(&list)
	if len(list) > 0 {
		type CronSql struct {
			Statement []string `json:"statement"` // sql语句多条
		}
		type CronConfigCommand struct {
			Sql *CronSql `json:"sql"`
		}
		for _, item := range list {
			cmd := &CronConfigCommand{Sql: &CronSql{Statement: []string{}}}
			if er := jsoniter.Unmarshal(item.Command, cmd); err != nil {
				log.Println("	sql 解析错误", item.Id, er.Error())
				continue
			}
			newStatement := make([]map[string]string, len(cmd.Sql.Statement))
			for i, statement := range cmd.Sql.Statement {
				newStatement[i] = map[string]string{
					"type":  "local",
					"local": statement,
				}
			}
			str, _ := jsoniter.MarshalToString(newStatement)
			updateSql := `UPDATE cron_config set command=JSON_SET(command, '$.sql.origin', 'local', '$.sql.statement', CAST(? as JSON)) WHERE id=?`
			if er := db.Exec(updateSql, str, item.Id).Error; er != nil {
				log.Println("	sql 修正错误: ", updateSql, er.Error())
			}
		}
	}
	set.Scene = "sys_tag_history_update"
	set.Content = `{"version":"0.6.1"}`
	db.Create(set)

}
