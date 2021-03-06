package config

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/auth"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/cgroups"
	internallog "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/log"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/prometheus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/sentry"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/env"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"golang.org/x/sys/unix"
)

const (
	// GitalyDataPrefix is the top-level directory we use to store system
	// (non-user) data. We need to be careful that this path does not clash
	// with any directory name that could be provided by a user. The '+'
	// character is not allowed in GitLab namespaces or repositories.
	GitalyDataPrefix = "+gitaly"
)

// DailyJob enables a daily task to be scheduled for specific storages
type DailyJob struct {
	Hour     uint     `toml:"start_hour"`
	Minute   uint     `toml:"start_minute"`
	Duration Duration `toml:"duration"`
	Storages []string `toml:"storages"`

	// Disabled will completely disable a daily job, even in cases where a
	// default schedule is implied
	Disabled bool `toml:"disabled"`
}

// Cfg is a container for all config derived from config.toml.
type Cfg struct {
	SocketPath             string            `toml:"socket_path" split_words:"true"`
	ListenAddr             string            `toml:"listen_addr" split_words:"true"`
	TLSListenAddr          string            `toml:"tls_listen_addr" split_words:"true"`
	PrometheusListenAddr   string            `toml:"prometheus_listen_addr" split_words:"true"`
	BinDir                 string            `toml:"bin_dir"`
	Git                    Git               `toml:"git" envconfig:"git"`
	Storages               []Storage         `toml:"storage" envconfig:"storage"`
	Logging                Logging           `toml:"logging" envconfig:"logging"`
	Prometheus             prometheus.Config `toml:"prometheus"`
	Auth                   auth.Config       `toml:"auth"`
	TLS                    TLS               `toml:"tls"`
	Ruby                   Ruby              `toml:"gitaly-ruby"`
	Gitlab                 Gitlab            `toml:"gitlab"`
	GitlabShell            GitlabShell       `toml:"gitlab-shell"`
	Hooks                  Hooks             `toml:"hooks"`
	Concurrency            []Concurrency     `toml:"concurrency"`
	GracefulRestartTimeout Duration          `toml:"graceful_restart_timeout"`
	InternalSocketDir      string            `toml:"internal_socket_dir"`
	DailyMaintenance       DailyJob          `toml:"daily_maintenance"`
	Cgroups                cgroups.Config    `toml:"cgroups"`
	PackObjectsCache       StreamCacheConfig `toml:"pack_objects_cache"`
}

// TLS configuration
type TLS struct {
	CertPath string `toml:"certificate_path,omitempty"`
	KeyPath  string `toml:"key_path,omitempty"`
}

// GitlabShell contains the settings required for executing `gitlab-shell`
type GitlabShell struct {
	Dir string `toml:"dir" json:"dir"`
}

// Gitlab contains settings required to connect to the Gitlab api
type Gitlab struct {
	URL             string       `toml:"url" json:"url"`
	RelativeURLRoot string       `toml:"relative_url_root" json:"relative_url_root"` // For UNIX sockets only
	HTTPSettings    HTTPSettings `toml:"http-settings" json:"http_settings"`
	SecretFile      string       `toml:"secret_file" json:"secret_file"`
}

// Hooks contains the settings required for hooks
type Hooks struct {
	CustomHooksDir string `toml:"custom_hooks_dir" json:"custom_hooks_dir"`
}

//nolint: revive,stylecheck // This is unintentionally missing documentation.
type HTTPSettings struct {
	ReadTimeout int    `toml:"read_timeout" json:"read_timeout"`
	User        string `toml:"user" json:"user"`
	Password    string `toml:"password" json:"password"`
	CAFile      string `toml:"ca_file" json:"ca_file"`
	CAPath      string `toml:"ca_path" json:"ca_path"`
	SelfSigned  bool   `toml:"self_signed_cert" json:"self_signed_cert"`
}

// Git contains the settings for the Git executable
type Git struct {
	UseBundledBinaries bool        `toml:"use_bundled_binaries"`
	BinPath            string      `toml:"bin_path"`
	CatfileCacheSize   int         `toml:"catfile_cache_size"`
	Config             []GitConfig `toml:"config"`
	// HooksPath is the location where Gitaly has its hooks. This variable cannot be set via the
	// config file and is only used in our tests.
	HooksPath string `toml:"-"`

	execEnv []string
}

// GitConfig contains a key-value pair which is to be passed to git as configuration.
type GitConfig struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

// Storage contains a single storage-shard
type Storage struct {
	Name string
	Path string
}

