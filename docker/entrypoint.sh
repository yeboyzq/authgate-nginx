# Copyright (c) 2026 authgate-nginx
# authgate-nginx is licensed under Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#         http://license.coscl.org.cn/MulanPSL2
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
# EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
# MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
# See the Mulan PSL v2 for more details.

set -euo pipefail

# 应用可执行文件
APP_BIN="/data/apps/authgate-nginx"
# 配置文件路径(优先 APP_CONF_FILE 环境变量, 兜底默认值)
CONF_FILE="${APP_CONF_FILE:-/data/apps/custom/conf/config.yaml}"

echo " 工作目录: $(pwd)"
echo " 配置文件: ${CONF_FILE}"

# 环境检查
echo "[1/2] 环境检查..."
if [[ ! -x "${APP_BIN}" ]]; then
    echo "错误: 应用可执行文件不存在或不可执行: ${APP_BIN}"
    exit 1
fi
if [[ ! -f "${CONF_FILE}" ]]; then
    echo "错误: 配置文件不存在: ${CONF_FILE}"
    echo "     请通过 configmap 挂载到该路径, 或设置 APP_CONF_FILE 环境变量"
    exit 1
fi
# 确保运行期可写目录存在
mkdir -p /data/logs /data/apps/custom/db
chown -R 1000:1000 /data/logs && chmod -R 755 /data/logs
echo "环境检查通过."

# 启动应用
echo "[2/2] 启动应用..."
exec "${APP_BIN}" start --config "${CONF_FILE}"
