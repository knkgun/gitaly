// +build postgres

package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/datastore/glsql"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/nodes"
	"gitlab.com/gitlab-org/gitaly/internal/praefect/service/info"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	"google.golang.org/grpc"
)

func getDB(t *testing.T) glsql.DB {
	return glsql.GetDB(t, "cmd_praefect")
}

func registerPraefectInfoServer(impl gitalypb.PraefectInfoServiceServer) svcRegistrar {
	return func(srv *grpc.Server) {
		gitalypb.RegisterPraefectInfoServiceServer(srv, impl)
	}
}

func TestDatalossSubcommand(t *testing.T) {
	mgr := &nodes.MockManager{
		GetShardFunc: func(vs string) (nodes.Shard, error) {
			var primary string
			switch vs {
			case "virtual-storage-1":
				primary = "gitaly-1"
			case "virtual-storage-2":
				primary = "gitaly-4"
			default:
				t.Error("unexpected virtual storage")
			}

			return nodes.Shard{Primary: &nodes.MockNode{
				GetStorageMethod: func() string {
					return primary
				},
			}}, nil
		},
	}

	cfg := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			{
				Name: "virtual-storage-1",
				Nodes: []*config.Node{
					{Storage: "gitaly-1"},
					{Storage: "gitaly-2"},
					{Storage: "gitaly-3"},
				},
			},
			{
				Name: "virtual-storage-2",
				Nodes: []*config.Node{
					{Storage: "gitaly-4"},
				},
			},
		},
	}

	gs := datastore.NewPostgresRepositoryStore(getDB(t), cfg.StorageNames())

	ctx, cancel := testhelper.Context()
	defer cancel()

	require.NoError(t, gs.SetGeneration(ctx, "virtual-storage-1", "repository-1", "gitaly-1", 1))
	require.NoError(t, gs.SetGeneration(ctx, "virtual-storage-1", "repository-1", "gitaly-2", 0))

	require.NoError(t, gs.SetGeneration(ctx, "virtual-storage-1", "repository-2", "gitaly-2", 0))
	require.NoError(t, gs.SetGeneration(ctx, "virtual-storage-1", "repository-2", "gitaly-3", 0))

	ln, clean := listenAndServe(t, []svcRegistrar{
		registerPraefectInfoServer(info.NewServer(mgr, cfg, nil, gs, nil))})
	defer clean()
	for _, tc := range []struct {
		desc            string
		args            []string
		virtualStorages []*config.VirtualStorage
		output          string
		error           error
	}{
		{
			desc:  "positional arguments",
			args:  []string{"-virtual-storage=virtual-storage-1", "positional-arg"},
			error: unexpectedPositionalArgsError{Command: "dataloss"},
		},
		{
			desc: "data loss with read-only repositories",
			args: []string{"-virtual-storage=virtual-storage-1"}, output: `Virtual storage: virtual-storage-1
  Outdated repositories:
    repository-2 (read-only):
      Primary: gitaly-1
      In-Sync Storages:
        gitaly-2, assigned host
        gitaly-3, assigned host
      Outdated Storages:
        gitaly-1 is behind by 1 change or less, assigned host
`,
		},
		{
			desc: "data loss with partially replicated repositories",
			args: []string{"-virtual-storage=virtual-storage-1", "-partially-replicated"}, output: `Virtual storage: virtual-storage-1
  Outdated repositories:
    repository-1 (writable):
      Primary: gitaly-1
      In-Sync Storages:
        gitaly-1, assigned host
      Outdated Storages:
        gitaly-2 is behind by 1 change or less, assigned host
        gitaly-3 is behind by 2 changes or less, assigned host
    repository-2 (read-only):
      Primary: gitaly-1
      In-Sync Storages:
        gitaly-2, assigned host
        gitaly-3, assigned host
      Outdated Storages:
        gitaly-1 is behind by 1 change or less, assigned host
`,
		},
		{
			desc:            "multiple virtual storages with read-only repositories",
			virtualStorages: []*config.VirtualStorage{{Name: "virtual-storage-2"}, {Name: "virtual-storage-1"}},
			output: `Virtual storage: virtual-storage-1
  Outdated repositories:
    repository-2 (read-only):
      Primary: gitaly-1
      In-Sync Storages:
        gitaly-2, assigned host
        gitaly-3, assigned host
      Outdated Storages:
        gitaly-1 is behind by 1 change or less, assigned host
Virtual storage: virtual-storage-2
  All repositories are writable!
`,
		},
		{
			desc:            "multiple virtual storages with partially replicated repositories",
			args:            []string{"-partially-replicated"},
			virtualStorages: []*config.VirtualStorage{{Name: "virtual-storage-2"}, {Name: "virtual-storage-1"}},
			output: `Virtual storage: virtual-storage-1
  Outdated repositories:
    repository-1 (writable):
      Primary: gitaly-1
      In-Sync Storages:
        gitaly-1, assigned host
      Outdated Storages:
        gitaly-2 is behind by 1 change or less, assigned host
        gitaly-3 is behind by 2 changes or less, assigned host
    repository-2 (read-only):
      Primary: gitaly-1
      In-Sync Storages:
        gitaly-2, assigned host
        gitaly-3, assigned host
      Outdated Storages:
        gitaly-1 is behind by 1 change or less, assigned host
Virtual storage: virtual-storage-2
  All repositories are up to date!
`,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			cmd := newDatalossSubcommand()
			output := &bytes.Buffer{}
			cmd.output = output

			fs := cmd.FlagSet()
			require.NoError(t, fs.Parse(tc.args))
			err := cmd.Exec(fs, config.Config{
				VirtualStorages: tc.virtualStorages,
				SocketPath:      ln.Addr().String(),
			})
			require.Equal(t, tc.error, err, err)
			require.Equal(t, tc.output, output.String())
		})
	}
}
