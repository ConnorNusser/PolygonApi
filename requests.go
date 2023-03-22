package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
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
	listAddr     string
	currDay      time.Time
	store        Storage
	polyInstance polygon.Client
}

func NewApiServer(listenAddress string, store Storage, client polygon.Client) *ApiServer {
	return &ApiServer{
		listAddr:     listenAddress,
		currDay:      time.Now(),
		store:        store,
		polyInstance: client,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/lastDaily/{ticker}", makeHttpRequestHandler(s.lastDaily))
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
func polygonApiCall(ticker string, date time.Time) *models.GetDailyOpenCloseAggParams {
	params := models.GetDailyOpenCloseAggParams{
		Ticker: ticker,
		Date:   models.Date(date),
	}

	return &params
}
func (s *ApiServer) getLastDaily(writer http.ResponseWriter, request *http.Request) error {
	days, err := getDays(request)
	if err == nil {
		return s.getLastDailyRange(writer, request, days)
	}
	return s.getLastDailyTen(writer, request)
}
func (s *ApiServer) getLastDailyTen(writer http.ResponseWriter, request *http.Request) error {
	getDaily := []*models.GetDailyOpenCloseAggResponse{}
	ticker := getTicker(request)
	n := 0
	for n < 10 {
		curr := time.Now().AddDate(0, 0, -n)
		if curr.Weekday() != time.Saturday || curr.Weekday() != time.Sunday {
			params := polygonApiCall(ticker, curr)
			getDailyInstance, err := s.polyInstance.GetDailyOpenCloseAgg(context.Background(), params)
			fmt.Println(getDailyInstance)
			fmt.Println("YO")
			if err != nil {
				getDaily = append(getDaily, getDailyInstance)
			}
		}
		n += 1
	}
	for _, value := range getDaily {
		fmt.Println(value)
		fmt.Println("hi this is connor")
		ds := newDailyStock(value.AfterHours, value.Close, value.From, value.High, value.Low, value.Open, value.PreMarket, value.Status, value.Symbol, value.Volume)
		s.store.CreateStock(ds)
	}

	return nil
}

func (s *ApiServer) getLastDailyRange(writer http.ResponseWriter, request *http.Request, days int) error {
	return nil
}

func (s *ApiServer) submitLastDaily(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func getTicker(r *http.Request) string {
	idStr := mux.Vars(r)["ticker"]
	return idStr
}
func getDays(r *http.Request) (int, error) {
	daysStr := mux.Vars(r)["days"]
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return days, fmt.Errorf("invalid days or none given %s", days)
	}
	return days, nil
}
