package pgsql

import (
	"database/sql"
	"jwtauthv2"
	"jwtauthv2/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TokenPgSQL struct {
	db *sqlx.DB
}

func NewTokenPgSQL(db *sqlx.DB) *TokenPgSQL {
	return &TokenPgSQL{
		db: db,
	}
}

func (r *TokenPgSQL) Add(id *entity.ID, token string) error {
	query := `INSERT INTO tokens(id, refreshToken) VALUES($1, $2)`
	_, err := r.db.Exec(query, id, token)

	if err != nil {
		return &jwtauthv2.Error{
			Op:  "TokenPgSQL.Add",
			Err: err,
		}
	}

	return nil
}

func (r *TokenPgSQL) FindByToken(token string) (*entity.ID, error) {
	var user uuid.UUID
	query := `SELECT id FROM tokens WHERE refreshToken=$1`

	if err := r.db.Get(token, query, &user); err != nil {
		if err == sql.ErrNoRows {
			return nil, &jwtauthv2.Error{
				Code:    jwtauthv2.ENOTFOUND,
				Message: "Token not found",
			}
		}

		return nil, &jwtauthv2.Error{
			Op:  "TokenPgSQL.FindByToken",
			Err: err,
		}
	}

	return &user, nil
}

func (r *TokenPgSQL) RemoveByToken(refreshToken string) error {
	query := `DELETE FROM tokens WHERE refreshToken=$1`
	_, err := r.db.Exec(query, refreshToken)

	if err != nil {
		return &jwtauthv2.Error{
			Op:  "TokenPgSQL.RemoveByToken",
			Err: err,
		}
	}

	return nil
}

func (r *TokenPgSQL) RemoveByID(id *entity.ID) error {
	query := `DELETE FROM tokens WHERE id=$1`
	_, err := r.db.Exec(query, id)

	if err != nil {
		return &jwtauthv2.Error{
			Op:  "TokenPgSQL.RemoveByID",
			Err: err,
		}
	}

	return nil
}
