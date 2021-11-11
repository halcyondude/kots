package types

import (
	"time"

	"github.com/blang/semver"
	v1beta1 "github.com/replicatedhq/kots/kotskinds/apis/kots/v1beta1"
	"github.com/replicatedhq/kots/pkg/cursor"
	storetypes "github.com/replicatedhq/kots/pkg/store/types"
)

type Downstream struct {
	ClusterID        string `json:"id"`
	ClusterSlug      string `json:"slug"`
	Name             string `json:"name"`
	CurrentSequence  int64  `json:"currentSequence"`
	SnapshotSchedule string `json:"snapshotSchedule,omitempty"`
	SnapshotTTL      string `json:"snapshotTtl,omitempty"`
}

type DownstreamVersion struct {
	ChannelID                string                             `json:"channelID"`
	ChannelName              string                             `json:"channelName"`
	VersionLabel             string                             `json:"versionLabel"`
	UpdateCursor             string                             `json:"updateCursor"`
	ParsedCursor             cursor.Cursor                      `json:"-"`
	Semver                   *semver.Version                    `json:"semver,omitempty"`
	Status                   storetypes.DownstreamVersionStatus `json:"status"`
	CreatedOn                *time.Time                         `json:"createdOn"`
	ParentSequence           int64                              `json:"parentSequence"`
	Sequence                 int64                              `json:"sequence"`
	ReleaseNotes             string                             `json:"releaseNotes"`
	DeployedAt               *time.Time                         `json:"deployedAt"`
	Source                   string                             `json:"source"`
	PreflightResult          string                             `json:"preflightResult,omitempty"`
	PreflightResultCreatedAt *time.Time                         `json:"preflightResultCreatedAt,omitempty"`
	PreflightSkipped         bool                               `json:"preflightSkipped"`
	DiffSummary              string                             `json:"diffSummary,omitempty"`
	DiffSummaryError         string                             `json:"diffSummaryError,omitempty"`
	CommitURL                string                             `json:"commitUrl,omitempty"`
	GitDeployable            bool                               `json:"gitDeployable,omitempty"`
	UpstreamReleasedAt       *time.Time                         `json:"upstreamReleasedAt,omitempty"`
	YamlErrors               []v1beta1.InstallationYAMLError    `json:"yamlErrors,omitempty"`
}

type DownstreamVersions struct {
	ChannelID       string
	ChannelName     string
	CurrentVersion  *DownstreamVersion
	PendingVersions []DownstreamVersion
	PastVersions    []DownstreamVersion
	AllVersions     []DownstreamVersion
}

func (d DownstreamVersions) Len() int { return len(d.AllVersions) }

// Less will attempt to sort by semver, then by cursor, then by sequence.
func (d DownstreamVersions) Less(i, j int) bool {
	if d.AllVersions[i].Semver != nil || d.AllVersions[j].Semver != nil {
		return d.lessSemver(i, j)
	}
	// if d.AllVersions[i].ParsedCursor != nil && d.AllVersions[j].ParsedCursor != nil {
	// 	return d.lessCursor(i, j)
	// }
	return d.lessSequence(i, j)
}

func (d DownstreamVersions) lessSemver(i, j int) bool {
	// Treating releases with semver as newer than those without if on different channel.
	// Releases on the same channel will be sorted using cursor rules.
	if d.AllVersions[i].Semver == nil {
		return d.lessSequence(i, j)
	}

	if d.AllVersions[j].Semver == nil {
		return d.lessSequence(i, j)
	}

	if d.AllVersions[i].Semver.EQ((*d.AllVersions[j].Semver)) {
		return d.lessSequence(i, j)
	}

	return d.AllVersions[i].Semver.LT((*d.AllVersions[j].Semver))
}

func (d DownstreamVersions) lessCursor(i, j int) bool {
	if d.AllVersions[i].ParsedCursor == nil || d.AllVersions[j].ParsedCursor == nil {
		return d.lessSequence(i, j)
	}
	if d.AllVersions[i].ParsedCursor.Equal(d.AllVersions[j].ParsedCursor) {
		return d.lessSequence(i, j)
	}
	return d.AllVersions[i].ParsedCursor.Before(d.AllVersions[j].ParsedCursor)
}

func (d DownstreamVersions) lessSequence(i, j int) bool {
	return d.AllVersions[i].Sequence < d.AllVersions[j].Sequence
}

func (d DownstreamVersions) Swap(i, j int) {
	tmp := d.AllVersions[i]
	d.AllVersions[i] = d.AllVersions[j]
	d.AllVersions[j] = tmp
}

type DownstreamOutput struct {
	DryrunStdout string `json:"dryrunStdout"`
	DryrunStderr string `json:"dryrunStderr"`
	ApplyStdout  string `json:"applyStdout"`
	ApplyStderr  string `json:"applyStderr"`
	HelmStdout   string `json:"helmStdout"`
	HelmStderr   string `json:"helmStderr"`
	RenderError  string `json:"renderError"`
}
