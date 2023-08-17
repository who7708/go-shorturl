package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

// App encapsulates Env, Router and middleware
type App struct {
	Router *mux.Router
	// Middlewares *Middleware
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	Shortlink string `json:"shortlink"`
}

// 初始化
func (a *App) Initialize() {
	// set log formatter
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	// a.Middlewares = &Middleware{}
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.shortlinkRedirect).Methods("GET")
}

func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest,
			fmt.Errorf("参数解析错误 %v", r.Body)})
		return
	}
	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest,
			fmt.Errorf("参数校验错误 %v", req)})
		return
	}
	defer r.Body.Close()

	fmt.Printf("%v\n", req)
}

func respondWithError(w http.ResponseWriter, err error) {
	// panic("unimplemented")
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s\n", e.Status(), e)
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// panic("unimplemented")

	resp, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	fmt.Printf("%s\n", s)
}

func (a *App) shortlinkRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Printf("%s\n", vars["shortlink"])
}

// Run starts listen and server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
