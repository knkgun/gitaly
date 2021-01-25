package git

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/internal/testhelper"
)

type ReaderFunc func([]byte) (int, error)

func (fn ReaderFunc) Read(b []byte) (int, error) { return fn(b) }

func TestLocalRepository_WriteBlob(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	pbRepo, repoPath, clean := testhelper.InitBareRepo(t)
	defer clean()

	// write attributes file so we can verify WriteBlob runs the files through filters as
	// appropriate
	require.NoError(t, ioutil.WriteFile(filepath.Join(repoPath, "info", "attributes"), []byte(`
crlf binary
lf   text
	`), os.ModePerm))

	repo := NewRepository(pbRepo, config.Config)

	for _, tc := range []struct {
		desc    string
		path    string
		input   io.Reader
		sha     string
		error   error
		content string
	}{
		{
			desc:  "error reading",
			input: ReaderFunc(func([]byte) (int, error) { return 0, assert.AnError }),
			error: fmt.Errorf("%w, stderr: %q", assert.AnError, []byte{}),
		},
		{
			desc:    "successful empty blob",
			input:   strings.NewReader(""),
			sha:     "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
			content: "",
		},
		{
			desc:    "successful blob",
			input:   strings.NewReader("some content"),
			sha:     "f0eec86f614944a81f87d879ebdc9a79aea0d7ea",
			content: "some content",
		},
		{
			desc:    "line endings not normalized",
			path:    "crlf",
			input:   strings.NewReader("\r\n"),
			sha:     "d3f5a12faa99758192ecc4ed3fc22c9249232e86",
			content: "\r\n",
		},
		{
			desc:    "line endings normalized",
			path:    "lf",
			input:   strings.NewReader("\r\n"),
			sha:     "8b137891791fe96927ad78e64b0aad7bded08bdc",
			content: "\n",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			sha, err := repo.WriteBlob(ctx, tc.path, tc.input)
			require.Equal(t, tc.error, err)
			if tc.error != nil {
				return
			}

			assert.Equal(t, tc.sha, sha)
			content, err := repo.ReadObject(ctx, sha)
			require.NoError(t, err)
			assert.Equal(t, tc.content, string(content))
		})
	}
}

func TestLocalRepository_FormatTag(t *testing.T) {
	for _, tc := range []struct {
		desc       string
		objectID   string
		objectType string
		tagName    []byte
		userName   []byte
		userEmail  []byte
		tagBody    []byte
		authorDate time.Time
		err        error
	}{
		// Just trivial tests here, most of this is tested in
		// internal/gitaly/service/operations/tags_test.go
		{
			desc:       "basic signature",
			objectID:   "0000000000000000000000000000000000000000",
			objectType: "commit",
			tagName:    []byte("my-tag"),
			userName:   []byte("root"),
			userEmail:  []byte("root@localhost"),
			tagBody:    []byte(""),
		},
		{
			desc:       "basic signature",
			objectID:   "0000000000000000000000000000000000000000",
			objectType: "commit",
			tagName:    []byte("my-tag\ninjection"),
			userName:   []byte("root"),
			userEmail:  []byte("root@localhost"),
			tagBody:    []byte(""),
			err:        FormatTagError{expectedLines: 4, actualLines: 5},
		},
		{
			desc:       "signature with fixed time",
			objectID:   "0000000000000000000000000000000000000000",
			objectType: "commit",
			tagName:    []byte("my-tag"),
			userName:   []byte("root"),
			userEmail:  []byte("root@localhost"),
			tagBody:    []byte(""),
			authorDate: time.Unix(12345, 0),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			signature, err := FormatTag(tc.objectID, tc.objectType, tc.tagName, tc.userName, tc.userEmail, tc.tagBody, tc.authorDate)
			if err != nil {
				require.Equal(t, tc.err, err)
				require.Equal(t, "", signature)
			} else {
				require.NoError(t, err)
				require.Contains(t, signature, "object ")
				require.Contains(t, signature, "tag ")
				require.Contains(t, signature, "tagger ")
			}
		})
	}
}

func TestLocalRepository_WriteTag(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	pbRepo, repoPath, clean := testhelper.NewTestRepo(t)
	defer clean()

	repo := NewRepository(pbRepo, config.Config)

	for _, tc := range []struct {
		desc       string
		objectID   string
		objectType string
		tagName    []byte
		userName   []byte
		userEmail  []byte
		tagBody    []byte
		authorDate time.Time
	}{
		// Just trivial tests here, most of this is tested in
		// internal/gitaly/service/operations/tags_test.go
		{
			desc:       "basic signature",
			objectID:   "c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd",
			objectType: "commit",
			tagName:    []byte("my-tag"),
			userName:   []byte("root"),
			userEmail:  []byte("root@localhost"),
			tagBody:    []byte(""),
		},
		{
			desc:       "signature with time",
			objectID:   "c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd",
			objectType: "commit",
			tagName:    []byte("tag-with-timestamp"),
			userName:   []byte("root"),
			userEmail:  []byte("root@localhost"),
			tagBody:    []byte(""),
			authorDate: time.Unix(12345, 0),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			tagObjID, err := repo.WriteTag(ctx, tc.objectID, tc.objectType, tc.tagName, tc.userName, tc.userEmail, tc.tagBody, tc.authorDate)
			require.NoError(t, err)

			repoTagObjID := testhelper.MustRunCommand(t, nil, "git", "-C", repoPath, "rev-parse", tagObjID)
			require.Equal(t, text.ChompBytes(repoTagObjID), tagObjID)
		})
	}
}

func TestLocalRepository_ReadObject(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	testRepo, _, cleanup := testhelper.NewTestRepo(t)
	defer cleanup()

	repo := NewRepository(testRepo, config.Config)

	for _, tc := range []struct {
		desc    string
		oid     string
		content string
		error   error
	}{
		{
			desc:  "invalid object",
			oid:   ZeroOID.String(),
			error: InvalidObjectError(ZeroOID.String()),
		},
		{
			desc: "valid object",
			// README in gitlab-test
			oid:     "3742e48c1108ced3bf45ac633b34b65ac3f2af04",
			content: "Sample repo for testing gitlab features\n",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			content, err := repo.ReadObject(ctx, tc.oid)
			require.Equal(t, tc.error, err)
			require.Equal(t, tc.content, string(content))
		})
	}
}
