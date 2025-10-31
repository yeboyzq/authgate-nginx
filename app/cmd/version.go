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

	"github.com/yeboyzq/authgate-nginx/app/utils"

	"github.com/spf13/cobra"
)

var versionFlag bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看应用版本信息",
	Long:  "查看应用版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nBuilt By: %s\nBuild Time: %s\n", utils.VersionInfo.AppVersion, utils.VersionInfo.BuiltBy, utils.VersionInfo.BuildTime)
	},
}
