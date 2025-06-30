package argv

// Argument represents a command line option
type Argument struct {
	Long        string `json:"long,omitempty"`        // the --long form for this option, or empty if none
	Short       string `json:"short,omitempty"`       // the -s short form for this option, or empty if none
	Cardinality string `json:"cardinality,omitempty"` // determines how many tokens will be present (possible values: zero, one, multiple)
	Required    bool   `json:"required,omitempty"`    // if true, this option must be present on the command line
	Positional  bool   `json:"positional,omitempty"`  // if true, this option will be looked for in the positional flags
	Separate    bool   `json:"separate,omitempty"`    // if true, each slice and map entry will have its own --flag
	Help        string `json:"help,omitempty"`        // the help text for this option
	Default     string `json:"default,omitempty"`     // default value for this option, in string form to be displayed in help text
	Placeholder string `json:"placeholder,omitempty"` // placeholder string in help
}

// Command represents a named subcommand, or the top-level Command
type Command struct {
	Name        string      `json:"name,omitempty"`
	Aliases     []string    `json:"aliases,omitempty"`
	Help        string      `json:"help,omitempty"`
	Arguments   []*Argument `json:"arguments,omitempty"`
	Subcommands []*Command  `json:"subcommands,omitempty"`
}
