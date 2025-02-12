package term

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

type CmdConfig struct {
	Cmd   string
	Alias string
	Desc  string
	Repl  bool
}

var CliCommands = []CmdConfig{
	{"", "", "start the Plandex REPL", false},
	{"new", "", "start a new plan", true},
	{"plans", "pl", "list plans", true},
	{"cd", "", "set current plan by name or index", true},
	{"current", "cu", "show current plan", true},
	{"rename", "", "rename the current plan", true},
	{"delete-plan", "dp", "delete plan by name or index", true},

	{"config", "", "show current plan config", true},
	{"set-config", "", "update current plan config", true},
	{"config default", "", "show default config for new plans", true},
	{"set-config default", "", "update default config for new plans", true},

	{"tell", "t", "describe a task to complete", false},
	{"chat", "ch", "ask a question or chat", false},

	{"load", "l", "load files/dirs/urls/notes/images or pipe data into context", true},
	{"ls", "", "list everything in context", true},
	{"rm", "", "remove context by index, range, name, or glob", true},
	{"clear", "", "remove all context", true},
	{"update", "u", "update outdated context", true},
	{"show", "", "show current context by name or index", true},

	// {"changes", "", "review pending changes in a TUI", false},
	{"diff", "", "review pending changes", true},
	{"diff --git", "", "review pending changes in 'git diff' format", true},
	{"diff --plain", "", "review pending changes in 'git diff' format with no color formatting", false},
	{"summary", "", "show the latest summary of the current plan", true},

	{"apply", "ap", "apply pending changes to project files", true},
	{"reject", "rj", "reject pending changes to one or more project files", true},

	{"log", "", "show log of plan updates", true},
	{"rewind", "rw", "rewind to a previous state", true},

	{"continue", "c", "continue the plan", true},
	{"debug", "db", "repeatedly run a command and auto-apply fixes until it succeeds", true},
	{"build", "b", "build any pending changes", true},

	{"convo", "", "show plan conversation", true},
	{"convo 1", "", "show a specific message in the conversation", false},
	{"convo 2-5", "", "show a range of messages in the conversation", false},
	{"convo --plain", "", "show conversation in plain text", false},

	{"branches", "br", "list plan branches", true},
	{"checkout", "co", "checkout or create a branch", true},
	{"delete-branch", "db", "delete a branch by name or index", true},

	{"plans --archived", "", "list archived plans", true},
	{"archive", "arc", "archive a plan", true},
	{"unarchive", "unarc", "unarchive a plan", true},

	{"models", "", "show current plan model settings", true},
	{"models default", "", "show org-wide default model settings for new plans", true},
	{"models available", "", "show all available models", true},
	{"models available --custom", "", "show available custom models only", true},
	{"models delete", "", "delete a custom model", true},
	{"models add", "", "add a custom model", true},
	{"model-packs", "", "show all available model packs", true},
	{"model-packs create", "", "create a new custom model pack", true},
	{"model-packs delete", "", "delete a custom model pack", true},
	{"model-packs --custom", "", "show custom model packs only", true},
	{"set-model", "", "update current plan model settings", true},
	{"set-model default", "", "update org-wide default model settings for new plans", true},

	{"ps", "", "list active and recently finished plan streams", true},
	{"stop", "", "stop an active plan stream", true},
	{"connect", "conn", "connect to an active plan stream", true},

	{"sign-in", "", "sign in, accept an invite, or create an account", true},
	{"invite", "", "invite a user to join your org", true},
	{"revoke", "", "revoke an invite or remove a user from your org", true},
	{"users", "", "list users and pending invites in your org", true},

	{"credits", "", "show Plandex Cloud credits balance", true},
	{"usage", "", "show Plandex Cloud credits transaction log", true},
	{"billing", "", "show Plandex Cloud billing settings", true},
}

var CmdDesc = map[string]CmdConfig{}

func init() {
	for _, cmd := range CliCommands {
		CmdDesc[cmd.Cmd] = cmd
	}
}

func PrintCmds(prefix string, cmds ...string) {
	printCmds(os.Stderr, prefix, []color.Attribute{color.Bold, color.FgHiWhite, color.BgCyan, color.FgHiWhite}, cmds...)
}

func PrintCmdsWithColors(prefix string, colors []color.Attribute, cmds ...string) {
	printCmds(os.Stderr, prefix, colors, cmds...)
}

