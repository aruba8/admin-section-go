package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"time"
	"strings"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(getConfig().App.Port), a.Router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	u := user{}
	users, err := u.getUsers(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) addAccount(w http.ResponseWriter, r *http.Request) {
	var acc account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&acc); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := acc.insertAccount(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	acc.Password = ""
	respondWithJSON(w, http.StatusOK, acc)
}

func Guard(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := parseRequest(r)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "Authorization Error: "+err.Error())
			return
		}
		if !token.Valid || claims.ExpiresAt < time.Now().Unix() {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		f(w, r)
	}
}

func (a *App) tokenAuth(w http.ResponseWriter, r *http.Request) {
	var acc account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&acc); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if validateAccount(acc.Username, acc.Password, a.DB) {
		claims := ClaimsStruct{
			acc.Username,
			acc.ID,
			jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))
		respondWithJSON(w, http.StatusOK, tokenStruct{
			"JWT " + tokenString,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Authorization Error")
	}
}

func (a *App) tokenRefresh(w http.ResponseWriter, r *http.Request) {

	token, claims, err := parseRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization Error: "+err.Error())
		return
	}

	if token.Valid && claims.ExpiresAt > time.Now().Unix() {
		claims.IssuedAt = time.Now().Unix()
		claims.ExpiresAt = time.Now().Add(time.Minute * 5).Unix()
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))
		respondWithJSON(w, http.StatusOK, tokenStruct{
			"JWT" + tokenString,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Authorization Error")
	}
}

func parseRequest(r *http.Request) (*jwt.Token, ClaimsStruct, error) {
	authString := r.Header.Get("Authorization")
	tokenString := "";
	if authStringArr := strings.Split(authString, " "); len(authStringArr) == 2 && authStringArr[0] == "JWT" {
		tokenString = authStringArr[1]
	}
	claims := ClaimsStruct{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return token, claims, err
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", Guard(a.getUsers)).Methods("GET")
	a.Router.HandleFunc("/accounts", a.addAccount).Methods("POST")
	a.Router.HandleFunc("/api-token-auth", a.tokenAuth).Methods("POST")
	a.Router.HandleFunc("/api-token-refresh", a.tokenRefresh).Methods("POST")
}
