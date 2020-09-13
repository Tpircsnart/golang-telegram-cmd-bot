package models

//TotalList : undertake TotalList to build struct
type TotalList struct {
	Remotes []struct {
		Label string `json:"label,omitempty"`
		Apps  []struct {
			Label   string `json:"label,omitempty"`
			TopCmds []struct {
				Label    string `json:"label,omitempty"`
				ShellSrc string `json:"shell_src,omitempty"`
			} `json:"top_cmds,omitempty"`
			MidCmds []struct {
				Label    string `json:"label,omitempty"`
				ShellSrc string `json:"shell_src,omitempty"`
			} `json:"mid_cmds,omitempty"`
			BotCmds []struct {
				Label    string `json:"label,omitempty"`
				ShellSrc string `json:"shell_src,omitempty"`
			} `json:"bot_cmds,omitempty"`
		} `json:"apps,omitempty"`
	} `json:"remotes,omitempty"`
}
