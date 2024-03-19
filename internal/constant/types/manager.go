package types

import "time"

type ManagerCreateRequest struct {
	Ident     string    `json:"ident"`
	Name      string    `json:"name"`
	Status    int8      `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

type ManagerCreateResponse struct {
	Id        int       `json:"id"`
	Ident     string    `json:"ident"`
	Name      string    `json:"name"`
	Key       uint64    `json:"key"`
	Secret    string    `json:"secret"`
	Status    int8      `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ManagerCreateData struct {
	Ident     string    `json:"ident"`
	Name      string    `json:"name"`
	Key       uint64    `json:"key"`
	Secret    string    `json:"secret"`
	Status    int8      `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}
