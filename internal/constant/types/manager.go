package types

import "time"

type CreateRequest struct {
	Ident     string    `json:"ident"`
	Name      string    `json:"name"`
	Status    uint32    `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CreateResponse struct {
	Id        int       `json:"id"`
	Ident     string    `json:"ident"`
	Name      string    `json:"name"`
	Key       uint32    `json:"key"`
	Secret    string    `json:"secret"`
	Status    int8      `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateData struct {
	Ident     string    `json:"ident"`
	Name      string    `json:"name"`
	Key       uint32    `json:"key"`
	Secret    string    `json:"secret"`
	Status    int8      `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}
