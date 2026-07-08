#!/bin/bash

set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "用法: $0 <config.yaml>"
  exit 2
fi

GO_MAIN="${GO_MAIN:-./rime-dict-sync}"
CONFIG_PATH="$1"

"${GO_MAIN}" -c "${CONFIG_PATH}"
