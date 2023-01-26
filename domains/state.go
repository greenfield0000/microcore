package domains

type State struct {
	Id      *uint64 `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Sysname *string `json:"sysname,omitempty"`
}
