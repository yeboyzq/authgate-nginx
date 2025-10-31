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

	"github.com/yeboyzq/authgate-nginx/app/utils"

	"github.com/spf13/cobra"
)

// 根入口
var rootCmd = &cobra.Command{
	Use:   utils.AppFileName(),
	Short: "AuthGate-Nginx",
	Long:  "AuthGate-Nginx",
	Run: func(cmd *cobra.Command, args []string) {
		rootFlag()
	},
}

// rootFlag 根标志处理
func rootFlag() {
	if versionFlag {
		fmt.Printf("Version: %s\nBuildTime: %s\n", utils.VersionInfo.AppVersion, utils.VersionInfo.BuildTime)
		os.Exit(0)
	}
}
