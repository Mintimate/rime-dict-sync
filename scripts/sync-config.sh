#!/bin/bash

set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "用法: $0 <config.yaml>"
  exit 2
fi

GO_MAIN="${GO_MAIN:-./rime-dict-sync}"
CONFIG_PATH="$1"

if [ ! -x "${GO_MAIN}" ]; then
  if [ "${GO_MAIN}" = "./rime-dict-sync" ]; then
    if ! go build -o "${GO_MAIN}" .; then
      echo "构建同步工具失败" >&2
      exit 2
    fi
  else
    echo "找不到可执行文件: ${GO_MAIN}"
    exit 2
  fi
fi

"${GO_MAIN}" -c "${CONFIG_PATH}"
