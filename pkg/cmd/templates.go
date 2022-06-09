package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/stripe/stripe-cli/pkg/ansi"
	"github.com/stripe/stripe-cli/pkg/terminal"
)

//
// Public functions
//

// WrappedInheritedFlagUsages returns a string containing the usage information
// for all flags which were inherited from parent commands, wrapped to the
// terminal's width.
func WrappedInheritedFlagUsages(cmd *cobra.Command) string {
	return cmd.InheritedFlags().FlagUsagesWrapped(terminal.GetTerminalWidth())
}

// WrappedLocalFlagUsages returns a string containing the usage information
// for all flags specifically set in the current command, wrapped to the
// terminal's width.
func WrappedLocalFlagUsages(cmd *cobra.Command) string {
	return cmd.LocalFlags().FlagUsagesWrapped(terminal.GetTerminalWidth())
}

// WrappedRequestParamsFlagUsages returns a string containing the usage
// information for all request parameters flags, i.e. flags used in operation
// commands to set values for request parameters. The string is wrapped to the
// terminal's width.
func WrappedRequestParamsFlagUsages(cmd *cobra.Command) string {
	var sb strings.Builder

	// We're cheating a little bit in thie method: we're not actually wrapping
	// anything, just printing out the flag names and assuming that no name
	// will be long enough to go over the terminal's width.
	// We do this instead of using pflag's `FlagUsagesWrapped` function because
	// we don't want to print the types (all request parameters flags are
	// defined as strings in the CLI, but it would be confusing to print that
	// out as a lot of them are not strings in the API).
	// If/when we do add help strings for request parameters flags, we'll have
	// to do actual wrapping.
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if _, ok := flag.Annotations["request"]; ok {
			sb.WriteString(fmt.Sprintf("      --%s\n", flag.Name))
		}
	})

	return sb.String()
}

// WrappedNonRequestParamsFlagUsages returns a string containing the usage
// information for all non-request parameters flags. The string is wrapped to
// the terminal's width.
func WrappedNonRequestParamsFlagUsages(cmd *cobra.Command) string {
	nonRequestParamsFlags := pflag.NewFlagSet("request", pflag.ExitOnError)

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if _, ok := flag.Annotations["request"]; !ok {
			nonRequestParamsFlags.AddFlag(flag)
		}
	})

	return nonRequestParamsFlags.FlagUsagesWrapped(terminal.GetTerminalWidth())
}

//
// Private functions
//

func getUsageTemplate() string {
	return fmt.Sprintf(`%s{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{if .Annotations}}

%s{{range $index, $cmd := .Commands}}{{if (eq (index $.Annotations $cmd.Name) "webhooks")}}
  {{rpad $cmd.Name $cmd.NamePadding}} {{$cmd.Short}}{{end}}{{end}}

%s{{range $index, $cmd := .Commands}}{{if (eq (index $.Annotations $cmd.Name) "stripe")}}
  {{rpad $cmd.Name $cmd.NamePadding}} {{$cmd.Short}}{{end}}{{end}}

%s
  {{rpad "get" 29}} Quickly retrieve resources from Stripe
  {{rpad "charges" 29}} Make requests (capture, create, list, etc) on charges
  {{rpad "customers" 29}} Make requests (create, delete, list, etc) on customers
  {{rpad "payment_intents" 29}} Make requests (cancel, capture, confirm, etc) on payment intents
  {{rpad "..." 29}} %s

%s{{range $index, $cmd := .Commands}}{{if (not (or (index $.Annotations $cmd.Name) $cmd.Hidden))}}
  {{rpad $cmd.Name $cmd.NamePadding}} {{$cmd.Short}}{{end}}{{end}}{{else}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{WrappedLocalFlagUsages . | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{WrappedInheritedFlagUsages . | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`,
		ansi.Bold("Usage:"),
		ansi.Bold("Aliases:"),
		ansi.Bold("Examples:"),
		ansi.Bold("Webhook commands:"),
		ansi.Bold("Stripe commands:"),
		ansi.Bold("Resource commands:"),
		ansi.Italic("To see more resource commands, run `stripe resources help`"),
		ansi.Bold("Other commands:"),
		ansi.Bold("Available commands:"),
		ansi.Bold("Flags:"),
		ansi.Bold("Global flags:"),
	)
}

func init() {
	cobra.AddTemplateFunc("WrappedInheritedFlagUsages", WrappedInheritedFlagUsages)
	cobra.AddTemplateFunc("WrappedLocalFlagUsages", WrappedLocalFlagUsages)
	cobra.AddTemplateFunc("WrappedRequestParamsFlagUsages", WrappedRequestParamsFlagUsages)
	cobra.AddTemplateFunc("WrappedNonRequestParamsFlagUsages", WrappedNonRequestParamsFlagUsages)
}
