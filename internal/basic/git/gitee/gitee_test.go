package gitee

import (
	"context"
	"cron/internal/basic/git"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var conf = &git.Config{AccessToken: "e6a28b06d79d492f9809069d5550b436"}

func TestUrl(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())
	res, err := api.FileGet(handler, &FileGetRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Path: "work/mysql.sql",
		Ref:  "master",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(handler)
	fmt.Println(res)
	fmt.Println(res.DecodeContent())
}

func TestUser(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())
	res, err := api.User(handler)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(handler)
	fmt.Println(string(res))
}

func TestApiV5_PullsCreate(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsCreate(handler, &PullsCreateRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Head:                  "master",
		Base:                  "test",
		Title:                 "test demo",
		Body:                  "pr body .",
		MilestoneNumber:       0,
		Labels:                "",
		Issue:                 "",
		Assignees:             "",
		Testers:               "",
		AssigneesNumber:       0,
		TestersNumber:         0,
		RefPullRequestNumbers: "",
		PruneSourceBranch:     false,
		CloseRelatedIssue:     false,
		Draft:                 false,
		Squash:                false,
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

func TestApiV5_PullsReview(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsReview(handler, &PullsReviewRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Number: 9,
		Force:  false,
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

func TestApiV5_PullsTest(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsTest(handler, &PullsTestRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Number: 9,
		Force:  false,
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

// 合并分支
func TestPullsMerge(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsMerge(handler, &PullsMergeRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Number:            9,
		MergeMethod:       "merge",
		PruneSourceBranch: false,
		Title:             "A",
		Description:       "B",
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

func TestName(t *testing.T) {
	// 最后一个字符串+1
	// 层级为4，不足时补0
	str := `release_v3.5.87.2`
	parts := strings.Split(str, ".")
	lastNumString := parts[len(parts)-1]
	fmt.Println(parts, lastNumString)

	// 将字符串转换为数字并加 1
	lastNum, err := strconv.Atoi(lastNumString)
	if err != nil {
		panic(err)
	}
	lastNum++

	// 将数字转换回字符串并重新组装版本号
	newLastNumString := strconv.Itoa(lastNum)
	parts[len(parts)-1] = newLastNumString
	newVersion := strings.Join(parts, ".")

	fmt.Println(newVersion) // 输出：release_v3.5.87.3
}
