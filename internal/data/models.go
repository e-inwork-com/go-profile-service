package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Profiles ProfileModelInterface
	Users    UserModelInterface
}

func InitModels(db *sql.DB, gRPCProfile string) Models {
	return Models{
		Profiles: ProfileModel{
			DB:          db,
			GRPCProfile: gRPCProfile,
		},
		Users: UserModel{DB: db},
	}
}
