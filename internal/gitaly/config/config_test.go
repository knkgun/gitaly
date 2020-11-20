package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/gitaly/config/sentry"
)

func TestLoad_doesntClearPreviousGlobalConfig(t *testing.T) {
	defer func(old Cfg) { Config = old }(Config)

	Config = Cfg{SocketPath: "/tmp"}
	cfg, err := Load(&bytes.Buffer{})
	require.NoError(t, err)

	require.Equal(t, "", cfg.SocketPath)
	require.Equal(t, "/tmp", Config.SocketPath)
}

func TestLoadBrokenConfig(t *testing.T) {
	tmpFile := strings.NewReader(`path = "/tmp"\nname="foo"`)
	_, err := Load(tmpFile)
	assert.Error(t, err)
}

func TestLoadEmptyConfig(t *testing.T) {
	cfg, err := Load(strings.NewReader(``))
	require.NoError(t, err)

	defaultConf := Cfg{InternalSocketDir: cfg.InternalSocketDir}
	require.NoError(t, defaultConf.setDefaults())

	assert.Equal(t, defaultConf, cfg)
}

func TestLoadURLs(t *testing.T) {
	tmpFile := strings.NewReader(`
[gitlab]
url = "unix:///tmp/test.socket"
relative_url_root = "/gitlab"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	defaultConf := Cfg{
		Gitlab: Gitlab{
			URL:             "unix:///tmp/test.socket",
			RelativeURLRoot: "/gitlab",
		},
	}
	require.NoError(t, defaultConf.setDefaults())

	assert.Equal(t, defaultConf.Gitlab, cfg.Gitlab)
}

func TestLoadStorage(t *testing.T) {
	tmpFile := strings.NewReader(`[[storage]]
name = "default"
path = "/tmp/"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	if assert.Equal(t, 1, len(cfg.Storages), "Expected one (1) storage") {
		expectedConf := Cfg{
			Storages: []Storage{
				{Name: "default", Path: "/tmp"},
			},
		}
		require.NoError(t, expectedConf.setDefaults())

		assert.Equal(t, expectedConf.Storages, cfg.Storages)
	}
}

func TestUncleanStoragePaths(t *testing.T) {
	cfg, err := Load(strings.NewReader(`[[storage]]
name="unclean-path-1"
path="/tmp/repos1//"

[[storage]]
name="unclean-path-2"
path="/tmp/repos2/subfolder/.."
`))
	require.NoError(t, err)

	require.Equal(t, []Storage{
		{Name: "unclean-path-1", Path: "/tmp/repos1"},
		{Name: "unclean-path-2", Path: "/tmp/repos2"},
	}, cfg.Storages)
}

func TestLoadMultiStorage(t *testing.T) {
	tmpFile := strings.NewReader(`[[storage]]
name="default"
path="/tmp/repos1"

[[storage]]
name="other"
path="/tmp/repos2/"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	if assert.Equal(t, 2, len(cfg.Storages), "Expected one (1) storage") {
		expectedConf := Cfg{
			Storages: []Storage{
				{Name: "default", Path: "/tmp/repos1"},
				{Name: "other", Path: "/tmp/repos2"},
			},
		}
		require.NoError(t, expectedConf.setDefaults())

		assert.Equal(t, expectedConf.Storages, cfg.Storages)
	}
}

func TestLoadSentry(t *testing.T) {
	tmpFile := strings.NewReader(`[logging]
sentry_environment = "production"
sentry_dsn = "abc123"
ruby_sentry_dsn = "xyz456"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	expectedConf := Cfg{
		Logging: Logging{
			Sentry: Sentry(sentry.Config{
				Environment: "production",
				DSN:         "abc123",
			}),
			RubySentryDSN: "xyz456",
		},
	}
	require.NoError(t, expectedConf.setDefaults())

	assert.Equal(t, expectedConf.Logging, cfg.Logging)
}

func TestLoadPrometheus(t *testing.T) {
	tmpFile := strings.NewReader(`prometheus_listen_addr=":9236"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, ":9236", cfg.PrometheusListenAddr)
}

func TestLoadSocketPath(t *testing.T) {
	tmpFile := strings.NewReader(`socket_path="/tmp/gitaly.sock"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, "/tmp/gitaly.sock", cfg.SocketPath)
}

func TestLoadListenAddr(t *testing.T) {
	tmpFile := strings.NewReader(`listen_addr=":8080"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, ":8080", cfg.ListenAddr)
}

func tempEnv(key, value string) func() {
	temp := os.Getenv(key)
	os.Setenv(key, value)

	return func() {
		os.Setenv(key, temp)
	}
}

func TestLoadOverrideEnvironment(t *testing.T) {
	// Test that this works since we still want this to work
	tempEnv1 := tempEnv("GITALY_SOCKET_PATH", "/tmp/gitaly2.sock")
	defer tempEnv1()
	tempEnv2 := tempEnv("GITALY_LISTEN_ADDR", ":8081")
	defer tempEnv2()
	tempEnv3 := tempEnv("GITALY_PROMETHEUS_LISTEN_ADDR", ":9237")
	defer tempEnv3()

	tmpFile := strings.NewReader(`socket_path = "/tmp/gitaly.sock"
listen_addr = ":8080"
prometheus_listen_addr = ":9236"`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	assert.Equal(t, ":9237", cfg.PrometheusListenAddr)
	assert.Equal(t, "/tmp/gitaly2.sock", cfg.SocketPath)
	assert.Equal(t, ":8081", cfg.ListenAddr)
}

func TestLoadOnlyEnvironment(t *testing.T) {
	// Test that this works since we still want this to work
	defer tempEnv("GITALY_SOCKET_PATH", "/tmp/gitaly2.sock")()
	defer tempEnv("GITALY_LISTEN_ADDR", ":8081")()
	defer tempEnv("GITALY_PROMETHEUS_LISTEN_ADDR", ":9237")()

	cfg, err := Load(&bytes.Buffer{})
	require.NoError(t, err)

	assert.Equal(t, ":9237", cfg.PrometheusListenAddr)
	assert.Equal(t, "/tmp/gitaly2.sock", cfg.SocketPath)
	assert.Equal(t, ":8081", cfg.ListenAddr)
}

func TestValidateStorages(t *testing.T) {
	repositories, err := filepath.Abs("testdata/repositories")
	require.NoError(t, err)

	repositories2, err := filepath.Abs("testdata/repositories2")
	require.NoError(t, err)

	invalidDir := filepath.Join(repositories, t.Name())

	testCases := []struct {
		desc     string
		storages []Storage
		invalid  bool
	}{
		{
			desc: "just 1 storage",
			storages: []Storage{
				{Name: "default", Path: repositories},
			},
		},
		{
			desc: "multiple storages",
			storages: []Storage{
				{Name: "default", Path: repositories},
				{Name: "other", Path: repositories2},
			},
		},
		{
			desc: "multiple storages pointing to same directory",
			storages: []Storage{
				{Name: "default", Path: repositories},
				{Name: "other", Path: repositories},
				{Name: "third", Path: repositories},
			},
		},
		{
			desc: "nested paths 1",
			storages: []Storage{
				{Name: "default", Path: "/home/git/repositories"},
				{Name: "other", Path: "/home/git/repositories"},
				{Name: "third", Path: "/home/git/repositories/third"},
			},
			invalid: true,
		},
		{
			desc: "nested paths 2",
			storages: []Storage{
				{Name: "default", Path: "/home/git/repositories/default"},
				{Name: "other", Path: "/home/git/repositories"},
				{Name: "third", Path: "/home/git/repositories"},
			},
			invalid: true,
		},
		{
			desc: "duplicate definition",
			storages: []Storage{
				{Name: "default", Path: repositories},
				{Name: "default", Path: repositories},
			},
			invalid: true,
		},
		{
			desc: "re-definition",
			storages: []Storage{
				{Name: "default", Path: repositories},
				{Name: "default", Path: repositories2},
			},
			invalid: true,
		},
		{
			desc: "empty name",
			storages: []Storage{
				{Name: "", Path: repositories},
			},
			invalid: true,
		},
		{
			desc: "empty path",
			storages: []Storage{
				{Name: "default", Path: ""},
			},
			invalid: true,
		},
		{
			desc: "non existing directory",
			storages: []Storage{
				{Name: "default", Path: repositories},
				{Name: "nope", Path: invalidDir},
			},
			invalid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := Cfg{Storages: tc.storages}

			err := cfg.validateStorages()
			if tc.invalid {
				assert.Error(t, err, "%+v", tc.storages)
				return
			}

			assert.NoError(t, err, "%+v", tc.storages)
		})
	}
}

func TestStoragePath(t *testing.T) {
	cfg := Cfg{Storages: []Storage{
		{Name: "default", Path: "/home/git/repositories1"},
		{Name: "other", Path: "/home/git/repositories2"},
		{Name: "third", Path: "/home/git/repositories3"},
	}}

	testCases := []struct {
		in, out string
		ok      bool
	}{
		{in: "default", out: "/home/git/repositories1", ok: true},
		{in: "third", out: "/home/git/repositories3", ok: true},
		{in: "", ok: false},
		{in: "foobar", ok: false},
	}

	for _, tc := range testCases {
		out, ok := cfg.StoragePath(tc.in)
		if !assert.Equal(t, tc.ok, ok, "%+v", tc) {
			continue
		}
		assert.Equal(t, tc.out, out, "%+v", tc)
	}
}

type hookFileMode int

const (
	hookFileExists hookFileMode = 1 << (4 - 1 - iota)
	hookFileExecutable
)

func setupTempHookDirs(t *testing.T, m map[string]hookFileMode) (string, func()) {
	tempDir, err := ioutil.TempDir("", "hooks")
	require.NoError(t, err)

	for hookName, mode := range m {
		if mode&hookFileExists > 0 {
			path := filepath.Join(tempDir, hookName)
			require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))

			require.NoError(t, ioutil.WriteFile(filepath.Join(tempDir, hookName), nil, 0644))

			if mode&hookFileExecutable > 0 {
				require.NoError(t, os.Chmod(filepath.Join(tempDir, hookName), 0755))
			}
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}

var (
	fileNotExistsErrRegexSnippit  = "no such file or directory"
	fileNotExecutableRegexSnippit = "not executable: .*"
)

func TestValidateHooks(t *testing.T) {
	testCases := []struct {
		desc             string
		expectedErrRegex string
		hookFiles        map[string]hookFileMode
	}{
		{
			desc: "everything is ✅",
			hookFiles: map[string]hookFileMode{
				"ruby/git-hooks/update":       hookFileExists | hookFileExecutable,
				"ruby/git-hooks/pre-receive":  hookFileExists | hookFileExecutable,
				"ruby/git-hooks/post-receive": hookFileExists | hookFileExecutable,
			},
			expectedErrRegex: "",
		},
		{
			desc: "missing git-hooks",
			hookFiles: map[string]hookFileMode{
				"ruby/git-hooks/update":       0,
				"ruby/git-hooks/pre-receive":  0,
				"ruby/git-hooks/post-receive": 0,
			},
			expectedErrRegex: fmt.Sprintf("%s, %s, %s", fileNotExistsErrRegexSnippit, fileNotExistsErrRegexSnippit, fileNotExistsErrRegexSnippit),
		},
		{
			desc: "git-hooks are not executable",
			hookFiles: map[string]hookFileMode{
				"ruby/git-hooks/update":       hookFileExists,
				"ruby/git-hooks/pre-receive":  hookFileExists,
				"ruby/git-hooks/post-receive": hookFileExists,
			},
			expectedErrRegex: fmt.Sprintf("%s, %s, %s", fileNotExecutableRegexSnippit, fileNotExecutableRegexSnippit, fileNotExecutableRegexSnippit),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			originalConfig := Config
			defer func() {
				Config = originalConfig
			}()

			tempHookDir, cleanup := setupTempHookDirs(t, tc.hookFiles)
			defer cleanup()

			cfg := Cfg{
				Ruby: Ruby{
					Dir: filepath.Join(tempHookDir, "ruby"),
				},
				GitlabShell: GitlabShell{
					Dir: filepath.Join(tempHookDir, "/gitlab-shell"),
				},
				BinDir: filepath.Join(tempHookDir, "/bin"),
			}

			err := cfg.validateHooks()
			if tc.expectedErrRegex != "" {
				require.Error(t, err)
				require.Regexp(t, tc.expectedErrRegex, err.Error(), "error should match regexp")
			}
		})
	}
}

func TestLoadGit(t *testing.T) {
	tmpFile := strings.NewReader(`[git]
bin_path = "/my/git/path"
catfile_cache_size = 50
`)

	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	require.Equal(t, "/my/git/path", cfg.Git.BinPath)
	require.Equal(t, 50, cfg.Git.CatfileCacheSize)
}

func TestSetGitPath(t *testing.T) {
	var resolvedGitPath string
	if path, ok := os.LookupEnv("GITALY_TESTING_GIT_BINARY"); ok {
		resolvedGitPath = path
	} else {
		path, err := exec.LookPath("git")
		require.NoError(t, err)
		resolvedGitPath = path
	}

	testCases := []struct {
		desc       string
		gitBinPath string
		expected   string
	}{
		{
			desc:       "With a Git Path set through the settings",
			gitBinPath: "/path/to/myGit",
			expected:   "/path/to/myGit",
		},
		{
			desc:       "When a git path hasn't been set",
			gitBinPath: "",
			expected:   resolvedGitPath,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := Cfg{Git: Git{BinPath: tc.gitBinPath}}
			require.NoError(t, cfg.SetGitPath())
			assert.Equal(t, tc.expected, cfg.Git.BinPath, tc.desc)
		})
	}
}

func TestValidateShellPath(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "gitaly-tests-")
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Join(tmpDir, "bin"), 0755))
	tmpFile := filepath.Join(tmpDir, "my-file")
	defer os.RemoveAll(tmpDir)
	fp, err := os.Create(tmpFile)
	require.NoError(t, err)
	require.NoError(t, fp.Close())

	testCases := []struct {
		desc      string
		path      string
		shouldErr bool
	}{
		{
			desc:      "When no Shell Path set",
			path:      "",
			shouldErr: true,
		},
		{
			desc:      "When Shell Path set to non-existing path",
			path:      "/non/existing/path",
			shouldErr: true,
		},
		{
			desc:      "When Shell Path set to non-dir path",
			path:      tmpFile,
			shouldErr: true,
		},
		{
			desc:      "When Shell Path set to a valid directory",
			path:      tmpDir,
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := Cfg{GitlabShell: GitlabShell{Dir: tc.path}}
			err := cfg.validateShell()
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigureRuby(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "gitaly-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "file")
	require.NoError(t, ioutil.WriteFile(tmpFile, nil, 0644))

	testCases := []struct {
		dir  string
		ok   bool
		desc string
	}{
		{dir: "", desc: "empty"},
		{dir: "/does/not/exist", desc: "does not exist"},
		{dir: tmpFile, desc: "exists but is not a directory"},
		{dir: ".", ok: true, desc: "relative path"},
		{dir: tmpDir, ok: true, desc: "ok"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := Cfg{Ruby: Ruby{Dir: tc.dir}}

			err := cfg.ConfigureRuby()
			if !tc.ok {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			dir := cfg.Ruby.Dir
			require.True(t, filepath.IsAbs(dir), "expected %q to be absolute path", dir)
		})
	}
}

func TestConfigureRubyNumWorkers(t *testing.T) {
	testCases := []struct {
		in, out int
	}{
		{in: -1, out: 2},
		{in: 0, out: 2},
		{in: 1, out: 2},
		{in: 2, out: 2},
		{in: 3, out: 3},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			cfg := Cfg{Ruby: Ruby{Dir: "/", NumWorkers: tc.in}}
			require.NoError(t, cfg.ConfigureRuby())
			require.Equal(t, tc.out, cfg.Ruby.NumWorkers)
		})
	}
}

func TestValidateListeners(t *testing.T) {
	testCases := []struct {
		desc string
		Cfg
		ok bool
	}{
		{desc: "empty"},
		{desc: "socket only", Cfg: Cfg{SocketPath: "/foo/bar"}, ok: true},
		{desc: "tcp only", Cfg: Cfg{ListenAddr: "a.b.c.d:1234"}, ok: true},
		{desc: "both socket and tcp", Cfg: Cfg{SocketPath: "/foo/bar", ListenAddr: "a.b.c.d:1234"}, ok: true},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.Cfg.validateListeners()
			if tc.ok {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestLoadGracefulRestartTimeout(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		expected time.Duration
	}{
		{
			name:     "default value",
			expected: 1 * time.Minute,
		},
		{
			name:     "8m03s",
			config:   `graceful_restart_timeout = "8m03s"`,
			expected: 8*time.Minute + 3*time.Second,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpFile := strings.NewReader(test.config)

			cfg, err := Load(tmpFile)
			assert.NoError(t, err)

			assert.Equal(t, test.expected, cfg.GracefulRestartTimeout.Duration())
		})
	}
}

func TestGitlabShellDefaults(t *testing.T) {
	gitlabShellDir := "/dir"
	expectedGitlab := Gitlab{
		SecretFile: filepath.Join(gitlabShellDir, ".gitlab_shell_secret"),
	}

	expectedHooks := Hooks{
		CustomHooksDir: filepath.Join(gitlabShellDir, "hooks"),
	}

	tmpFile := strings.NewReader(fmt.Sprintf(`[gitlab-shell]
dir = '%s'`, gitlabShellDir))
	cfg, err := Load(tmpFile)
	require.NoError(t, err)

	require.Equal(t, expectedGitlab, cfg.Gitlab)
	require.Equal(t, expectedHooks, cfg.Hooks)
}

func TestValidateInternalSocketDir(t *testing.T) {
	// create a valid socket directory
	tempDir, err := ioutil.TempDir("testdata", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// create a symlinked socket directory
	dirName := "internal_socket_dir"
	validSocketDirSymlink := filepath.Join(tempDir, dirName)
	tmpSocketDir, err := ioutil.TempDir(tempDir, "")
	require.NoError(t, err)
	tmpSocketDir, err = filepath.Abs(tmpSocketDir)
	require.NoError(t, err)
	require.NoError(t, os.Symlink(tmpSocketDir, validSocketDirSymlink))

	// create a broken symlink
	dirName = "internal_socket_dir_broken"
	brokenSocketDirSymlink := filepath.Join(tempDir, dirName)
	require.NoError(t, os.Symlink("/does/not/exist", brokenSocketDirSymlink))

	testCases := []struct {
		desc              string
		internalSocketDir string
		shouldError       bool
	}{
		{
			desc:              "empty socket dir",
			internalSocketDir: "",
			shouldError:       false,
		},
		{
			desc:              "non existing directory",
			internalSocketDir: "/tmp/relative/path/to/nowhere",
			shouldError:       true,
		},
		{
			desc:              "valid socket directory",
			internalSocketDir: tempDir,
			shouldError:       false,
		},
		{
			desc:              "valid symlinked directory",
			internalSocketDir: validSocketDirSymlink,
			shouldError:       false,
		},
		{
			desc:              "broken symlinked directory",
			internalSocketDir: brokenSocketDirSymlink,
			shouldError:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := Cfg{InternalSocketDir: tc.internalSocketDir}
			if tc.shouldError {
				assert.Error(t, cfg.validateInternalSocketDir())
				return
			}
			assert.NoError(t, cfg.validateInternalSocketDir())
		})
	}
}

func TestInternalSocketDir(t *testing.T) {
	cfg, err := Load(bytes.NewReader(nil))
	require.NoError(t, err)
	socketDir := cfg.InternalSocketDir

	require.NoError(t, trySocketCreation(socketDir))
	require.NoError(t, os.RemoveAll(socketDir))
}

func TestLoadDailyMaintenance(t *testing.T) {
	for _, tt := range []struct {
		name        string
		rawCfg      string
		expect      DailyJob
		loadErr     error
		validateErr error
	}{
		{
			name: "success",
			rawCfg: `[[storage]]
name = "default"
path = "/"

			[daily_maintenance]
start_hour = 11
start_minute = 23
duration = "45m"
storages = ["default"]
`,
			expect: DailyJob{
				Hour:     11,
				Minute:   23,
				Duration: Duration(45 * time.Minute),
				Storages: []string{"default"},
			},
		},
		{
			rawCfg: `[daily_maintenance]
			start_hour = 24`,
			expect: DailyJob{
				Hour: 24,
			},
			validateErr: errors.New("daily maintenance specified hour '24' outside range (0-23)"),
		}, {
			rawCfg: `[daily_maintenance]
			start_hour = 60`,
			expect: DailyJob{
				Hour: 60,
			},
			validateErr: errors.New("daily maintenance specified hour '60' outside range (0-23)"),
		},
		{
			rawCfg: `[daily_maintenance]
			duration = "meow"`,
			expect:  DailyJob{},
			loadErr: errors.New("load toml: (2, 4): unmarshal text: time: invalid duration meow"),
		}, {
			rawCfg: `[daily_maintenance]
			storages = ["default"]`,
			expect: DailyJob{
				Storages: []string{"default"},
			},
			validateErr: errors.New(`daily maintenance specified storage "default" does not exist in configuration`),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := strings.NewReader(tt.rawCfg)
			cfg, err := Load(tmpFile)
			require.Equal(t, tt.loadErr, err)
			require.Equal(t, tt.expect, cfg.DailyMaintenance)
			require.Equal(t, tt.validateErr, cfg.validateMaintenance())
		})
	}
}
