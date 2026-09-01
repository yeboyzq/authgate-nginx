/*
Copyright (c) 2025 authgate-nginx
authgate-nginx is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
        http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Cfg     *viper.Viper // 配置实例
	CfgFile string       // 配置文件
)

// Init 初始化配置
func Init() {
	Cfg = readConfig()
	defaultConfig(Cfg)
	go dynamicConfig()
}

// readConfig 读取配置
func readConfig() *viper.Viper {
	Cfg = viper.New()

	// 启用环境变量读取，环境变量优先于配置文件
	Cfg.SetEnvPrefix("app")
	Cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Cfg.AutomaticEnv()

	// 配置文件优先级：
	// 1. 命令行参数 CfgFile
	// 2. 环境变量 APP_CONF_FILE
	// 3. 默认配置文件 ./custom/conf/config.yaml
	var configPath string
	if CfgFile != "" {
		configPath = CfgFile
	} else if envCfgFile := os.Getenv("APP_CONF_FILE"); envCfgFile != "" {
		configPath = envCfgFile
	} else {
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exeDir := filepath.Dir(exePath)    // ./bin
		projectDir := filepath.Dir(exeDir) // ./
		configPath = filepath.Join(projectDir, "custom", "conf", "config.yaml")
	}

	Cfg.SetConfigFile(configPath)
	if err := Cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("WARN: 配置文件未找到: %v, 继续使用环境变量或默认值\n", configPath)
		} else {
			fmt.Printf("ERROR: 配置初始化失败: %v\n", err)
		}
	} else {
		fmt.Printf("INFO: 配置初始化完成.\n")
		fmt.Printf("当前配置文件: %v\n", Cfg.ConfigFileUsed())
	}
	return Cfg
}

// dynamicConfig 动态加载配置
func dynamicConfig() {
	Cfg.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("INFO: 配置文件发生变更: %v\n", e.Name)
	})
	Cfg.WatchConfig()
}
