package common

type Command struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

type Response struct {
	Command Command  `json:"command"`
	Type    string   `json:"type"`
	Answers []string `json:"answers"`
}
