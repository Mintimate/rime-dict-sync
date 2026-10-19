#!/bin/bash

# Exercise the actual CI shell blocks without downloading or publishing dictionaries.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TEST_DIR="$(mktemp -d)"
trap 'rm -rf -- "${TEST_DIR}"' EXIT

mkdir -p "${TEST_DIR}/blocks" "${TEST_DIR}/workspace/scripts" "${TEST_DIR}/bin"
awk -v output="${TEST_DIR}/blocks" '
  /^            - \|$/ { block++; active = 1; next }
  active && /^              / {
    print substr($0, 15) > (output "/" block ".sh")
    next
  }
  { active = 0 }
  END { if (block != 7) exit 1 }
' "${REPO_ROOT}/.cnb.yml"

cat > "${TEST_DIR}/workspace/scripts/sync-config.sh" <<'EOF'
#!/bin/bash
set -euo pipefail
if [ -e dl_dicts ]; then
  echo "CI did not clear the staging directory" >&2
  exit 99
fi
mkdir dl_dicts
printf 'fresh dictionary\n' > dl_dicts/current.dict.yaml
exit "${SYNC_TEST_STATUS}"
EOF
cp "${TEST_DIR}/workspace/scripts/sync-config.sh" "${TEST_DIR}/workspace/scripts/custom-sync.sh"

assert_status() {
  local expected="$1"
  shift
  local actual=0
  "$@" > "${TEST_DIR}/output.log" 2>&1 || actual=$?
  if [ "${actual}" -ne "${expected}" ]; then
    cat "${TEST_DIR}/output.log" >&2
    echo "Expected exit ${expected}, got ${actual}: $*" >&2
    exit 1
  fi
}

cd "${TEST_DIR}/workspace"
for block in "${TEST_DIR}"/blocks/*.sh; do
  for status in 0 1 2 7; do
    mkdir -p dl_dicts
    printf 'retired dictionary\n' > dl_dicts/rime_mint.ext.dict.yaml
    expected="${status}"
    if [ "${status}" -eq 1 ]; then
      expected=78
    fi
    assert_status "${expected}" env SYNC_TEST_STATUS="${status}" \
      sync_shell=scripts/custom-sync.sh bash -e "${block}"
    test -f dl_dicts/current.dict.yaml
    test ! -e dl_dicts/rime_mint.ext.dict.yaml
  done
done

# The config wrapper must preserve sync statuses and distinguish build failure
# (go build commonly returns 1) from a successful no-change result.
cat > "${TEST_DIR}/sync-stub" <<'EOF'
#!/bin/bash
exit "${SYNC_TEST_STATUS}"
EOF
cat > "${TEST_DIR}/bin/go" <<'EOF'
#!/bin/bash
exit 1
EOF
chmod +x "${TEST_DIR}/sync-stub" "${TEST_DIR}/bin/go"

for status in 0 1 2 7; do
  assert_status "${status}" env SYNC_TEST_STATUS="${status}" \
    GO_MAIN="${TEST_DIR}/sync-stub" bash "${REPO_ROOT}/scripts/sync-config.sh" config.yaml
done
assert_status 2 env GO_MAIN=./rime-dict-sync PATH="${TEST_DIR}/bin:${PATH}" \
  bash "${REPO_ROOT}/scripts/sync-config.sh" config.yaml
assert_status 2 env GO_MAIN=./missing-sync \
  bash "${REPO_ROOT}/scripts/sync-config.sh" config.yaml
assert_status 2 bash "${REPO_ROOT}/scripts/sync-config.sh"

echo "PASS: 7 CI entry points clear stale output and handle sync statuses; wrapper handles build failures."
