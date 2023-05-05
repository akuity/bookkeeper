package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/akuity/bookkeeper"
)

func TestRequest(t *testing.T) {
	// We need to start by clearing these out, because these are all actually set
	// during a GitHub Actions Run -- which means these are sometimes set when
	// these tests run.
	t.Setenv("GITHUB_REPOSITORY", "")
	t.Setenv("INPUT_PERSONALACCESSTOKEN", "")
	t.Setenv("GITHUB_SHA", "")
	t.Setenv("INPUT_TARGETBRANCH", "")
	const (
		testRepo   = "krancour/foo"
		testImage1 = "krancour/foo:blue"
		testImage2 = "krancour/foo:green"
	)
	testReq := bookkeeper.RenderRequest{
		RepoURL: fmt.Sprintf("https://github.com/%s", testRepo),
		RepoCreds: bookkeeper.RepoCredentials{
			Username: "git",
			Password: "12345", // Like something an idiot would use for their luggage
		},
		Ref:          "1234567",
		TargetBranch: "env/dev",
		Images:       []string{testImage1, testImage2},
	}
	testCases := []struct {
		name       string
		setup      func()
		assertions func(bookkeeper.RenderRequest, error)
	}{
		{
			name: "GITHUB_REPOSITORY not specified",
			assertions: func(_ bookkeeper.RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "GITHUB_REPOSITORY")
			},
		},
		{
			name: "INPUT_PERSONALACCESSTOKEN not specified",
			setup: func() {
				t.Setenv("GITHUB_REPOSITORY", testRepo)
			},
			assertions: func(_ bookkeeper.RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "INPUT_PERSONALACCESSTOKEN")
			},
		},
		{
			name: "GITHUB_SHA not specified",
			setup: func() {
				t.Setenv("INPUT_PERSONALACCESSTOKEN", testReq.RepoCreds.Password)
			},
			assertions: func(_ bookkeeper.RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "GITHUB_SHA")
			},
		},
		{
			name: "INPUT_TARGETBRANCH not specified",
			setup: func() {
				t.Setenv("GITHUB_SHA", testReq.Ref)
			},
			assertions: func(_ bookkeeper.RenderRequest, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "INPUT_TARGETBRANCH")
			},
		},
		{
			name: "success",
			setup: func() {
				t.Setenv("INPUT_TARGETBRANCH", testReq.TargetBranch)
				t.Setenv(
					"INPUT_IMAGES",
					fmt.Sprintf("%s,%s", testImage1, testImage2))
			},
			assertions: func(req bookkeeper.RenderRequest, err error) {
				require.NoError(t, err)
				require.Equal(t, testReq, req)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setup != nil {
				testCase.setup()
			}
			testCase.assertions(request())
		})
	}
}
