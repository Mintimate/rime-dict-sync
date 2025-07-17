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
3. Create configuration file `config.yaml` (refer to Configuration section below)
4. Run program:
   ```bash
   go mod tidy
   go run .
   ```

## Configuration

Create `config.yaml` file with following format:

```yaml
# Dictionaries to download: save_name-download_url
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

## Example Output

After running the program:

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

Downloaded dictionary files are saved in `dl_dicts` directory and ready for use with Rime input method.

## LICENSE

This project is licensed under the MIT License. For details, please see the [LICENSE](./LICENSE) file.