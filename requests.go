package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error
type apiError struct {
	Error string
}

func makeHttpRequestHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle the error
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type ApiServer struct {
	listAddr string
	currDay  time.Time
}

func NewApiServer(listenAddress string) *ApiServer {
	return &ApiServer{
		listAddr: listenAddress,
		currDay:  time.Now(),
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/lastDaily/{ticker}/{days}", makeHttpRequestHandler(s.lastDaily))
	http.ListenAndServe(s.listAddr, router)
}
func (s *ApiServer) lastDaily(writer http.ResponseWriter, request *http.Request) error {
	if request.Method == "GET" {
		return s.getLastDaily(writer, request)
	}
	if request.Method == "POST" {
		return s.submitLastDaily(writer, request)
	}
	return nil
}
func (s *ApiServer) getLastDaily(writer http.ResponseWriter, request *http.Request) error {
	days, err := getDays(request)
	if err == nil {
		return s.getLastDailyRange(writer, request, days)
	}
	return s.getLastDailyTen(writer, request)
}
func (s *ApiServer) getLastDailyTen(writer http.ResponseWriter, request *http.Request) error {
	return nil
}
func (s *ApiServer) getLastDailyRange(writer http.ResponseWriter, request *http.Request, days int) error {
	return nil
}

func (s *ApiServer) submitLastDaily(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func getTicker(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["ticker"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
func getDays(r *http.Request) (int, error) {
	daysStr := mux.Vars(r)["days"]
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return days, fmt.Errorf("invalid days or none given %s", days)
	}
	return days, nil
}
func Hello() {
	var apiString = "https://api.polygon.io/v1/open-close/"
	fmt.Println(apiString)
}
