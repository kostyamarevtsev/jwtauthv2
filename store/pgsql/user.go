package pgsql

import (
	"database/sql"
	"jwtauthv2"
	"jwtauthv2/entity"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

type UserPgSQL struct {
	db *sqlx.DB
}

func NewUserPgSQL(db *sqlx.DB) *UserPgSQL {
	return &UserPgSQL{
		db: db,
	}
}

func (r *UserPgSQL) Create(user *entity.User) error {
	query := `INSERT INTO users(id, name, password) VALUES(:id, :name, :password)`
	_, err := r.db.NamedExec(query, user)

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
			return &jwtauthv2.Error{
				Code:    jwtauthv2.ECONFLICT,
				Message: "Name is already in use. Please choose a different Name",
			}
		}

		return &jwtauthv2.Error{
			Op:  "UserPgSQL.Create",
			Err: err,
		}
	}

	return nil
}

func (r *UserPgSQL) FindByName(name string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT id, name, password FROM users WHERE name=$1`

	if err := r.db.Get(user, query, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, &jwtauthv2.Error{
				Code:    jwtauthv2.ENOTFOUND,
				Message: "Name not found",
			}
		}

		return nil, &jwtauthv2.Error{
			Op:  "db.FindByName",
			Err: err,
		}
	}

	return user, nil
}
