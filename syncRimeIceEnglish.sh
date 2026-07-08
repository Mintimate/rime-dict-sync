#!/bin/bash

set -euo pipefail

DOWNLOAD_DIR="dl_dicts"
REMOTE_REPO="https://github.com/Mintimate/oh-my-rime.git"
REMOTE_DIR="temp_oh_my_rime"
RAW_BASE="https://raw.githubusercontent.com/iDvel/rime-ice/main/en_dicts"
TMP_DIR="$(mktemp -d)"

trap 'rm -rf "${TMP_DIR}" "${REMOTE_DIR}"' EXIT

rm -rf "${DOWNLOAD_DIR}" "${REMOTE_DIR}"
mkdir -p "${DOWNLOAD_DIR}"

git clone --depth 1 "${REMOTE_REPO}" "${REMOTE_DIR}"

download() {
  local source_name="$1"
  local output_path="$2"

  curl -fsSL "${RAW_BASE}/${source_name}" -o "${output_path}"
}

merge_yaml_dict() {
  local remote_path="$1"
  local upstream_path="$2"
  local output_path="$3"
  local dict_name="$4"

  awk '1; /^---[[:space:]]*$/{exit}' "${remote_path}" > "${output_path}"
  awk -v dict_name="${dict_name}" '
    found && /^\.\.\.[[:space:]]*$/ { print; exit }
    found {
      if ($0 ~ /^name:[[:space:]]*/) {
        print "name: " dict_name
      } else {
        print
      }
    }
    /^---[[:space:]]*$/ { found = 1 }
  ' "${upstream_path}" >> "${output_path}"
  awk 'seen { print } /^\.\.\.[[:space:]]*$/ { seen = 1 }' "${upstream_path}" >> "${output_path}"
}

merge_table_dict() {
  local remote_path="$1"
  local upstream_path="$2"
  local output_path="$3"

  awk '1; /此行之后不能写注释/{exit}' "${remote_path}" > "${output_path}"
  awk 'seen { print } /此行之后不能写注释/ { seen = 1 }' "${upstream_path}" >> "${output_path}"
}

download "en.dict.yaml" "${TMP_DIR}/en.dict.yaml"
download "en_ext.dict.yaml" "${TMP_DIR}/en_ext.dict.yaml"
download "cn_en.txt" "${TMP_DIR}/cn_en.txt"
download "cn_en_flypy.txt" "${TMP_DIR}/cn_en_flypy.txt"

merge_yaml_dict "${REMOTE_DIR}/dicts/rime_ice.en.dict.yaml" "${TMP_DIR}/en.dict.yaml" "${DOWNLOAD_DIR}/rime_ice.en.dict.yaml" "rime_ice.en"
merge_yaml_dict "${REMOTE_DIR}/dicts/rime_ice.en_ext.dict.yaml" "${TMP_DIR}/en_ext.dict.yaml" "${DOWNLOAD_DIR}/rime_ice.en_ext.dict.yaml" "rime_ice.en_ext"
merge_table_dict "${REMOTE_DIR}/dicts/rime_ice.cn_en.txt" "${TMP_DIR}/cn_en.txt" "${DOWNLOAD_DIR}/rime_ice.cn_en.txt"
merge_table_dict "${REMOTE_DIR}/dicts/rime_ice.cn_en_flypy.txt" "${TMP_DIR}/cn_en_flypy.txt" "${DOWNLOAD_DIR}/rime_ice.cn_en_flypy.txt"

has_changes=0
for dict_file in "${DOWNLOAD_DIR}"/*; do
  filename="$(basename "${dict_file}")"
  if cmp -s "${dict_file}" "${REMOTE_DIR}/dicts/${filename}"; then
    echo "无变化: ${filename}"
  else
    echo "检测到变化: ${filename}"
    has_changes=1
  fi
done

if [ "${has_changes}" -eq 0 ]; then
  echo "词库无变化，跳过更新"
  exit 1
fi

echo "词库同步完成，有变化需要更新"
