package util

import (
	"github.com/forever765/clickhouse_sinker_nali/ipHandle/pkg/common"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/saracen/go7z"
)

func Zxipv6wry_Download(filePath string) (data []byte, err error) {
	data, err = zxipv6wryGetData()
	if err != nil {
		Logger.Info("ZX IPv6数据库下载失败，请手动下载解压后保存到本地:", zap.String("",filePath))
		Logger.Info("\n下载链接： https://ip.zxinc.org/ip.7z")
		return
	}
	common.ExistThenRemove(filePath)
	if err = ioutil.WriteFile(filePath, data, 0644); err == nil {
		Logger.Info("已将最新的 ZX IPv6数据库 保存到本地: ", zap.String("",filePath))
	}
	return

}

func zxipv6wryGetData() (data []byte, err error) {
	resp, err := http.Get("https://ip.zxinc.org/ip.7z")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	file7z, err := ioutil.TempFile("", "*")
	if err != nil {
		errLog := err.Error() + "\n"
		Logger.Error(errLog)
	}
	defer os.Remove(file7z.Name())

	if err := ioutil.WriteFile(file7z.Name(), body, 0644); err == nil {
		return Un7z(file7z.Name())
	}
	return
}

func Un7z(filePath string) (data []byte, err error) {
	sz, err := go7z.OpenReader(filePath)
	if err != nil {
		errLog := err.Error() + "\n"
		Logger.Error(errLog)
	}
	defer sz.Close()

	fileNoNeed, err := ioutil.TempFile("", "*")
	if err != nil {
		errLog := err.Error() + "\n"
		Logger.Error(errLog)
	}
	fileNeed, err := ioutil.TempFile("", "*")
	if err != nil {
		errLog := err.Error() + "\n"
		Logger.Error(errLog)
	}

	if err != nil {
		errLog := err.Error() + "\n"
		Logger.Error(errLog)
	}
	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			errLog := err.Error() + "\n"
			Logger.Error(errLog)
		}

		if hdr.Name == "ipv6wry.db" {
			if _, err := io.Copy(fileNeed, sz); err != nil {
				Logger.Fatal("ZX ipv6数据库解压出错：", zap.Error(err))
			}
		} else {
			if _, err := io.Copy(fileNoNeed, sz); err != nil {
				Logger.Fatal("ZX ipv6数据库解压出错：", zap.Error(err))
			}
		}
	}
	err = fileNoNeed.Close()
	if err != nil {
		errLog := err.Error() + "\n"
		Logger.Error(errLog)
	}
	defer os.Remove(fileNoNeed.Name())
	defer os.Remove(fileNeed.Name())
	return ioutil.ReadFile(fileNeed.Name())
}
