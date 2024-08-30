package vinfo

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

// Info struct provides version info.
type Info struct {
	GitVersion string `json:"git_version,omitempty"`
	GitCommit  string `json:"git_commit,omitempty"`
	GitBranch  string `json:"git_branch,omitempty"`
	BuildDate  string `json:"build_date,omitempty"`
	Platform   string `json:"platform,omitempty"`
}

func (i *Info) JSONString() (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("failed to marshal version: %w", err)
	}

	return string(b), err
}

func (i Info) String() string {
	return fmt.Sprintf("%s\n%s",
		color.GreenString("%s %s", i.GitVersion, i.Platform),
		color.YellowString("%s %s %s", i.BuildDate, i.GitCommit, i.GitBranch),
	)
}
