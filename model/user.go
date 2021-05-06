package model

import (
	"time"
)

type SafeUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Usage    int64  `json:"usage"`
	Cap      int64  `json:"cap"`
}

type User struct {
	SafeUser
	Password string `json:"password"`
	// absolute path
	RootDir string `json:"root_dir"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
