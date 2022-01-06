package common

import (
	"go.uber.org/zap"
	"log"
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
			log.Fatal("旧文件删除失败", zap.Error(err))
			//此处不能用util的logger，循环引入会导致编译报错
			//util.Logger.Fatal("旧文件删除失败",zap.Error(err))
			os.Exit(444)
		}
	}
}
