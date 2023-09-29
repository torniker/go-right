package response

// Country model
type Country struct {
	Name        string            `json:"name"`
	Code        string            `json:"code"`
	Region      string            `json:"region"`
	Subregion   string            `json:"subregion"`
}
