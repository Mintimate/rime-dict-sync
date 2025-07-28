#!/bin/bash

# 编译好的可执行文件路径
GO_MAIN="./rime-dict-sync"

# config.yaml配置文件内容
CONFIG_FILE_CONTENT='''
# 下载的词库，保存名称-下载地址
DOWNLOAD_DIR: "dl_dicts"
TARGET_DICT:
  - name: "radical_pinyin.dict.yaml"
    url: "https://github.com/mirtlecn/rime-radical-pinyin/raw/master/radical_pinyin.dict.yaml"
'''

# 保存config.yaml的路径
CONFIG_FILE_PATH="./config.yaml"

cat <<EOF > ${CONFIG_FILE_PATH}
${CONFIG_FILE_CONTENT}
EOF


# 运行可执行文件
${GO_MAIN} -c ${CONFIG_FILE_PATH}