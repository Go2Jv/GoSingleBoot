package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64  `json:"id" bun:"id,pk,autoincrement"`
	Username      string `json:"username" bun:"username,notnull,unique"`
	Password      string `json:"password" bun:"password,notnull"`
}
