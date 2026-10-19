package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Exercise the real CLI exit codes in a child process so os.Exit is observable.
func TestSyncCLIHelper(t *testing.T) {
	configPath := os.Getenv("RIME_SYNC_TEST_CONFIG")
	if configPath == "" {
		return
	}
	flag.CommandLine = flag.NewFlagSet("rime-dict-sync", flag.ExitOnError)
	os.Args = []string{"rime-dict-sync", "-c", configPath}
	main()
	os.Exit(0)
}

func TestSyncExitCodes(t *testing.T) {
	const upstream = "# Rime dictionary\n---\nname: jichu\nversion: LTS\nsort: by_weight\n...\n测试\tcè shì\t100\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(upstream))
	}))
	defer server.Close()

	for _, tc := range []struct {
		name       string
		remoteBody string
		missing    bool
		wantCode   int
	}{
		{"changed", "旧词\tjiù cí\t10\n", false, 0},
		{"unchanged", "测试\tcè shì\t100\n", false, 1},
		{"removed_upstream", "旧词\tjiù cí\t10\n", true, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			repo := filepath.Join(dir, "remote")
			if err := os.MkdirAll(filepath.Join(repo, "dicts"), 0755); err != nil {
				t.Fatal(err)
			}
			const name = "rime_mint.base.dict.yaml"
			remote := "---\nname: rime_mint.base\nversion: old\nsort: by_weight\n...\n" + tc.remoteBody
			if err := os.WriteFile(filepath.Join(repo, "dicts", name), []byte(remote), 0644); err != nil {
				t.Fatal(err)
			}
			for _, args := range [][]string{
				{"init", "--quiet"},
				{"add", "."},
				{"-c", "user.name=Sync Test", "-c", "user.email=sync-test@example.invalid", "-c", "commit.gpgSign=false", "commit", "--quiet", "-m", "fixture"},
			} {
				cmd := exec.Command("git", args...)
				cmd.Dir = repo
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git %v: %v\n%s", args, err, output)
				}
			}
			url := server.URL + "/jichu"
			if tc.missing {
				url = server.URL + "/missing"
			}
			config := DictConfig{
				DOWNLOAD_DIR: filepath.Join(dir, "dl_dicts"),
				TARGET_DICT:  []DictTarget{{Name: name, URL: url, RemoteRepo: repo, RemotePath: "dicts/" + name}},
			}
			data, err := yaml.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			configPath := filepath.Join(dir, "config.yaml")
			if err := os.WriteFile(configPath, data, 0644); err != nil {
				t.Fatal(err)
			}
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(executable, "-test.run=^TestSyncCLIHelper$")
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "RIME_SYNC_TEST_CONFIG="+configPath)
			output, err := cmd.CombinedOutput()
			code := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					code = exitErr.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			if code != tc.wantCode {
				t.Fatalf("exit = %d, want %d\n%s", code, tc.wantCode, output)
			}
			if _, err := os.Stat(filepath.Join(dir, "temp_remote_repo_"+generateRepoId(repo))); !os.IsNotExist(err) {
				t.Fatalf("temporary clone was not cleaned up: %v", err)
			}
			content, err := os.ReadFile(filepath.Join(config.DOWNLOAD_DIR, name))
			if tc.missing {
				if !os.IsNotExist(err) {
					t.Fatalf("404 must not produce a dictionary: %v", err)
				}
			} else if err != nil || !strings.Contains(string(content), "name: rime_mint.base\n") || !strings.HasSuffix(string(content), "测试\tcè shì\t100\n") {
				t.Fatalf("unexpected converted dictionary: %s, %v", content, err)
			}
		})
	}
}

func TestRejectInvalidConfig(t *testing.T) {
	for _, content := range []string{
		"",
		"TARGET_DICT: []\n",
		"TARGET_DICT:\n  - name: example.dict.yaml\n",
		"TARGET_DICT:\n  - url: https://example.invalid/dict\n",
		"TARGET_DICT:\n  - name: ../example.dict.yaml\n    url: https://example.invalid/dict\n",
		"TARGET_DICT:\n  - name: example.dict.yaml\n    url: https://example.invalid/one\n  - name: example.dict.yaml\n    url: https://example.invalid/two\n",
	} {
		t.Run(content, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
			if _, err := loadConfig(path); err == nil {
				t.Fatal("invalid config was accepted")
			}
		})
	}
}
