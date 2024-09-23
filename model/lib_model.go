package model

type Library struct {
	Size    int    `json:"size"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Perm    string `json:"perm"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
