#!/usr/bin/env bash
set -e
export selected

APP_HOME="$(cd "$(dirname "${0}")" && cd .. && pwd -P)"
cd ${APP_HOME}

source ${APP_HOME}/script/util.sh

main(){
  local name=${1}
  shift
  source ${APP_HOME}/script/${name}.sh
}

main $@