// Sentry is a sentry.Config. We redefine this type to a different name so
// we can embed both structs into Logging
type Sentry sentry.Config

// Logging contains the logging configuration for Gitaly
type Logging struct {
	internallog.Config
	Sentry

	RubySentryDSN string `toml:"ruby_sentry_dsn"`
}

// Concurrency allows endpoints to be limited to a maximum concurrency per repo
type Concurrency struct {
	RPC        string `toml:"rpc"`
	MaxPerRepo int    `toml:"max_per_repo"`
}

// StreamCacheConfig contains settings for a streamcache instance.
type StreamCacheConfig struct {
	Enabled bool     `toml:"enabled"` // Default: false
	Dir     string   `toml:"dir"`     // Default: <FIRST STORAGE PATH>/+gitaly/PackObjectsCache
	MaxAge  Duration `toml:"max_age"` // Default: 5m
}

// Load initializes the Config variable from file and the environment.
//  Environment variables take precedence over the file.
func Load(file io.Reader) (Cfg, error) {
	cfg := Cfg{
		Prometheus: prometheus.DefaultConfig(),
	}

	if err := toml.NewDecoder(file).Decode(&cfg); err != nil {
		return Cfg{}, fmt.Errorf("load toml: %v", err)
	}

	if err := cfg.setDefaults(); err != nil {
		return Cfg{}, err
	}

	for i := range cfg.Storages {
		cfg.Storages[i].Path = filepath.Clean(cfg.Storages[i].Path)
	}

	return cfg, nil
}

// Validate checks the current Config for sanity.
func (cfg *Cfg) Validate() error {
	for _, run := range []func() error{
		cfg.validateListeners,
		cfg.validateStorages,
		cfg.validateToken,
		cfg.validateGit,
		cfg.validateShell,
		cfg.ConfigureRuby,
		cfg.validateBinDir,
		cfg.validateInternalSocketDir,
		cfg.validateHooks,
		cfg.validateMaintenance,
		cfg.validateCgroups,
		cfg.configurePackObjectsCache,
	} {
		if err := run(); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Cfg) setDefaults() error {
	if cfg.GracefulRestartTimeout.Duration() == 0 {
		cfg.GracefulRestartTimeout = Duration(time.Minute)
	}

	if cfg.Gitlab.SecretFile == "" {
		cfg.Gitlab.SecretFile = filepath.Join(cfg.GitlabShell.Dir, ".gitlab_shell_secret")
	}

	if cfg.Hooks.CustomHooksDir == "" {
		cfg.Hooks.CustomHooksDir = filepath.Join(cfg.GitlabShell.Dir, "hooks")
	}

	if cfg.InternalSocketDir == "" {
		// The socket path must be short-ish because listen(2) fails on long
		// socket paths. We hope/expect that os.MkdirTemp creates a directory
		// that is not too deep. We need a directory, not a tempfile, because we
		// will later want to set its permissions to 0700

		tmpDir, err := os.MkdirTemp("", "gitaly-internal")
		if err != nil {
			return fmt.Errorf("create internal socket directory: %w", err)
		}
		cfg.InternalSocketDir = tmpDir
	}

	if reflect.DeepEqual(cfg.DailyMaintenance, DailyJob{}) {
		cfg.DailyMaintenance = defaultMaintenanceWindow(cfg.Storages)
	}

	return nil
}

func (cfg *Cfg) validateListeners() error {
	if len(cfg.SocketPath) == 0 && len(cfg.ListenAddr) == 0 && len(cfg.TLSListenAddr) == 0 {
		return fmt.Errorf("at least one of socket_path, listen_addr or tls_listen_addr must be set")
	}
	return nil
}

func (cfg *Cfg) validateShell() error {
	if len(cfg.GitlabShell.Dir) == 0 {
		return fmt.Errorf("gitlab-shell.dir: is not set")
	}

	return validateIsDirectory(cfg.GitlabShell.Dir, "gitlab-shell.dir")
}

func checkExecutable(path string) error {
	if err := unix.Access(path, unix.X_OK); err != nil {
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("not executable: %v", path)
		}
		return err
	}

	return nil
}

type hookErrs struct {
	errors []error
}

func (h *hookErrs) Error() string {
	var errStrings []string
	for _, err := range h.errors {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, ", ")
}

func (h *hookErrs) Add(err error) {
	h.errors = append(h.errors, err)
}

