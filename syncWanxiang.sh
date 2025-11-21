#!/bin/bash

# 编译好的可执行文件路径
GO_MAIN="./rime-dict-sync"

# config.yaml配置文件内容
CONFIG_FILE_CONTENT='''
# 下载的词库，保存名称-下载地址
DOWNLOAD_DIR: "dl_dicts"
TARGET_DICT:
  - name: "rime_mint.base.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/dicts/基础词库.dict.yaml"
    remote_repo: "https://github.com/Mintimate/oh-my-rime.git"
    remote_path: "dicts/rime_mint.base.dict.yaml"
  - name: "rime_mint.chars.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/dicts/单字表.dict.yaml"
    remote_repo: "https://github.com/Mintimate/oh-my-rime.git"
    remote_path: "dicts/rime_mint.chars.dict.yaml"
  - name: "rime_mint.correlation.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/dicts/错音错字.dict.yaml"
    remote_repo: "https://github.com/Mintimate/oh-my-rime.git"
    remote_path: "dicts/rime_mint.correlation.dict.yaml"
  - name: "rime_mint.ext.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/dicts/联想词库.dict.yaml"
    remote_repo: "https://github.com/Mintimate/oh-my-rime.git"
    remote_path: "dicts/rime_mint.ext.dict.yaml"
'''

# 保存config.yaml的路径
CONFIG_FILE_PATH="./config.yaml"

cat <<EOF > ${CONFIG_FILE_PATH}
${CONFIG_FILE_CONTENT}
EOF

# 运行可执行文件
${GO_MAIN} -c ${CONFIG_FILE_PATH}

# 检查退出码
if [ $? -eq 0 ]; then
    echo "词库同步完成，有变化需要更新"
    exit 0
else
    echo "无变化或执行失败"
    exit 1
fi