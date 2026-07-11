package services

import "testing"

func Test_Reader(t *testing.T) {
	filepath := "/Users/charlie/Desktop/服务升级人员挂靠及训练战赛规划 (2).xlsx"
	reader := &Reader{}
	reader.read("id", filepath)
}
