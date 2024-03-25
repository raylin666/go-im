package types

type HeaderAppID struct {
	Key     uint64 `json:"key"`
	Secret  string `json:"secret"`
	UserId  string `json:"user_id"`
	Usersig string `json:"usersig"`
}
