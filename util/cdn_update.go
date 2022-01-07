package util

import (
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/common"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func CdnDownload(filePath string) (data []byte, err error) {
	data, err = CdnGetData()
	if err != nil {
		Logger.Info("CDN数据库下载失败，请手动下载解压后保存到本地: ", zap.String("", filePath))
		Logger.Info("\n下载链接：", zap.String("",githubUrl))
		return
	}

	common.ExistThenRemove(filePath)
	if err := ioutil.WriteFile(filePath, data, 0644); err == nil {
		Logger.Info("CDN数据库更新完毕，", zap.String("本地路径: ", filePath))
	}
	return
}

const (
	githubUrl   = "https://raw.githubusercontent.com/SukkaLab/cdn/master/dist/cdn.json"
	jsdelivrUrl = "https://cdn.jsdelivr.net/gh/SukkaLab/cdn/dist/cdn.json"
)

func CdnGetData() (data []byte, err error) {
	resp, err := http.Get(jsdelivrUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp, err = http.Get(githubUrl)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
