package backup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/client"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	gitalylog "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/log"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/setup"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/storage"
	praefectConfig "gitlab.com/gitlab-org/gitaly/v14/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/datastore/glsql"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testserver"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/protobuf/proto"
)

func TestManager_Create(t *testing.T) {
	cfg := testcfg.Build(t)

	gitalyAddr := testserver.RunGitalyServer(t, cfg, nil, setup.RegisterAll)

	path := testhelper.TempDir(t)

	hooksRepo, hooksRepoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0], gittest.CloneRepoOpts{
		RelativePath: "hooks",
	})
	require.NoError(t, os.Mkdir(filepath.Join(hooksRepoPath, "custom_hooks"), os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(hooksRepoPath, "custom_hooks/pre-commit.sample"), []byte("Some hooks"), os.ModePerm))

	noHooksRepo, _ := gittest.CloneRepo(t, cfg, cfg.Storages[0], gittest.CloneRepoOpts{
		RelativePath: "no-hooks",
	})
	emptyRepo, _ := gittest.InitRepo(t, cfg, cfg.Storages[0])
	nonexistentRepo := proto.Clone(emptyRepo).(*gitalypb.Repository)
	nonexistentRepo.RelativePath = "nonexistent"

	for _, tc := range []struct {
		desc               string
		repo               *gitalypb.Repository
		createsBundle      bool
		createsCustomHooks bool
		err                error
	}{
		{
			desc:               "no hooks",
			repo:               noHooksRepo,
			createsBundle:      true,
			createsCustomHooks: false,
		},
		{
			desc:               "hooks",
			repo:               hooksRepo,
			createsBundle:      true,
			createsCustomHooks: true,
		},
		{
			desc:               "empty repo",
			repo:               emptyRepo,
			createsBundle:      false,
			createsCustomHooks: false,
			err:                fmt.Errorf("manager: repository empty: %w", ErrSkipped),
		},
		{
			desc:               "nonexistent repo",
			repo:               nonexistentRepo,
			createsBundle:      false,
			createsCustomHooks: false,
			err:                fmt.Errorf("manager: repository empty: %w", ErrSkipped),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			repoPath := filepath.Join(cfg.Storages[0].Path, tc.repo.RelativePath)
			refsPath := filepath.Join(path, tc.repo.RelativePath+".refs")
			bundlePath := filepath.Join(path, tc.repo.RelativePath+".bundle")
			customHooksPath := filepath.Join(path, tc.repo.RelativePath, "custom_hooks.tar")

			ctx, cancel := testhelper.Context()
			defer cancel()

			fsBackup := NewManager(NewFilesystemSink(path), LegacyLocator{})
			err := fsBackup.Create(ctx, &CreateRequest{
				Server:     storage.ServerInfo{Address: gitalyAddr, Token: cfg.Auth.Token},
				Repository: tc.repo,
			})
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.err, err)
			}

			if tc.createsBundle {
				require.FileExists(t, refsPath)
				require.FileExists(t, bundlePath)

				dirInfo, err := os.Stat(filepath.Dir(bundlePath))
				require.NoError(t, err)
				require.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm(), "expecting restricted directory permissions")

				bundleInfo, err := os.Stat(bundlePath)
				require.NoError(t, err)
				require.Equal(t, os.FileMode(0o600), bundleInfo.Mode().Perm(), "expecting restricted file permissions")

				output := gittest.Exec(t, cfg, "-C", repoPath, "bundle", "verify", bundlePath)
				require.Contains(t, string(output), "The bundle records a complete history")

				expectedRefs := gittest.Exec(t, cfg, "-C", repoPath, "show-ref", "--head")
				actualRefs := testhelper.MustReadFile(t, refsPath)
				require.Equal(t, string(expectedRefs), string(actualRefs))
			} else {
				require.NoFileExists(t, bundlePath)
			}

			if tc.createsCustomHooks {
				require.FileExists(t, customHooksPath)
			} else {
				require.NoFileExists(t, customHooksPath)
			}
		})
	}
}

func TestManager_Restore(t *testing.T) {
	cfg := testcfg.Build(t)
	testhelper.BuildGitalyHooks(t, cfg)

	gitalyAddr := testserver.RunGitalyServer(t, cfg, nil, setup.RegisterAll)

	testManagerRestore(t, cfg, gitalyAddr)
}

