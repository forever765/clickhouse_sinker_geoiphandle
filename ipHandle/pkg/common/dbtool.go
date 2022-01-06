package common

import (
	"github.com/forever765/clickhouse_sinker_nali/util"
	"go.uber.org/zap"
	"os"
)

func ByteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

func ExistThenRemove(filePath string) {
	_, err := os.Stat(filePath)
	if err == nil {
		err = os.Remove(filePath)
		if err != nil {
			util.Logger.Fatal("旧文件删除失败", zap.Error(err))
			os.Exit(1)
		}
	}
}
