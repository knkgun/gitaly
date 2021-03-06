package gittest

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

const (
	committerName  = "Scrooge McDuck"
	committerEmail = "scrooge@mcduck.com"
)

type writeCommitConfig struct {
	branch             string
	parents            []git.ObjectID
	committerName      string
	message            string
	treeEntries        []TreeEntry
	alternateObjectDir string
}

// WriteCommitOption is an option which can be passed to WriteCommit.
type WriteCommitOption func(*writeCommitConfig)

// WithBranch is an option for WriteCommit which will cause it to update the update the given branch
// name to the new commit.
func WithBranch(branch string) WriteCommitOption {
	return func(cfg *writeCommitConfig) {
		cfg.branch = branch
	}
}

// WithMessage is an option for WriteCommit which will set the commit message.
func WithMessage(message string) WriteCommitOption {
	return func(cfg *writeCommitConfig) {
		cfg.message = message
	}
}

// WithParents is an option for WriteCommit which will set the parent OIDs of the resulting commit.
func WithParents(parents ...git.ObjectID) WriteCommitOption {
	return func(cfg *writeCommitConfig) {
		if parents != nil {
			cfg.parents = parents
		} else {
			// We're explicitly initializing parents here such that we can discern the
			// case where the commit should be created with no parents.
			cfg.parents = []git.ObjectID{}
		}
	}
}

// WithTreeEntries is an option for WriteCommit which will cause it to create a new tree and use it
// as root tree of the resulting commit.
func WithTreeEntries(entries ...TreeEntry) WriteCommitOption {
	return func(cfg *writeCommitConfig) {
		cfg.treeEntries = entries
	}
}

// WithCommitterName is an option for WriteCommit which will set the committer name.
func WithCommitterName(name string) WriteCommitOption {
	return func(cfg *writeCommitConfig) {
		cfg.committerName = name
	}
}

// WithAlternateObjectDirectory will cause the commit to be written into the given alternate object
// directory. This can either be an absolute path or a relative path. In the latter case the path
// is considered to be relative to the repository path.
func WithAlternateObjectDirectory(alternateObjectDir string) WriteCommitOption {
	return func(cfg *writeCommitConfig) {
		cfg.alternateObjectDir = alternateObjectDir
	}
}

// WriteCommit writes a new commit into the target repository.
func WriteCommit(t testing.TB, cfg config.Cfg, repoPath string, opts ...WriteCommitOption) git.ObjectID {
	t.Helper()

	var writeCommitConfig writeCommitConfig
	for _, opt := range opts {
		opt(&writeCommitConfig)
	}

	message := "message"
	if writeCommitConfig.message != "" {
		message = writeCommitConfig.message
	}
	stdin := bytes.NewBufferString(message)

	// The ID of an arbitrary commit known to exist in the test repository.
	parents := []git.ObjectID{"1a0b36b3cdad1d2ee32457c102a8c0b7056fa863"}
	if writeCommitConfig.parents != nil {
		parents = writeCommitConfig.parents
	}

	var tree string
	if len(writeCommitConfig.treeEntries) > 0 {
		tree = WriteTree(t, cfg, repoPath, writeCommitConfig.treeEntries).String()
	} else if len(parents) == 0 {
		// If there are no parents, then we set the root tree to the empty tree.
		tree = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
	} else {
		tree = parents[0].String() + "^{tree}"
	}

	if writeCommitConfig.committerName == "" {
		writeCommitConfig.committerName = committerName
	}

	// Use 'commit-tree' instead of 'commit' because we are in a bare
	// repository. What we do here is the same as "commit -m message
	// --allow-empty".
	commitArgs := []string{
		"-c", fmt.Sprintf("user.name=%s", writeCommitConfig.committerName),
		"-c", fmt.Sprintf("user.email=%s", committerEmail),
		"-C", repoPath,
		"commit-tree", "-F", "-", tree,
	}

	var env []string
	if writeCommitConfig.alternateObjectDir != "" {
		require.True(t, filepath.IsAbs(writeCommitConfig.alternateObjectDir),
			"alternate object directory must be an absolute path")
		require.NoError(t, os.MkdirAll(writeCommitConfig.alternateObjectDir, 0o755))

		env = append(env,
			fmt.Sprintf("GIT_OBJECT_DIRECTORY=%s", writeCommitConfig.alternateObjectDir),
			fmt.Sprintf("GIT_ALTERNATE_OBJECT_DIRECTORIES=%s", filepath.Join(repoPath, "objects")),
		)
	}

	for _, parent := range parents {
		commitArgs = append(commitArgs, "-p", parent.String())
	}

	stdout := ExecOpts(t, cfg, ExecConfig{
		Stdin: stdin,
		Env:   env,
	}, commitArgs...)
	oid, err := git.NewObjectIDFromHex(text.ChompBytes(stdout))
	require.NoError(t, err)

	if writeCommitConfig.branch != "" {
		ExecOpts(t, cfg, ExecConfig{
			Env: env,
		}, "-C", repoPath, "update-ref", "refs/heads/"+writeCommitConfig.branch, oid.String())
	}

	return oid
}

func authorEqualIgnoringDate(t testing.TB, expected *gitalypb.CommitAuthor, actual *gitalypb.CommitAuthor) {
	t.Helper()
	require.Equal(t, expected.GetName(), actual.GetName(), "author name does not match")
	require.Equal(t, expected.GetEmail(), actual.GetEmail(), "author mail does not match")
}

// CommitEqual tests if two `GitCommit`s are equal
func CommitEqual(t testing.TB, expected, actual *gitalypb.GitCommit) {
	t.Helper()

	authorEqualIgnoringDate(t, expected.GetAuthor(), actual.GetAuthor())
	authorEqualIgnoringDate(t, expected.GetCommitter(), actual.GetCommitter())
	require.Equal(t, expected.GetBody(), actual.GetBody(), "body does not match")
	require.Equal(t, expected.GetSubject(), actual.GetSubject(), "subject does not match")
	require.Equal(t, expected.GetId(), actual.GetId(), "object ID does not match")
	require.Equal(t, expected.GetParentIds(), actual.GetParentIds(), "parent IDs do not match")
}