func TestManager_Restore_praefect(t *testing.T) {
	gitalyCfg := testcfg.Build(t, testcfg.WithStorages("gitaly-1"))

	testhelper.BuildPraefect(t, gitalyCfg)
	testhelper.BuildGitalyHooks(t, gitalyCfg)

	gitalyAddr := testserver.RunGitalyServer(t, gitalyCfg, nil, setup.RegisterAll, testserver.WithDisablePraefect())

	db := glsql.NewDB(t)
	var database string
	require.NoError(t, db.QueryRow(`SELECT current_database()`).Scan(&database))
	dbConf := glsql.GetDBConfig(t, database)

	conf := praefectConfig.Config{
		SocketPath: testhelper.GetTemporaryGitalySocketFileName(t),
		VirtualStorages: []*praefectConfig.VirtualStorage{
			{
				Name: "default",
				Nodes: []*praefectConfig.Node{
					{Storage: gitalyCfg.Storages[0].Name, Address: gitalyAddr},
				},
			},
		},
		DB: dbConf,
		Failover: praefectConfig.Failover{
			Enabled:          true,
			ElectionStrategy: praefectConfig.ElectionStrategyPerRepository,
		},
		Replication: praefectConfig.DefaultReplicationConfig(),
		Logging: gitalylog.Config{
			Format: "json",
			Level:  "panic",
		},
	}

	tempDir := testhelper.TempDir(t)
	configFilePath := filepath.Join(tempDir, "config.toml")
	configFile, err := os.Create(configFilePath)
	require.NoError(t, err)
	defer testhelper.MustClose(t, configFile)

	require.NoError(t, toml.NewEncoder(configFile).Encode(&conf))
	require.NoError(t, configFile.Sync())

	cmd := exec.Command(filepath.Join(gitalyCfg.BinDir, "praefect"), "-config", configFilePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	require.NoError(t, cmd.Start())

	t.Cleanup(func() { _ = cmd.Wait() })
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	testManagerRestore(t, gitalyCfg, "unix://"+conf.SocketPath)
}

func testManagerRestore(t *testing.T, cfg config.Cfg, gitalyAddr string) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	cc, err := client.Dial(gitalyAddr, nil)
	require.NoError(t, err)
	defer func() { require.NoError(t, cc.Close()) }()
	repoClient := gitalypb.NewRepositoryServiceClient(cc)

	createRepo := func(t testing.TB, relativePath string) *gitalypb.Repository {
		t.Helper()

		repo := &gitalypb.Repository{
			StorageName:  "default",
			RelativePath: relativePath,
		}

		for i := 0; true; i++ {
			_, err := repoClient.CreateRepository(ctx, &gitalypb.CreateRepositoryRequest{Repository: repo})
			if err != nil {
				require.Regexp(t, "(no healthy nodes)|(no such file or directory)|(connection refused)", err.Error())
				require.Less(t, i, 100, "praefect doesn't serve for too long")
				time.Sleep(50 * time.Millisecond)
			} else {
				break
			}
		}

		return repo
	}

	path := testhelper.TempDir(t)

	existingRepo := createRepo(t, "existing")
	require.NoError(t, os.MkdirAll(filepath.Join(path, existingRepo.RelativePath), os.ModePerm))
	existingRepoBundlePath := filepath.Join(path, existingRepo.RelativePath+".bundle")
	gittest.BundleTestRepo(t, cfg, "gitlab-test.git", existingRepoBundlePath)

	existingRepoHooks := createRepo(t, "existing_hooks")
	existingRepoHooksBundlePath := filepath.Join(path, existingRepoHooks.RelativePath+".bundle")
	existingRepoHooksCustomHooksPath := filepath.Join(path, existingRepoHooks.RelativePath, "custom_hooks.tar")
	require.NoError(t, os.MkdirAll(filepath.Join(path, existingRepoHooks.RelativePath), os.ModePerm))
	testhelper.CopyFile(t, existingRepoBundlePath, existingRepoHooksBundlePath)
	testhelper.CopyFile(t, "../gitaly/service/repository/testdata/custom_hooks.tar", existingRepoHooksCustomHooksPath)

	missingBundleRepo := createRepo(t, "missing_bundle")
	missingBundleRepoAlwaysCreate := createRepo(t, "missing_bundle_always_create")

	nonexistentRepo := &gitalypb.Repository{
		StorageName:  "default",
		RelativePath: "nonexistent",
	}
	nonexistentRepoBundlePath := filepath.Join(path, nonexistentRepo.RelativePath+".bundle")
	testhelper.CopyFile(t, existingRepoBundlePath, nonexistentRepoBundlePath)

	for _, tc := range []struct {
		desc          string
		repo          *gitalypb.Repository
		alwaysCreate  bool
		expectExists  bool
		expectedPaths []string
		expectedErrAs error
		expectVerify  bool
	}{
		{
			desc:         "existing repo, without hooks",
			repo:         existingRepo,
			expectVerify: true,
			expectExists: true,
		},
		{
			desc: "existing repo, with hooks",
			repo: existingRepoHooks,
			expectedPaths: []string{
				"custom_hooks/pre-commit.sample",
				"custom_hooks/prepare-commit-msg.sample",
				"custom_hooks/pre-push.sample",
			},
			expectExists: true,
			expectVerify: true,
		},
		{
			desc:          "missing bundle",
			repo:          missingBundleRepo,
			expectedErrAs: ErrSkipped,
		},
		{
			desc:         "missing bundle, always create",
			repo:         missingBundleRepoAlwaysCreate,
			alwaysCreate: true,
			expectExists: true,
		},
		{
			desc:         "nonexistent repo",
			repo:         nonexistentRepo,
			expectVerify: true,
			expectExists: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			repoPath := filepath.Join(cfg.Storages[0].Path, tc.repo.RelativePath)
			bundlePath := filepath.Join(path, tc.repo.RelativePath+".bundle")

			fsBackup := NewManager(NewFilesystemSink(path), LegacyLocator{})
			err := fsBackup.Restore(ctx, &RestoreRequest{
				Server:       storage.ServerInfo{Address: gitalyAddr, Token: cfg.Auth.Token},
				Repository:   tc.repo,
				AlwaysCreate: tc.alwaysCreate,
			})
			if tc.expectedErrAs != nil {
				require.True(t, errors.Is(err, tc.expectedErrAs), err.Error())
			} else {
				require.NoError(t, err)
			}

			exists, err := repoClient.RepositoryExists(ctx, &gitalypb.RepositoryExistsRequest{
				Repository: tc.repo,
			})
			require.NoError(t, err)
			require.Equal(t, tc.expectExists, exists.Exists)

			if tc.expectVerify {
				output := gittest.Exec(t, cfg, "-C", repoPath, "bundle", "verify", bundlePath)
				require.Contains(t, string(output), "The bundle records a complete history")
			}

			for _, p := range tc.expectedPaths {
				require.FileExists(t, filepath.Join(repoPath, p))
			}
		})
	}
}

