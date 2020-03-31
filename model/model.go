package model

//Git struct
type Git struct {
	Repo []*Repo `json:"repo"`
}

//Repo struct
type Repo struct {
	RepoName  string       `json:"reponame"`
	Forks     int          `json:"forks"`
	Committee []*Committee `json:"committee"`
}

//Committee struct
type Committee struct {
	Name    string `json:"name"`
	Commits int    `json:"commits"`
}

//ResponseResult struct
type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
