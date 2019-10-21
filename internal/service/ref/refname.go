package ref

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"gitlab.com/gitlab-org/gitaly/internal/git"
	"gitlab.com/gitlab-org/gitaly/internal/helper"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
)

// FindRefName returns a ref that starts with the given prefix, if one exists.
//  If there is more than one such ref there is no guarantee which one is
//  returned or that the same one is returned on each call.
func (s *server) FindRefName(ctx context.Context, in *gitalypb.FindRefNameRequest) (*gitalypb.FindRefNameResponse, error) {
	if in.CommitId == "" {
		return nil, helper.ErrInvalidArgument(fmt.Errorf("empty commit sha"))
	}

	ref, err := findRefName(ctx, in.Repository, in.CommitId, string(in.Prefix))
	if err != nil {
		return nil, helper.ErrInternal(err)
	}

	return &gitalypb.FindRefNameResponse{Name: []byte(ref)}, nil
}

// We assume `repo` and `commitID` and `prefix` are non-empty
func findRefName(ctx context.Context, repo *gitalypb.Repository, commitID, prefix string) (string, error) {
	flags := []git.Option{
		git.Flag{"--format=%(refname)"},
		git.Flag{"--count=1"},
		git.ValueFlag{Name: "--contains", Value: commitID},
	}

	cmd, err := git.SafeCmd(ctx, repo, nil, git.SubCmd{
		Name:  "for-each-ref",
		Flags: flags,
		Args:  []string{prefix},
	})
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(cmd)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	refName := scanner.Text()

	if err := cmd.Wait(); err != nil {
		// We're suppressing the error since invalid commits isn't an error
		//  according to Rails
		return "", nil
	}

	// Trailing spaces are not allowed per the documentation
	//  https://www.kernel.org/pub/software/scm/git/docs/git-check-ref-format.html
	return strings.TrimSpace(refName), nil
}