func TestResolveSink(t *testing.T) {
	isStorageServiceSink := func(expErrMsg string) func(t *testing.T, sink Sink) {
		return func(t *testing.T, sink Sink) {
			t.Helper()
			sssink, ok := sink.(*StorageServiceSink)
			require.True(t, ok)
			_, err := sssink.bucket.List(nil).Next(context.TODO())
			ierr, ok := err.(interface{ Unwrap() error })
			require.True(t, ok)
			terr := ierr.Unwrap()
			require.Contains(t, terr.Error(), expErrMsg)
		}
	}

	tmpDir := testhelper.TempDir(t)
	gsCreds := filepath.Join(tmpDir, "gs.creds")
	require.NoError(t, os.WriteFile(gsCreds, []byte(`
{
  "type": "service_account",
  "project_id": "hostfactory-179005",
  "private_key_id": "6253b144ccd94f50ce1224a73ffc48bda256d0a7",
  "private_key": "-----BEGIN PRIVATE KEY-----\nXXXX<KEY CONTENT OMMIT HERR> \n-----END PRIVATE KEY-----\n",
  "client_email": "303721356529-compute@developer.gserviceaccount.com",
  "client_id": "116595416948414952474",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/303724477529-compute%40developer.gserviceaccount.com"
}`), 0o655))

	for _, tc := range []struct {
		desc   string
		envs   map[string]string
		path   string
		verify func(t *testing.T, sink Sink)
		errMsg string
	}{
		{
			desc: "AWS S3",
			envs: map[string]string{
				"AWS_ACCESS_KEY_ID":     "test",
				"AWS_SECRET_ACCESS_KEY": "test",
				"AWS_REGION":            "us-east-1",
			},
			path:   "s3://bucket",
			verify: isStorageServiceSink("The AWS Access Key Id you provided does not exist in our records."),
		},
		{
			desc: "Google Cloud Storage",
			envs: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": gsCreds,
			},
			path:   "blob+gs://bucket",
			verify: isStorageServiceSink("storage.googleapis.com"),
		},
		{
			desc: "Azure Cloud File Storage",
			envs: map[string]string{
				"AZURE_STORAGE_ACCOUNT":   "test",
				"AZURE_STORAGE_KEY":       "test",
				"AZURE_STORAGE_SAS_TOKEN": "test",
			},
			path:   "blob+bucket+azblob://bucket",
			verify: isStorageServiceSink("https://test.blob.core.windows.net"),
		},
		{
			desc: "Filesystem",
			path: "/some/path",
			verify: func(t *testing.T, sink Sink) {
				require.IsType(t, &FilesystemSink{}, sink)
			},
		},
		{
			desc:   "undefined",
			path:   "some:invalid:path\x00",
			errMsg: `parse "some:invalid:path\x00": net/url: invalid control character in URL`,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Cleanup(testhelper.ModifyEnvironment(t, k, v))
			}
			sink, err := ResolveSink(context.TODO(), tc.path)
			if tc.errMsg != "" {
				require.EqualError(t, err, tc.errMsg)
				return
			}
			tc.verify(t, sink)
		})
	}
}

func TestResolveLocator(t *testing.T) {
	for _, tc := range []struct {
		locator     string
		expectedErr string
	}{
		{locator: "legacy"},
		{locator: "pointer"},
		{
			locator:     "unknown",
			expectedErr: "unknown locator: \"unknown\"",
		},
	} {
		t.Run(tc.locator, func(t *testing.T) {
			l, err := ResolveLocator(tc.locator, nil)

			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NotNil(t, l)
		})
	}
}
