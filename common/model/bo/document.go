package bo

type DocumentBo struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DisplayPath  string `json:"displayPath"`
	Path         string `json:"path"`
	Suffix       string `json:"suffix"`
	Icon         string `json:"icon"`
	ResourceType int    `json:"resourceType"`
	RecycleFlag  int    `json:"recycleFlag"`
}
