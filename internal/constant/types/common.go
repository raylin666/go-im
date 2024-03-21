package types

type HeaderAppID struct {
	Key     uint64 `json:"key"`
	Secret  string `json:"secret"`
	Usersig string `json:"usersig"`
}
