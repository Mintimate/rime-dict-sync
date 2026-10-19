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
3. 使用 `configs/` 中的任务配置，或按配置说明新增 YAML
4. 运行程序：
   ```bash
   go mod tidy
   go build .
   bash scripts/sync-config.sh configs/wanxiang.yaml
   ```

## 配置文件说明

普通 `.dict.yaml` 词库使用 `configs/*.yaml` 配置，格式如下：

```yaml
# 下载的词库，保存名称-下载地址
DOWNLOAD_DIR: "dl_dicts"
TARGET_DICT:
  - name: "rime_mint.base.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/jichu.dict.yaml"
  - name: "rime_mint.chars.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/zi.dict.yaml"
  - name: "rime_mint.correlation.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/cuoyin.dict.yaml"
  - name: "rime_mint.compatible.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/duoyin.dict.yaml"
  - name: "rime_mint.places.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/diming.dict.yaml"
  - name: "rime_mint.artists.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/yiren.dict.yaml"
  - name: "rime_mint.famous.dict.yaml"
    url: "https://raw.githubusercontent.com/amzxyz/rime-wanxiang-data/refs/heads/wanxiang/dicts/mingren.dict.yaml"
```

万象词库从 [rime-wanxiang-data](https://github.com/amzxyz/rime-wanxiang-data) 的 `wanxiang` 分支同步，映射如下：

| 上游词库 | 输出词库 |
| --- | --- |
| `jichu` | `rime_mint.base` |
| `zi` | `rime_mint.chars` |
| `cuoyin` | `rime_mint.correlation` |
| `duoyin` | `rime_mint.compatible` |
| `diming` | `rime_mint.places` |
| `yiren` | `rime_mint.artists` |
| `mingren` | `rime_mint.famous` |

上游于 2026-10-15 取消独立的 `lianxiang.dict.yaml`，因此不再同步 `rime_mint.ext.dict.yaml`。本次 oh-my-rime 将 ext 改为空兼容占位文件，保留旧用户 `import_tables` 中的 `rime_mint.ext` 引用仍可正常加载。新配置可移除该引用；只有确认所有引用已移除后，才可删除占位文件。工具不会自动清理目标仓库中的退役词库。上游对内容做了筛选与重组，不保证旧 ext 中的每条词条都进入基础库。

## 使用示例

运行程序后，输出示例：

```txt
正在下载: rime_mint.base.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.base.dict.yaml
正在下载: rime_mint.chars.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.chars.dict.yaml
正在下载: rime_mint.correlation.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.correlation.dict.yaml
```

下载的字典文件将保存在 `DOWNLOAD_DIR` 指定的目录中，默认是 `dl_dicts`。它是本次同步的输出暂存目录；本地运行不会自动清理已有文件，迁移或切换配置前应先清理旧输出。CI 的全部同步入口会在运行前清空 `dl_dicts`，以免通配复制时带上已退役或其他任务残留的文件；请勿将个人文件存放在此目录。

程序及 `scripts/sync-config.sh` 的退出码为：`0` 表示检测到变化，`1` 表示无变化，`2` 表示参数、构建、下载或处理失败。CI 仅把 `1` 转为跳过状态 `78`，其他失败保留原退出码并中止流水线。自定义同步脚本也必须遵守此约定，并将输出写入 `dl_dicts`。

需要保留目标仓库头部注释、只同步正文的特殊词库，可以使用 `merge_mode: "keep_remote_header"`，例如：

```bash
bash scripts/sync-config.sh configs/rime-ice-english.yaml
```

## 许可证

本项目使用 MIT 许可证。详情请查看 [LICENSE](./LICENSE) 文件。
