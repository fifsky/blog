#!/usr/bin/env bash
set -e

deploy() {
  local list=$(util::get_depoly_list $1)
  if [ -z "${list[0]}" ]; then # Check if the first element of the list is an empty string
    echo "没有匹配的项目"
    return
  fi

  local items=()
  # 正确处理项目名称，特别是包含空格的情况
  for p_name in $list; do
    items+=("$p_name")
  done

  if [ "${#items[@]}" -eq 1 ]; then
    selected=${items[0]}
    _rundeploy ${selected} --tail
    exit
  fi

  util::select "请选择部署项目" "${list[*]}"
  _rundeploy ${selected} --tail
}

deployall() {
  local list=$(util::get_depoly_list $1)
  for p_name in ${list[*]}; do
    _rundeploy ${p_name}
  done
}

_rundeploy() {
  if [ ! -d "$1" ]; then
    log::fatal "部署项目不存在: $1"
  fi

  # 如果是 web 项目，先执行构建
  if [[ "$1" == "deploy/web" ]]; then
    log::info "检测到 web 项目，正在执行构建..."
    make buildui
  fi

  log::info "$(date "+%Y-%m-%d %H:%M:%S") 开始部署: $1"
  uname=$(uname -m)
  if [ "$uname" == "arm64" ]; then
    skaffold run -f "$1/skaffold.yaml" $2 --platform=linux/arm64
  else
    skaffold run -f "$1/skaffold.yaml" $2 --platform=linux/amd64
  fi
}

main() {
  case $1 in
  all)
    deployall $2
    ;;
  *)
    deploy $1
    ;;
  esac
}

main $@