func (cfg *Cfg) validateHooks() error {
	if SkipHooks() {
		return nil
	}

	errs := &hookErrs{}

	for _, hookName := range []string{"pre-receive", "post-receive", "update"} {
		if err := checkExecutable(filepath.Join(cfg.Ruby.Dir, "git-hooks", hookName)); err != nil {
			errs.Add(err)
			continue
		}
	}

	if len(errs.errors) > 0 {
		return errs
	}

	return nil
}

func validateIsDirectory(path, name string) error {
	s, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%s: path doesn't exist: %q", name, path)
		}
		return fmt.Errorf("%s: %w", name, err)
	}
	if !s.IsDir() {
		return fmt.Errorf("%s: not a directory: %q", name, path)
	}

	log.WithField("dir", path).Debugf("%s set", name)

	return nil
}

func (cfg *Cfg) validateStorages() error {
	if len(cfg.Storages) == 0 {
		return fmt.Errorf("no storage configurations found. Are you using the right format? https://gitlab.com/gitlab-org/gitaly/issues/397")
	}

	for i, storage := range cfg.Storages {
		if storage.Name == "" {
			return fmt.Errorf("empty storage name at declaration %d", i+1)
		}

		if storage.Path == "" {
			return fmt.Errorf("empty storage path for storage %q", storage.Name)
		}

		fs, err := os.Stat(storage.Path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("storage path %q for storage %q doesn't exist", storage.Path, storage.Name)
			}
			return fmt.Errorf("storage %q: %w", storage.Name, err)
		}

		if !fs.IsDir() {
			return fmt.Errorf("storage path %q for storage %q is not a dir", storage.Path, storage.Name)
		}

		for _, other := range cfg.Storages[:i] {
			if other.Name == storage.Name {
				return fmt.Errorf("storage %q is defined more than once", storage.Name)
			}

			if storage.Path == other.Path {
				// This is weird but we allow it for legacy gitlab.com reasons.
				continue
			}

			if strings.HasPrefix(storage.Path, other.Path) || strings.HasPrefix(other.Path, storage.Path) {
				// If storages have the same sub directory, that is allowed
				if filepath.Dir(storage.Path) == filepath.Dir(other.Path) {
					continue
				}
				return fmt.Errorf("storage paths may not nest: %q and %q", storage.Name, other.Name)
			}
		}
	}

	return nil
}

//nolint: revive,stylecheck // This is unintentionally missing documentation.
func SkipHooks() bool {
	enabled, _ := env.GetBool("GITALY_TESTING_NO_GIT_HOOKS", false)
	return enabled
}

// HooksPath returns the path where Gitaly's Git hooks are located.
func (cfg *Cfg) HooksPath() string {
	if len(cfg.Git.HooksPath) > 0 {
		return cfg.Git.HooksPath
	}

	if SkipHooks() {
		return "/var/empty"
	}

	return filepath.Join(cfg.Ruby.Dir, "git-hooks")
}

