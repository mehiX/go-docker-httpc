package data

//Image a Docker image representation
type Image struct {
	Containers  int               `json:"Containers"`
	Created     int               `json:"Created"`
	ID          string            `json:"Id"`
	Labels      map[string]string `json:"Labels"`
	ParentID    string            `json:"ParentId"`
	RepoDigests []string          `json:"RepoDigests"`
	RepoTags    []string          `json:"RepoTags"`
	SharedSize  int               `json:"SharedSize"`
	Size        int               `json:"Size"`
	VirtualSize int               `json:"VirtualSize"`
}
