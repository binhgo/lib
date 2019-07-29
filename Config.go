package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	flag bool
}

func (this *Config) Init(modelName string) (error) {
	if this.flag {
		return nil
	}

	exePath, _ := os.Executable()
	dir, err1 := filepath.Abs(exePath + "/../../")
	if err1 != nil {
		log.Fatal(err1)
	}

	onlineConfPath := dir + "/bin/conf/" + modelName + ".yaml"
	fmt.Println(onlineConfPath)

	viper.SetConfigType("yaml")

	if _, err := os.Stat(onlineConfPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("dev confg mode")
			// 加载开发环境配置
			viper.SetConfigName(modelName + ".dev") // name of config file (without extension)
			// 同一个变量,优先出现,优先使用
			viper.AddConfigPath(dir + "/src/" + modelName + "/conf/")
		} else {
			panic("config error")
		}
	} else {
		fmt.Println("online config mode")
		viper.SetConfigName(modelName) // name of config file (without extension)
		viper.AddConfigPath(dir + "/bin/conf/")
	}

	err := viper.ReadInConfig() // Find and read the config file
	if (err != nil) {
		panic(err.Error())
	}

	this.flag = true

	return err
}
