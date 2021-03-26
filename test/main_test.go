package jwtauthv2_test

import (
	"jwtauthv2/api"
	"jwtauthv2/store/pgsql"
	"jwtauthv2/usecase/auth"
	"log"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Jwtauthv2Suite struct {
	db          *sqlx.DB
	router      *mux.Router
	userStore   *pgsql.UserPgSQL
	tokeStore   *pgsql.TokenPgSQL
	authService *auth.AuthService
}

var Start = new(Jwtauthv2Suite)

// Поднять БД, накатить данные и миграции

func TestMain(m *testing.M) {
	db, err := sqlx.Connect("pgx", os.Getenv("DSN_DB"))
	if err != nil {
		log.Fatal("sqlx.Connect ", err)
	}

	userStore := pgsql.NewUserPgSQL(db)
	tokenStore := pgsql.NewTokenPgSQL(db)
	authService := auth.NewAuthService(userStore, tokenStore)
	handler := api.NewHandler(authService)

	Start.db = db
	Start.userStore = userStore
	Start.tokeStore = tokenStore
	Start.authService = authService
	Start.router = handler.InitRouter()

	exitVal := m.Run()

	db.Close()
	os.Exit(exitVal)
}
