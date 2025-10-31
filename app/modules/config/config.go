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
	"log/slog"
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
	if CfgFile != "" {
		// 从命令行加载配置文件
		Cfg.SetConfigFile(CfgFile)
	} else {
		// 从环境变量中读取
		Cfg.SetEnvPrefix("app")
		Cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		Cfg.AutomaticEnv()
		// 定义文件名
		Cfg.SetConfigName("config.yaml")
		// 定义文件类型
		Cfg.SetConfigType("yaml")
		// 定义查找路径
		// Cfg.AddConfigPath(os.Getenv("APP_CONF_PATH"))
		Cfg.AddConfigPath("./conf.d")
		Cfg.AddConfigPath(".")
	}
	// 查找并读取配置文件
	err := Cfg.ReadInConfig()
	// 读取文件错误处理
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("配置初始化异常: 配置文件没有找到, 将以默认配置启动.")
		} else {
			slog.Error("配置初始化失败: " + err.Error())
		}
	} else {
		configFile := Cfg.ConfigFileUsed()
		slog.Info("配置初始化完成.")
		slog.Info("当前配置文件: " + configFile)
	}
	return Cfg
}

// dynamicConfig 动态加载配置
func dynamicConfig() {
	Cfg.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("配置文件发生变更: " + e.Name)
	})
	Cfg.WatchConfig()
}
