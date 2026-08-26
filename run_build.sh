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

# 定义变量
OUTPUT_DIR="./build"
OUTPUT_FILE_PREFIX="authgate-nginx"
# 定义目标平台矩阵(GOOS/GOARCH), 留空则默认构建当前平台
TARGETS=(
    "linux/amd64"
    # "linux/arm64"
    "windows/amd64"
    # "darwin/amd64"
    # "darwin/arm64"
)

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

# 清理版本号中的特殊字符(用于文件名, 分支名可能含/)
safe_version() {
    local version="$1"
    echo "${version//\//-}"
}

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

# 确保目标平台非空(空则填充当前平台)
ensure_targets() {
    if [[ ${#TARGETS[@]} -eq 0 ]]; then
        TARGETS=("$(go env GOOS)/$(go env GOARCH)")
        echo "INFO: 未指定目标平台, 默认构建当前平台: ${TARGETS[0]}"
    fi
}

# 构建单个目标平台
build_target() {
    local goos="$1" goarch="$2" version="$3" info="$4" build_time="$5" module="$6" core="$7"
    local output_file="${OUTPUT_DIR}/${OUTPUT_FILE_PREFIX}-${version}-${goos}-${goarch}"
    # Windows可执行文件加.exe后缀
    if [[ "$goos" == "windows" ]]; then
        output_file="${output_file}.exe"
    fi

    echo "INFO: 开始构建 ${goos}/${goarch} -> ${output_file} ..."
    CGO_ENABLED=1 GOOS="$goos" GOARCH="$goarch" \
        go build -v -p "$core" \
        -ldflags "-X '${module}/app/utils.version=${version}' -X '${module}/app/utils.build_info=${info}' -X '${module}/app/utils.build_time=${build_time}'" \
        -o "$output_file" ./app
    echo "INFO: 完成 ${goos}/${goarch} -> ${output_file}"
}

# 列出构建产物(容错处理)
list_outputs() {
    echo "INFO: 产物列表:"
    ls -lh "${OUTPUT_DIR}/${OUTPUT_FILE_PREFIX}"-* 2>/dev/null || echo "WARN: 无构建产物"
}

# 主函数
main() {
    ensure_targets

    # 收集编译信息(循环外统一获取, 保证多平台版本信息一致)
    local module_path cpu_core build_version build_info build_time
    module_path=$(get_module_path)
    cpu_core=$(get_cpu_core)
    build_version=$(safe_version "$(get_build_version)")
    build_info=$(get_build_info)
    build_time=$(get_build_time)

    echo "INFO: 目标平台数: ${#TARGETS[@]}, 构建核心数: ${cpu_core}, 构建版本号: ${build_version}, 构建信息: ${build_info}, 构建时间: ${build_time}"

    # 确保输出目录存在
    mkdir -p "$OUTPUT_DIR"

    # 循环构建每个目标平台
    local target
    for target in "${TARGETS[@]}"; do
        build_target "${target%/*}" "${target#*/}" "$build_version" "$build_info" "$build_time" "$module_path" "$cpu_core"
    done

    echo "INFO: 全部编译完成, 版本: ${build_version}"
    list_outputs
}

main "$@"