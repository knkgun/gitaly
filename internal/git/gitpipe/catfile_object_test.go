package gitpipe

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
)

func TestCatfileObject(t *testing.T) {
	cfg := testcfg.Build(t)

	repoProto, _ := gittest.CloneRepo(t, cfg, cfg.Storages[0])
	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	for _, tc := range []struct {
		desc              string
		catfileInfoInputs []CatfileInfoResult
		expectedResults   []CatfileObjectResult
		expectedErr       error
	}{
		{
			desc: "single blob",
			catfileInfoInputs: []CatfileInfoResult{
				{ObjectInfo: &catfile.ObjectInfo{Oid: lfsPointer1, Type: "blob", Size: 133}},
			},
			expectedResults: []CatfileObjectResult{
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: lfsPointer1, Type: "blob", Size: 133}}},
			},
		},
		{
			desc: "multiple blobs",
			catfileInfoInputs: []CatfileInfoResult{
				{ObjectInfo: &catfile.ObjectInfo{Oid: lfsPointer1, Type: "blob", Size: 133}},
				{ObjectInfo: &catfile.ObjectInfo{Oid: lfsPointer2, Type: "blob", Size: 127}},
				{ObjectInfo: &catfile.ObjectInfo{Oid: lfsPointer3, Type: "blob", Size: 127}},
				{ObjectInfo: &catfile.ObjectInfo{Oid: lfsPointer4, Type: "blob", Size: 129}},
			},
			expectedResults: []CatfileObjectResult{
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: lfsPointer1, Type: "blob", Size: 133}}},
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: lfsPointer2, Type: "blob", Size: 127}}},
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: lfsPointer3, Type: "blob", Size: 127}}},
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: lfsPointer4, Type: "blob", Size: 129}}},
			},
		},
		{
			desc: "revlist result with object names",
			catfileInfoInputs: []CatfileInfoResult{
				{ObjectInfo: &catfile.ObjectInfo{Oid: "b95c0fad32f4361845f91d9ce4c1721b52b82793", Type: "tree", Size: 43}},
				{ObjectInfo: &catfile.ObjectInfo{Oid: "93e123ac8a3e6a0b600953d7598af629dec7b735", Type: "blob", Size: 59}, ObjectName: []byte("branch-test.txt")},
			},
			expectedResults: []CatfileObjectResult{
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: "b95c0fad32f4361845f91d9ce4c1721b52b82793", Type: "tree", Size: 43}}},
				{Object: &catfile.Object{ObjectInfo: catfile.ObjectInfo{Oid: "93e123ac8a3e6a0b600953d7598af629dec7b735", Type: "blob", Size: 59}}, ObjectName: []byte("branch-test.txt")},
			},
		},
		{
			desc: "invalid object ID",
			catfileInfoInputs: []CatfileInfoResult{
				{ObjectInfo: &catfile.ObjectInfo{Oid: "invalidobjectid", Type: "blob"}},
			},
			expectedErr: errors.New("requesting object: object not found"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()

			catfileCache := catfile.NewCache(cfg)
			defer catfileCache.Stop()

			objectReader, err := catfileCache.ObjectReader(ctx, repo)
			require.NoError(t, err)

			it, err := CatfileObject(ctx, objectReader, NewCatfileInfoIterator(tc.catfileInfoInputs))
			require.NoError(t, err)

			var results []CatfileObjectResult
			for it.Next() {
				result := it.Result()

				objectData, err := io.ReadAll(result)
				require.NoError(t, err)
				require.Len(t, objectData, int(result.ObjectSize()))

				// We only really want to compare the publicly visible fields
				// containing info about the object itself, and not the object's
				// private state. We thus need to reconstruct the objects here.
				results = append(results, CatfileObjectResult{
					Object: &catfile.Object{
						ObjectInfo: catfile.ObjectInfo{
							Oid:  result.ObjectID(),
							Type: result.ObjectType(),
							Size: result.ObjectSize(),
						},
					},
					ObjectName: result.ObjectName,
				})
			}

			// We're converting the error here to a plain un-nested error such
			// that we don't have to replicate the complete error's structure.
			err = it.Err()
			if err != nil {
				err = errors.New(err.Error())
			}

			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedResults, results)
		})
	}
}
