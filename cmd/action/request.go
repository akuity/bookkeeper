package action

import (
	"fmt"

	"github.com/akuity/bookkeeper"
	libOS "github.com/akuity/bookkeeper/internal/os"
)

func request() (bookkeeper.RenderRequest, error) {
	req := bookkeeper.RenderRequest{
		RepoCreds: bookkeeper.RepoCredentials{
			Username: "git",
		},
		Images: libOS.GetStringSliceFromEnvVar("INPUT_IMAGES", nil),
	}
	repo, err := libOS.GetRequiredEnvVar("GITHUB_REPOSITORY")
	if err != nil {
		return req, err
	}
	req.RepoURL = fmt.Sprintf("https://github.com/%s", repo)
	if req.RepoCreds.Password, err =
		libOS.GetRequiredEnvVar("INPUT_PERSONALACCESSTOKEN"); err != nil {
		return req, err
	}
	if req.Ref, err = libOS.GetRequiredEnvVar("GITHUB_SHA"); err != nil {
		return req, err
	}
	req.TargetBranch, err = libOS.GetRequiredEnvVar("INPUT_TARGETBRANCH")
	return req, err
}
