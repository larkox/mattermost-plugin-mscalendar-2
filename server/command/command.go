// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package command

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"

	"github.com/mattermost/mattermost-plugin-msoffice/server/config"
)

// Handler handles commands
type Handler struct {
	Config    *config.Config
	Context   *plugin.Context
	Args      *model.CommandArgs
	ChannelId string
}

func getHelp() string {
	help := `
TODO: help text.
`
	return codeBlock(fmt.Sprintf(
		help,
	))
}

// Register is a function that allows the runner to register commands with the mattermost server.
type RegisterFunc func(*model.Command) error

// Init should be called by the plugin to register all necessary commands
func Init(registerFunc RegisterFunc) {
	_ = registerFunc(&model.Command{
		Trigger:          config.CommandTrigger,
		DisplayName:      "TODO display name",
		Description:      "TODO description",
		AutoComplete:     true,
		AutoCompleteDesc: "TODO autocomplete desc",
		AutoCompleteHint: "TODO autocomplete hint",
	})
}

// Execute should be called by the plugin when a command invocation is received from the Mattermost server.
func (h *Handler) Handle() (string, error) {
	cmd, parameters, err := h.isValid()
	if err != nil {
		return "", err
	}

	auth, err := h.Config.IsAuthorizedAdmin(h.Args.UserId)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to get authorization. Please contact your system administrator.\nFailure")
	}
	if !auth {
		return "", errors.New("You must be authorized to use the plugin. Please contact your system administrator.")
	}

	handler := h.help
	switch cmd {
	case "info":
		handler = h.info
	case "connect":
		handler = h.connect
	}
	out, err := handler(parameters...)
	if err != nil {
		return "", errors.WithMessagef(err, "Command /%s %s failed", config.CommandTrigger, cmd)
	}

	return out, nil
}

func (h *Handler) isValid() (subcommand string, parameters []string, err error) {
	if h.Context == nil || h.Args == nil || h.Config.MattermostSiteURL == "" {
		return "", nil, errors.New("Invalid arguments to command.Handler")
	}

	split := strings.Fields(h.Args.Command)
	command := split[0]
	if command != "/"+config.CommandTrigger {
		return "", nil, errors.Errorf("%q is not a supported command. Please contact your system administrator.", command)
	}

	parameters = []string{}
	subcommand = ""
	if len(split) > 1 {
		subcommand = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}

	return subcommand, parameters, nil
}

func codeBlock(in string) string {
	return fmt.Sprintf("```\n%s\n```", in)
}
