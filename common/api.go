package common

//Command is the object for issuing commands
type Command struct {
	Name string `json:"name"`
	Args string `json:"args"`
}

//Response is the object for responding to commands
type Response struct {
	Command Command  `json:"command"`
	Type    string   `json:"type"`
	Answers []string `json:"answers"`
}
