package projectvo

type ProjectMetaInfo struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	IsFilling     int    `json:"isFilling"`
	ProjectTypeId int    `json:"projectTypeId"`
}

type TableSimpleInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
