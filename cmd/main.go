package main

import (
	"jwtauthv2/api"
	"jwtauthv2/store/pgsql"
	"jwtauthv2/usecase/auth"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("pgx", os.Getenv("DSN_DB"))
	if err != nil {
		log.Fatal("sqlx.Connect ", err)
	}

	defer db.Close()

	userStore := pgsql.NewUserPgSQL(db)
	tokenStore := pgsql.NewTokenPgSQL(db)
	authService := auth.NewAuthService(userStore, tokenStore)
	handler := api.NewHandler(authService)

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), handler.InitRouter()))

}