// SetGitPath populates the variable GitPath with the path to the `git`
// executable. It warns if no path was specified in the configuration.
func (cfg *Cfg) SetGitPath() error {
	switch {
	case cfg.Git.BinPath != "":
		// Nothing to do.
	case os.Getenv("GITALY_TESTING_GIT_BINARY") != "":
		cfg.Git.BinPath = os.Getenv("GITALY_TESTING_GIT_BINARY")
	case os.Getenv("GITALY_TESTING_BUNDLED_GIT_PATH") != "":
		if cfg.BinDir == "" {
			return errors.New("cannot use bundled binaries without bin path being set")
		}

		// We need to symlink pre-built Git binaries into Gitaly's binary directory.
		// Normally they would of course already exist there, but in tests we create a new
		// binary directory for each server and thus need to populate it first.
		for _, binary := range []string{"gitaly-git", "gitaly-git-remote-http", "gitaly-git-http-backend"} {
			bundledGitBinary := filepath.Join(os.Getenv("GITALY_TESTING_BUNDLED_GIT_PATH"), binary)
			if _, err := os.Stat(bundledGitBinary); err != nil {
				return fmt.Errorf("statting %q: %w", binary, err)
			}

			if err := os.Symlink(bundledGitBinary, filepath.Join(cfg.BinDir, binary)); err != nil {
				// While Gitaly's Go tests use a temporary binary directory, Ruby
				// rspecs set up the binary directory to point to our build
				// directory. They thus already contain the Git binaries and don't
				// need symlinking.
				if errors.Is(err, os.ErrExist) {
					continue
				}
				return fmt.Errorf("symlinking bundled %q: %w", binary, err)
			}
		}

		cfg.Git.UseBundledBinaries = true

		fallthrough
	case cfg.Git.UseBundledBinaries:
		if cfg.Git.BinPath != "" {
			return errors.New("cannot set Git path and use bundled binaries")
		}

		// In order to support having a single Git binary only as compared to a complete Git
		// installation, we create our own GIT_EXEC_PATH which contains symlinks to the Git
		// binary for executables which Git expects to be present.
		gitExecPath, err := os.MkdirTemp("", "gitaly-git-exec-path-*")
		if err != nil {
			return fmt.Errorf("creating Git exec path: %w", err)
		}

		for executable, target := range map[string]string{
			"git":                "gitaly-git",
			"git-receive-pack":   "gitaly-git",
			"git-upload-pack":    "gitaly-git",
			"git-upload-archive": "gitaly-git",
			"git-http-backend":   "gitaly-git-http-backend",
			"git-remote-http":    "gitaly-git-remote-http",
			"git-remote-https":   "gitaly-git-remote-http",
			"git-remote-ftp":     "gitaly-git-remote-http",
			"git-remote-ftps":    "gitaly-git-remote-http",
		} {
			if err := os.Symlink(
				filepath.Join(cfg.BinDir, target),
				filepath.Join(gitExecPath, executable),
			); err != nil {
				return fmt.Errorf("linking Git executable %q: %w", executable, err)
			}
		}

		cfg.Git.BinPath = filepath.Join(gitExecPath, "git")
		cfg.Git.execEnv = []string{
			"GIT_EXEC_PATH=" + gitExecPath,
		}
	default:
		resolvedPath, err := exec.LookPath("git")
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				return fmt.Errorf(`"git" executable not found, set path to it in the configuration file or add it to the PATH`)
			}
		}

		log.WithFields(log.Fields{
			"resolvedPath": resolvedPath,
		}).Warn("git path not configured. Using default path resolution")

		cfg.Git.BinPath = resolvedPath
	}

	return nil
}

// GitExecEnv returns environment variables required to be set in the environment when Git executes.
func (cfg *Cfg) GitExecEnv() []string {
	return cfg.Git.execEnv
}

// StoragePath looks up the base path for storageName. The second boolean
// return value indicates if anything was found.
func (cfg *Cfg) StoragePath(storageName string) (string, bool) {
	storage, ok := cfg.Storage(storageName)
	return storage.Path, ok
}

// Storage looks up storageName.
func (cfg *Cfg) Storage(storageName string) (Storage, bool) {
	for _, storage := range cfg.Storages {
		if storage.Name == storageName {
			return storage, true
		}
	}
	return Storage{}, false
}

// GitalyInternalSocketPath is the path to the internal gitaly socket
func (cfg *Cfg) GitalyInternalSocketPath() string {
	return filepath.Join(cfg.InternalSocketDir, fmt.Sprintf("internal_%d.sock", os.Getpid()))
}

func (cfg *Cfg) validateBinDir() error {
	if len(cfg.BinDir) == 0 {
		return fmt.Errorf("bin_dir: is not set")
	}

	if err := validateIsDirectory(cfg.BinDir, "bin_dir"); err != nil {
		return err
	}

	var err error
	cfg.BinDir, err = filepath.Abs(cfg.BinDir)
	return err
}

// validateGitConfigKey does a best-effort check whether or not a given git config key is valid. It
// does not allow for assignments in keys, which is overly strict and does not allow some valid
// keys. It does avoid misinterpretation of keys though and should catch many cases of
// misconfiguration.
func validateGitConfigKey(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if strings.Contains(key, "=") {
		return errors.New("key cannot contain assignment")
	}
	if !strings.Contains(key, ".") {
		return errors.New("key must contain at least one section")
	}
	if strings.HasPrefix(key, ".") || strings.HasSuffix(key, ".") {
		return errors.New("key must not start or end with a dot")
	}
	return nil
}

func (cfg *Cfg) validateGit() error {
	if err := cfg.SetGitPath(); err != nil {
		return err
	}

	for _, configPair := range cfg.Git.Config {
		if err := validateGitConfigKey(configPair.Key); err != nil {
			return fmt.Errorf("invalid configuration key %q: %w", configPair.Key, err)
		}
		if configPair.Value == "" {
			return fmt.Errorf("invalid configuration value: %q", configPair.Value)
		}
	}

	return nil
}

