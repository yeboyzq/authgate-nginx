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

# 定义变量(由调用方通过环境变量注入, 均带兜底默认值)
TARGETOS="${TARGETOS:-linux}"
TARGETARCH="${TARGETARCH:-$(dpkg --print-architecture)}"
APP="${APP:-authgate-nginx}"
OUTPUT_FILE="${OUTPUT_FILE:-./build/${APP}}"
CGO_ENABLED="${CGO_ENABLED:-1}"

# 获取Go模块路径(从go.mod解析)
get_module_path() {
    local path
    # 优先使用 go list -m (最权威, 依赖go环境)
    path=$(go list -m 2>/dev/null) || {
        # 降级到文本解析(纯grep+awk, 不依赖go命令)
        path=$(grep -m1 '^module ' go.mod | awk '{print $2}')
    }
    if [[ -z "$path" ]]; then
        echo "ERROR: 无法从 go.mod 获取模块路径" >&2
        exit 1
    fi
    echo "$path"
}

# 获取编译时间
get_build_time() {
    date "+%Y-%m-%dT%H:%M:%S%:z"
}

# 获取编译版本(优先tag, 否则分支名)
get_build_version() {
    local tag
    tag=$(git tag --points-at HEAD 2>/dev/null | head -1)
    if [[ -n $tag ]]; then
        echo "$tag"
    else
        git branch --show-current 2>/dev/null || echo "unknown"
    fi
}

# 获取编译信息
get_build_info() {
    local branch info
    branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    info=$(git describe --tags --always --dirty --long 2>/dev/null || echo "no-git")
    echo "${branch}:${info}"
}

# 获取编译CPU核心数
get_cpu_core() {
    local num
    if command -v nproc >/dev/null 2>&1; then
        num=$(nproc)
    else
        num=$(getconf _NPROCESSORS_ONLN 2>/dev/null || echo 1)
    fi
    if [[ $num -gt 6 ]]; then
        echo $((num - 2))
    elif [[ $num -gt 2 ]]; then
        echo $((num - 1))
    else
        echo "$num"
    fi
}

# 按目标架构选择 CGO 交叉编译器(目标=构建机本机架构时用本机 gcc)
get_cc() {
    local build_arch="$1" target_arch="$2"
    if [[ "$target_arch" == "$build_arch" ]]; then
        echo "gcc"
        return
    fi
    case "$target_arch" in
        arm64) echo "aarch64-linux-gnu-gcc" ;;
        amd64) echo "x86_64-linux-gnu-gcc" ;;
        *) echo "gcc" ;;
    esac
}

# 编译目标平台
build_target() {
    local module version info build_time build_arch cc core
    # 模块路径优先环境变量, 否则从 go.mod 解析
    module="${MODULE:-$(get_module_path)}"
    version=$(get_build_version)
    info=$(get_build_info)
    build_time=$(get_build_time)
    build_arch=$(dpkg --print-architecture)
    cc=$(get_cc "$build_arch" "$TARGETARCH")
    core=$(get_cpu_core)

    echo "INFO: 目标平台 ${TARGETOS}/${TARGETARCH}, 构建机 ${build_arch}, CGO=${CGO_ENABLED}, CC=${cc}, 核心数 ${core}"
    echo "INFO: 版本 ${version}, 编译信息 ${info}, 编译时间 ${build_time}"

    mkdir -p "$(dirname "$OUTPUT_FILE")"

    CGO_ENABLED="$CGO_ENABLED" GOOS="$TARGETOS" GOARCH="$TARGETARCH" CC="$cc" \
        go build -v -trimpath -p "$core" \
        -ldflags "-s -w \
            -X '${module}/app/utils.version=${version}' \
            -X '${module}/app/utils.build_info=${info}' \
            -X '${module}/app/utils.build_time=${build_time}'" \
        -o "$OUTPUT_FILE" ./app

    echo "INFO: 编译完成 -> ${OUTPUT_FILE}"
}

# 主函数
main() {
    build_target
}

main "$@"
