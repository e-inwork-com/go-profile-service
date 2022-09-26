package data

import (
	"database/sql"
	"errors"

	dataUser "github.com/e-inwork-com/go-user-service/pkg/data"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Profiles	ProfileModel
	Users 		dataUser.UserModel
}

func InitModels(db *sql.DB) Models {
	return Models{
		Profiles:	ProfileModel{DB: db},
		Users: 		dataUser.UserModel{DB: db},
	}
}

