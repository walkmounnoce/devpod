package provider

import (
	"fmt"
	"strings"

	"github.com/loft-sh/devpod/pkg/config"
	"github.com/loft-sh/devpod/pkg/git"
	"github.com/loft-sh/devpod/pkg/types"
)

var (
	WorkspaceSourceGit   = "git:"
	WorkspaceSourceLocal = "local:"
	WorkspaceSourceImage = "image:"
)

type Workspace struct {
	// ID is the workspace id to use
	ID string `json:"id,omitempty"`

	// UID is used to identify this specific workspace
	UID string `json:"uid,omitempty"`

	// Folder is the local folder where workspace related contents will be stored
	Folder string `json:"folder,omitempty"`

	// Picture is the project social media image
	Picture string `json:"picture,omitempty"`

	// Provider is the provider used to create this workspace
	Provider WorkspaceProviderConfig `json:"provider,omitempty"`

	// Machine is the machine to use for this workspace
	Machine WorkspaceMachineConfig `json:"machine,omitempty"`

	// IDE holds IDE specific settings
	IDE WorkspaceIDEConfig `json:"ide,omitempty"`

	// Source is the source where this workspace will be created from
	Source WorkspaceSource `json:"source,omitempty"`

	// DevContainerPath is the relative path where the devcontainer.json is located.
	DevContainerPath string `json:"devContainerPath,omitempty"`

	// CreationTimestamp is the timestamp when this workspace was created
	CreationTimestamp types.Time `json:"creationTimestamp,omitempty"`

	// LastUsedTimestamp holds the timestamp when this workspace was last accessed
	LastUsedTimestamp types.Time `json:"lastUsed,omitempty"`

	// Context is the context where this config file was loaded from
	Context string `json:"context,omitempty"`

	// Origin is the place where this config file was loaded from
	Origin string `json:"-"`
}

type WorkspaceIDEConfig struct {
	// Name is the name of the IDE
	Name string `json:"name,omitempty"`

	// Options are the local options that override the global ones
	Options map[string]config.OptionValue `json:"options,omitempty"`
}

type WorkspaceMachineConfig struct {
	// ID is the machine ID to use for this workspace
	ID string `json:"machineId,omitempty"`

	// UID is the machine UID to use for this workspace
	UID string `json:"machineUid,omitempty"`

	// AutoDelete specifies if the machine should get destroyed when
	// the workspace is destroyed
	AutoDelete bool `json:"autoDelete,omitempty"`
}

type WorkspaceProviderConfig struct {
	// Name is the provider name
	Name string `json:"name,omitempty"`

	// Options are the local options that override the global ones
	Options map[string]config.OptionValue `json:"options,omitempty"`
}

type WorkspaceSource struct {
	// GitRepository is the repository to clone
	GitRepository string `json:"gitRepository,omitempty"`

	// GitBranch is the branch to use
	GitBranch string `json:"gitBranch,omitempty"`

	// LocalFolder is the local folder to use
	LocalFolder string `json:"localFolder,omitempty"`

	// Image is the docker image to use
	Image string `json:"image,omitempty"`
}

type AgentWorkspaceInfo struct {
	// Workspace holds the workspace info
	Workspace *Workspace `json:"workspace,omitempty"`

	// Machine holds the machine info
	Machine *Machine `json:"machine,omitempty"`

	// Agent holds the agent info
	Agent ProviderAgentConfig `json:"agent,omitempty"`

	// Options holds the filled provider options for this workspace
	Options map[string]config.OptionValue `json:"options,omitempty"`

	// ContentFolder holds the folder where the content is stored
	ContentFolder string `json:"contentFolder,omitempty"`

	// Origin holds the folder where this config was loaded from
	Origin string `json:"-"`
}

func (w WorkspaceSource) String() string {
	if w.GitRepository != "" {
		if w.GitBranch != "" {
			return WorkspaceSourceGit + w.GitRepository + "@" + w.GitBranch
		}

		return WorkspaceSourceGit + w.GitRepository
	}

	if w.LocalFolder != "" {
		return WorkspaceSourceLocal + w.LocalFolder
	}

	return WorkspaceSourceImage + w.Image
}

func ParseWorkspaceSource(source string) (*WorkspaceSource, error) {
	if strings.HasPrefix(source, WorkspaceSourceGit) {
		gitRepo, gitBranch := git.NormalizeRepository(strings.TrimPrefix(source, WorkspaceSourceGit))
		return &WorkspaceSource{
			GitRepository: gitRepo,
			GitBranch:     gitBranch,
		}, nil
	} else if strings.HasPrefix(source, WorkspaceSourceLocal) {
		return &WorkspaceSource{
			LocalFolder: strings.TrimPrefix(source, WorkspaceSourceLocal),
		}, nil
	} else if strings.HasPrefix(source, WorkspaceSourceImage) {
		return &WorkspaceSource{
			Image: strings.TrimPrefix(source, WorkspaceSourceImage),
		}, nil
	}

	return nil, fmt.Errorf("invalid workspace source: %s", source)
}