func printCmds(w io.Writer, prefix string, colors []color.Attribute, cmds ...string) {
	if os.Getenv("PLANDEX_DISABLE_SUGGESTIONS") != "" {
		return
	}

	for _, cmd := range cmds {
		config, ok := CmdDesc[cmd]
		if !ok {
			continue
		}

		if IsRepl && !config.Repl {
			continue
		}

		alias := config.Alias
		desc := config.Desc

		if alias != "" {
			if IsRepl {
				cmd = fmt.Sprintf("%s (\\%s)", cmd, alias)
			} else {
				containsFull := strings.Contains(cmd, alias)

				if containsFull {
					cmd = strings.Replace(cmd, alias, fmt.Sprintf("(%s)", alias), 1)
				} else {
					cmd = fmt.Sprintf("%s (%s)", cmd, alias)
				}
			}

			// desc += color.New(color.FgWhite).Sprintf(" • alias → %s", color.New(color.Bold).Sprint(alias))
		}

		var styled string
		if IsRepl {
			styled = color.New(colors...).Sprintf(" \\%s ", cmd)
		} else if cmd == "" { // special case for the repl
			styled = color.New(colors...).Sprintf(" plandex ")
		} else {
			styled = color.New(colors...).Sprintf(" plandex %s ", cmd)
		}

		fmt.Fprintf(w, "%s%s 👉 %s\n", prefix, styled, desc)
	}

}

func PrintCustomCmd(prefix, cmd, alias, desc string) {
	cmd = strings.Replace(cmd, alias, fmt.Sprintf("(%s)", alias), 1)
	// desc += color.New(color.FgWhite).Sprintf(" • alias → %s", color.New(color.Bold).Sprint(alias))
	styled := color.New(color.Bold, color.FgHiWhite, color.BgCyan, color.FgHiWhite).Sprintf(" plandex %s ", cmd)
	fmt.Printf("%s%s 👉 %s\n", prefix, styled, desc)
}

// PrintCustomHelp prints the custom help output for the Plandex CLI
func PrintCustomHelp(all bool) {
	builder := &strings.Builder{}

	color.New(color.Bold, color.BgGreen, color.FgHiWhite).Fprintln(builder, " Usage ")
	color.New(color.Bold).Fprintln(builder, "  plandex [command] [flags]")
	color.New(color.Bold).Fprintln(builder, "  pdx [command] [flags]")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgGreen, color.FgHiWhite).Fprintln(builder, " Help ")
	color.New(color.Bold).Fprintln(builder, "  plandex help # show basic usage")
	color.New(color.Bold).Fprintln(builder, "  plandex help --all # show all commands")
	color.New(color.Bold).Fprintln(builder, "  plandex [command] --help")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgMagenta, color.FgHiWhite).Fprintln(builder, " Getting Started ")
	fmt.Fprintln(builder)
	fmt.Fprintf(builder, "  1 - Create a new plan in your project's root directory with %s\n", color.New(color.Bold, color.BgCyan, color.FgHiWhite).Sprint(" plandex new "))
	fmt.Fprintln(builder)
	fmt.Fprintf(builder, "  2 - Load any relevant context with %s\n", color.New(color.Bold, color.BgCyan, color.FgHiWhite).Sprint(" plandex load [file-path-or-url] "))
	fmt.Fprintln(builder)
	fmt.Fprintf(builder, "  3 - Describe a task to complete with %s\n", color.New(color.Bold, color.BgCyan, color.FgHiWhite).Sprint(" plandex tell "))
	fmt.Fprintln(builder)

	if all {

	} else {

		// in the same style as 'getting started' section, output See All Commands

		color.New(color.Bold, color.BgHiBlue, color.FgHiWhite).Fprintln(builder, " Use 'plandex help --all' or 'plandex help -a' for a list of all commands ")

	}

	fmt.Print(builder.String())
}

func PrintHelpAllCommands() {
	builder := &strings.Builder{}

	color.New(color.Bold, color.BgMagenta, color.FgHiWhite).Fprintln(builder, " Key Commands ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiMagenta}, "new", "load", "tell", "diff", "diff --ui", "apply", "reject", "debug", "chat")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Plans ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "new", "plans", "cd", "current", "delete-plan", "rename", "archive", "plans --archived", "unarchive")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Changes ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "diff", "diff --ui", "diff --plain", "changes", "apply", "reject")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Context ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "load", "ls", "rm", "update", "clear")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Branches ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "branches", "checkout", "delete-branch")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " History ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "log", "rewind", "convo", "convo 1", "convo 2-5", "convo --plain", "summary")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Control ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "tell", "continue", "build", "debug", "chat")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Streams ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "ps", "connect", "stop")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Config ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "config", "set-config", "config default", "set-config default")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " AI Models ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "models", "models default", "models available", "set-model", "set-model default", "models available --custom", "models add", "models delete", "model-packs", "model-packs --custom", "model-packs create", "model-packs delete")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Accounts ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "sign-in", "invite", "revoke", "users")
	fmt.Fprintln(builder)

	color.New(color.Bold, color.BgCyan, color.FgHiWhite).Fprintln(builder, " Cloud ")
	printCmds(builder, " ", []color.Attribute{color.Bold, ColorHiCyan}, "credits", "usage", "billing")

	fmt.Print(builder.String())
}
