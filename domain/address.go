package domain

type Address struct {
	ID      uint64 `json:"id" pg:",pk"`
	Address string `json:"address" pg:",unique; type:varchar(64)"`
}
