package conflicts

import (
	"testing"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/commit"
	hookservice "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/repository"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/ssh"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testserver"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	testhelper.Run(m)
}

func SetupConfigAndRepo(t testing.TB, bare bool) (config.Cfg, *gitalypb.Repository, string) {
	cfg := testcfg.Build(t)

	testcfg.BuildGitalyGit2Go(t, cfg)

	repo, repoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0], gittest.CloneRepoOpts{
		WithWorktree: !bare,
	})

	return cfg, repo, repoPath
}

func SetupConflictsServiceWithConfig(t testing.TB, cfg *config.Cfg, hookManager hook.Manager) gitalypb.ConflictsServiceClient {
	serverSocketPath := runConflictsServer(t, *cfg, hookManager)
	cfg.SocketPath = serverSocketPath

	client, conn := NewConflictsClient(t, serverSocketPath)
	t.Cleanup(func() { conn.Close() })

	return client
}

func SetupConflictsService(t testing.TB, bare bool, hookManager hook.Manager) (config.Cfg, *gitalypb.Repository, string, gitalypb.ConflictsServiceClient) {
	cfg, repo, repoPath := SetupConfigAndRepo(t, bare)

	client := SetupConflictsServiceWithConfig(t, &cfg, hookManager)

	return cfg, repo, repoPath, client
}

func runConflictsServer(t testing.TB, cfg config.Cfg, hookManager hook.Manager) string {
	return testserver.RunGitalyServer(t, cfg, nil, func(srv *grpc.Server, deps *service.Dependencies) {
		gitalypb.RegisterConflictsServiceServer(srv, NewServer(
			deps.GetCfg(),
			deps.GetHookManager(),
			deps.GetLocator(),
			deps.GetGitCmdFactory(),
			deps.GetCatfileCache(),
			deps.GetConnsPool(),
		))
		gitalypb.RegisterRepositoryServiceServer(srv, repository.NewServer(
			deps.GetCfg(),
			deps.GetRubyServer(),
			deps.GetLocator(),
			deps.GetTxManager(),
			deps.GetGitCmdFactory(),
			deps.GetCatfileCache(),
			deps.GetConnsPool(),
		))
		gitalypb.RegisterSSHServiceServer(srv, ssh.NewServer(
			deps.GetCfg(),
			deps.GetLocator(),
			deps.GetGitCmdFactory(),
			deps.GetTxManager(),
		))
		gitalypb.RegisterHookServiceServer(srv, hookservice.NewServer(deps.GetCfg(), deps.GetHookManager(), deps.GetGitCmdFactory(), deps.GetPackObjectsCache()))
		gitalypb.RegisterCommitServiceServer(srv, commit.NewServer(
			deps.GetCfg(),
			deps.GetLocator(),
			deps.GetGitCmdFactory(),
			deps.GetLinguist(),
			deps.GetCatfileCache(),
		))
	}, testserver.WithHookManager(hookManager))
}

func NewConflictsClient(t testing.TB, serverSocketPath string) (gitalypb.ConflictsServiceClient, *grpc.ClientConn) {
	connOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(serverSocketPath, connOpts...)
	if err != nil {
		t.Fatal(err)
	}

	return gitalypb.NewConflictsServiceClient(conn), conn
}