func (cfg *Cfg) validateToken() error {
	if !cfg.Auth.Transitioning || len(cfg.Auth.Token) == 0 {
		return nil
	}

	log.Warn("Authentication is enabled but not enforced because transitioning=true. Gitaly will accept unauthenticated requests.")
	return nil
}

func (cfg *Cfg) validateInternalSocketDir() error {
	if cfg.InternalSocketDir == "" {
		return nil
	}

	if err := validateIsDirectory(cfg.InternalSocketDir, "internal_socket_dir"); err != nil {
		return err
	}

	if err := trySocketCreation(cfg.InternalSocketDir); err != nil {
		return fmt.Errorf("internal_socket_dir: try create socket: %w", err)
	}
	return nil
}

func trySocketCreation(dir string) error {
	// To validate the socket can actually be created, we open and close a socket.
	// Any error will be assumed persistent for when the gitaly-ruby sockets are created
	// and thus fatal at boot time
	b, err := text.RandomHex(4)
	if err != nil {
		return err
	}

	socketPath := filepath.Join(dir, fmt.Sprintf("test-%s.sock", b))
	defer func() { _ = os.Remove(socketPath) }()

	// Attempt to create an actual socket and not just a file to catch socket path length problems
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("socket could not be created in %s: %s", dir, err)
	}

	return l.Close()
}

// defaultMaintenanceWindow specifies a 10 minute job that runs daily at +1200
// GMT time
func defaultMaintenanceWindow(storages []Storage) DailyJob {
	storageNames := make([]string, len(storages))
	for i, s := range storages {
		storageNames[i] = s.Name
	}

	return DailyJob{
		Hour:     12,
		Minute:   0,
		Duration: Duration(10 * time.Minute),
		Storages: storageNames,
	}
}

func (cfg *Cfg) validateMaintenance() error {
	dm := cfg.DailyMaintenance

	sNames := map[string]struct{}{}
	for _, s := range cfg.Storages {
		sNames[s.Name] = struct{}{}
	}
	for _, sName := range dm.Storages {
		if _, ok := sNames[sName]; !ok {
			return fmt.Errorf("daily maintenance specified storage %q does not exist in configuration", sName)
		}
	}

	if dm.Hour > 23 {
		return fmt.Errorf("daily maintenance specified hour '%d' outside range (0-23)", dm.Hour)
	}
	if dm.Minute > 59 {
		return fmt.Errorf("daily maintenance specified minute '%d' outside range (0-59)", dm.Minute)
	}
	if dm.Duration.Duration() > 24*time.Hour {
		return fmt.Errorf("daily maintenance specified duration %s must be less than 24 hours", dm.Duration.Duration())
	}

	return nil
}

func (cfg *Cfg) validateCgroups() error {
	cg := cfg.Cgroups

	if cg.Count == 0 {
		return nil
	}

	if cg.Mountpoint == "" {
		return fmt.Errorf("cgroups.mountpoint: cannot be empty")
	}

	if cg.HierarchyRoot == "" {
		return fmt.Errorf("cgroups.hierarchy_root: cannot be empty")
	}

	if cg.CPU.Enabled && cg.CPU.Shares == 0 {
		return fmt.Errorf("cgroups.cpu.shares: has to be greater than zero")
	}

	if cg.Memory.Enabled && (cg.Memory.Limit == 0 || cg.Memory.Limit < -1) {
		return fmt.Errorf("cgroups.memory.limit: has to be greater than zero or equal to -1")
	}

	return nil
}

var (
	errPackObjectsCacheNegativeMaxAge = errors.New("pack_objects_cache.max_age cannot be negative")
	errPackObjectsCacheNoStorages     = errors.New("pack_objects_cache: cannot pick default cache directory: no storages")
	errPackObjectsCacheRelativePath   = errors.New("pack_objects_cache: storage directory must be absolute path")
)

func (cfg *Cfg) configurePackObjectsCache() error {
	poc := &cfg.PackObjectsCache
	if !poc.Enabled {
		return nil
	}

	if poc.MaxAge < 0 {
		return errPackObjectsCacheNegativeMaxAge
	}

	if poc.MaxAge == 0 {
		poc.MaxAge = Duration(5 * time.Minute)
	}

	if poc.Dir == "" {
		if len(cfg.Storages) == 0 {
			return errPackObjectsCacheNoStorages
		}

		poc.Dir = filepath.Join(cfg.Storages[0].Path, GitalyDataPrefix, "PackObjectsCache")
	}

	if !filepath.IsAbs(poc.Dir) {
		return errPackObjectsCacheRelativePath
	}

	return nil
}
