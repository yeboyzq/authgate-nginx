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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yeboyzq/authgate-nginx/app/modules/config"
)

// Execute cmd入口
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init 初始化(隐式调用)
func init() {
	// 执行命令时运行配置
	cobra.OnInitialize(config.Init)

	// 注册标志
	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "V", false, "查看应用版本")
	rootCmd.PersistentFlags().StringVarP(&config.CfgFile, "config", "C", "", "配置文件(默认: ./conf/config.yaml)")
	startCmd.Flags().BoolVarP(&debugFlag, "debug", "D", false, "Debug模式")
	startCmd.Flags().IntVarP(&serverPortFlag, "port", "P", 8000, "监听端口")

	// 绑定标志到配置
	config.Cfg.BindPFlag("base.debug", rootCmd.PersistentFlags().Lookup("debug"))
	config.Cfg.BindPFlag("base.server.port", rootCmd.PersistentFlags().Lookup("port"))

	// 注册命令
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)

}
