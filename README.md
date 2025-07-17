# Rime 字典自动同步工具

这个工具用于自动下载 Rime 输入法的字典文件，并更新其版本号到当前日期。

## 功能

- 从指定的URL下载字典文件
- 自动修改字典的`name`和`version`字段
  - `name`修改为字典文件名（不含扩展名）
  - `version`修改为当前日期（格式：YYYY-MM-DD）
- 将修改后的字典保存到指定目录

## 使用方法

1. 安装 Go 语言环境（1.24+）
2. 克隆仓库：
   ```bash
   git clone https://cnb.cool/Mintimate/rime/rime-dict-sync.git
   cd rime-dict-sync
   ```
3. 创建配置文件 `config.yaml`（参考下方配置说明）
4. 运行程序：
   ```bash
   go mod tidy
   go run .
   ```

## 配置文件说明

创建 `config.yaml` 文件，格式如下：

```yaml
# 下载的词库，保存名称-下载地址
DOWNLOAD_DIR: "dl_dicts"
TARGET_DICT:
  - name: "rime_mint.base.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/zh_dicts/base.dict.yaml"
  - name: "rime_mint.chars.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/zh_dicts/chars.dict.yaml"
  - name: "rime_mint.correlation.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/zh_dicts/correlation.dict.yaml"
  - name: "rime_mint.ext.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime_wanxiang/refs/heads/wanxiang/zh_dicts/suggestion.dict.yaml"
```

## 使用示例

运行程序后，输出示例：

```txt
正在下载: rime_mint.base.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.base.dict.yaml
正在下载: rime_mint.chars.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.chars.dict.yaml
正在下载: rime_mint.correlation.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.correlation.dict.yaml
正在下载: rime_mint.ext.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.ext.dict.yaml
```

下载的字典文件将保存在 `dl_dicts` 目录中，可以直接用于 Rime 输入法。