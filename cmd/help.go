package cmd

var RootCommandHelpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.FullName}} {{if .VisibleFlags}}[global options]{{end}}{{if .VisibleCommands}} [command [command options]]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Arguments}} [arguments...]{{end}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

{{range .VisibleFlagCategories}}
   {{if .Name}}{{.Name}}

   {{end}}{{$flglen := len .Flags}}{{range $i, $e := .Flags}}{{if eq (subtract $flglen $i) 1}}{{$e}}
{{else}}{{$e}}
   {{end}}{{end}}{{end}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}
{{- if len .Authors}}

AUTHOR{{template "authorsTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}{{if .VisibleFlagCategories}}

GLOBAL OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}

GLOBAL OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{template "copyrightTemplate" .}}{{end}}
`
