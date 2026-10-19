# Rime Dictionary Auto-Sync Tool

This tool automatically downloads dictionary files for Rime input method and updates their version to the current date.

## Features

- Downloads dictionary files from specified URLs
- Automatically modifies the `name` and `version` fields of dictionaries
  - `name` changed to the dictionary filename (without extension)
  - `version` updated to current date (format: YYYY-MM-DD)
- Saves modified dictionaries to specified directory

## Usage

1. Install Go environment (1.24+)
2. Clone repository:
   ```bash
   git clone https://cnb.cool/Mintimate/rime/rime-dict-sync.git
   cd rime-dict-sync
   ```
3. Use a task config from `configs/`, or add a YAML file following the configuration section
4. Run program:
   ```bash
   go mod tidy
   go build .
   bash scripts/sync-config.sh configs/wanxiang.yaml
   ```

## Configuration

Regular `.dict.yaml` dictionaries use `configs/*.yaml` files with the following format:

```yaml
# Dictionaries to download: save_name-download_url
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

Wanxiang dictionaries are synced from the `wanxiang` branch of [rime-wanxiang-data](https://github.com/amzxyz/rime-wanxiang-data), using these mappings:

| Upstream dictionary | Output dictionary |
| --- | --- |
| `jichu` | `rime_mint.base` |
| `zi` | `rime_mint.chars` |
| `cuoyin` | `rime_mint.correlation` |
| `duoyin` | `rime_mint.compatible` |
| `diming` | `rime_mint.places` |
| `yiren` | `rime_mint.artists` |
| `mingren` | `rime_mint.famous` |

Upstream removed the standalone `lianxiang.dict.yaml` on 2026-10-15, so this tool no longer syncs `rime_mint.ext.dict.yaml`. This oh-my-rime migration replaces ext with an empty compatibility placeholder, so existing `rime_mint.ext` references in users' `import_tables` continue to load. New configurations can omit this reference; delete the placeholder only after all references have been removed. The tool does not automatically clean up retired dictionaries in the target repository. Upstream filtered and reorganized the content, so the base dictionary is not guaranteed to contain every entry from the old ext dictionary.

## Example Output

After running the program:

```txt
正在下载: rime_mint.base.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.base.dict.yaml
正在下载: rime_mint.chars.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.chars.dict.yaml
正在下载: rime_mint.correlation.dict.yaml ...
已下载并修改: dl_dicts/rime_mint.correlation.dict.yaml
```

Downloaded files are saved in the staging directory specified by `DOWNLOAD_DIR`, which defaults to `dl_dicts`. Local runs do not clear existing files automatically; clear old output before migrating or switching configurations. Every CI sync entry point clears `dl_dicts` before running, so wildcard copies cannot include retired files or files left by another task. Do not store personal files in this directory.

The program and `scripts/sync-config.sh` return `0` when changes are detected, `1` when there are no changes, and `2` for argument, build, download, or processing errors. CI converts only `1` to the skip status `78`; other failures retain their exit codes and stop the pipeline. Custom sync scripts must follow this convention and write their output to `dl_dicts`.

Special dictionaries that need to keep the target repository header and only sync body content can use `merge_mode: "keep_remote_header"`, for example:

```bash
bash scripts/sync-config.sh configs/rime-ice-english.yaml
```

## LICENSE

This project is licensed under the MIT License. For details, please see the [LICENSE](./LICENSE) file.
