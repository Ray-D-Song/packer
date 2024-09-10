package model

type Library struct {
	Size    int    `json:"size"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
