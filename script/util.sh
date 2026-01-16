#!/usr/bin/env bash

# get relative path from current script
# returns an path
function util::rpath() {
  local path="$1"
  local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && cd .. && pwd)"
  echo ${path} | sed "s%${script_dir}/%%g"
}

# select a random item from an array
function util::select() {
  local arr=$2
  trap 'selected=""; return 0' INT TERM
  # 如果安装了gum，则使用gum来选择
  if command -v gum &>/dev/null; then
    selected=$(gum filter --header="**************** $1 ****************" ${arr} || true)
    return
  elif command -v fzf &>/dev/null; then
    selected=$(printf '%s\n' ${arr[@]} | fzf --height 40% --reverse --header="**************** $1 ****************" || true)
    return
  else
    log::info "**************** $1 ****************"
    # shellcheck disable=SC2068
    select selected in ${arr[@]}; do
      # shellcheck disable=SC2076
      if [[ "${arr[*]}" =~ "${selected}" ]]; then
        log::info "您选择了: ${selected}"
        break
      else
        log::info "请输入编号"
      fi
    done
  fi
}

# get depoly list
function util::get_depoly_list() {
  find deploy -name Dockerfile | xargs -I F dirname "F" | grep "deploy/$1"
}

# Log an error but keep going.  Don't dump the stack or exit.
log::fatal() {
  for message; do
    echo -e "\033[41m ${message} \033[0m" >&2
  done
  exit 1
}

log::error() {
  for message; do
    echo -e "\033[41m ${message} \033[0m" >&2
  done
}

# Print out some info that isn't a top level status line
log::info() {
  for message; do
    echo -e "\033[32m ${message} \033[0m"
  done
}