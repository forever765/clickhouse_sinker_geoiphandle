package constant

import (
	"log"
	"os"
	"path/filepath"
)

var (
	// HomePath database home path
	HomePath string
)



// get it from config file
func init() {
	HomePath = os.Getenv("NALI_DB_HOME")
	if HomePath == "" {
		// homeDir, err := os.UserHomeDir()
		// homeDir, _ := os.UserHomeDir()
		// if err != nil {
		// 	// errLog := err.Error() + "\n"
		// 	log.Fatal(err.Error())
		// 	// util.Logger.Error("Get homePath error: ", zap.Error(err))
		// }
		HomePath = filepath.Join("/usr/share/ch_sinker/geoip_db")
	}
	if _, err := os.Stat(HomePath); os.IsNotExist(err) {
		if err := os.MkdirAll(HomePath, 0777); err != nil {
			log.Fatal("can not create", HomePath, ", use bin dir instead")
		}
	}
}
