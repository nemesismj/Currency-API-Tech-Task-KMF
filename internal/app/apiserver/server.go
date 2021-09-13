package apiserver

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"io/ioutil"
	"net/http"
	_ "techtask/docs"
	"techtask/internal/app/model"
	"techtask/internal/app/model/dto"
	"techtask/internal/app/store"
	"time"
)

// Server struct
type Server struct {
	config     *Config
	logger     *logrus.Logger
	router     *mux.Router
	store      *store.Store
	httpServer *http.Server
}

// HttpResponse struct
type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

// New object of server
func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Run - run api server func
func (s *Server) Run() error {

	if err := s.configureLogger(); err != nil {
		s.logger.Error(fmt.Sprintf("Произошла ошибка при конфигурации логгера %s", err.Error()))
		return err
	}
	//s.configureRouter()
	for {
		db, _ := s.configureStore()
		if db != nil {
			break
		}
	}
	if _, err := s.configureStore(); err != nil {
		s.logger.Error(fmt.Sprintf("Произошла ошибка при конфигурации стора %s", err.Error()))
	}
	s.logger.Info(fmt.Sprintf("Сервер запустился на порту %s, хост %s ", s.config.Port, s.config.Host))
	s.httpServer = &http.Server{
		Addr:           ":" + s.config.Port,
		Handler:        s.configureRouter(),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return http.ListenAndServe(s.config.Port, s.router)
}

// Shutdown - shutdown api server func
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.store.Close(); err != nil {
		logrus.Errorf("Произошла ошибка при отключении подключения к базе данных: %s", err.Error())
	}

	return s.httpServer.Shutdown(ctx)
}

// configureLogger - configure logger func
func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	} else {
		s.logger.SetLevel(level)
	}
	return nil
}

// configureStore - configure store func
func (s *Server) configureStore() (*store.Store, error) {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return nil, err
	}
	s.store = st
	return s.store, nil
}

// ServeHTTP - Serve HTTP  func
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter - configure router func
func (s *Server) configureRouter() *mux.Router {
	s.router.HandleFunc("/api/v1/currency/save/{date}", s.handleRCurrencyCreate()).Methods("GET")
	s.router.HandleFunc("/api/v1/currency/{date}/{code}", s.handleRCurrenciesList()).Methods("GET")
	s.router.HandleFunc("/api/v1/currency/{date}", s.handleRCurrenciesList()).Methods("GET")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return s.router
}

// @Summary RCurrency create
// @Tags create a new rcurrencies
// @Description создание r_currency на основе данных с api
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 422 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /api/v1/currency/save/{date} [get]
func (s *Server) handleRCurrencyCreate() http.HandlerFunc {
	type Response struct {
		Success bool               `json:"success"`
		Data    []*model.RCurrency `json:"data,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		dateUrl := vars["date"]
		results := s.AsyncGetRequestToGetRates(dateUrl)
		if results.response.StatusCode == http.StatusOK {
			resp := Response{
				Success: true,
			}
			bodyBytes, err := ioutil.ReadAll(results.response.Body)
			if err != nil {
				s.logger.Error(fmt.Sprintf("Произошла ошибка при считывании XML респонса с API %s", err.Error()))
			}
			rcurrencies := &dto.RCurrencies{}
			_ = xml.Unmarshal(bodyBytes, &rcurrencies)

			dateParsed, err := time.Parse("02.01.2006", rcurrencies.Date)
			var slice []*model.RCurrency
			for _, currency := range rcurrencies.Item {
				rc := &model.RCurrency{
					A_DATE: dateParsed,
				}
				rc.CODE = currency.CODE
				rc.VALUE = currency.VALUE
				rc.TITLE = currency.TITLE
				go func() {
					_, err := s.store.RCurrency().Create(rc)
					if err != nil {
						s.logger.Error(fmt.Sprintf("Произошла ошибка при сохранении сущности %s", err.Error()))
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
				}()
				slice = append(slice, rc)
			}
			resp.Data = slice
			s.respond(w, r, http.StatusOK, resp)
			s.logger.Info("Запрос к API выполнен успешно")

		} else {
			s.logger.Error("Запрос к API выполнен не успешно")
			s.customErrorMessage(w, r, http.StatusConflict, "Запрос к API выполнен не успешно")
			return
		}
		defer results.response.Body.Close()
	}
}

// @Summary RCurrencies list
// @Tags List of rcurrencies
// @Description Получение списка r_currency, pathvariable code  не обязательный
// @Produce json
// @Success 200 {object} []model.RCurrency
// @Success 204  {string} string	"No content"
// @Failure 409 {object} model.ErrorResponse
// @Router /api/v1/currency/{date}/{code} [get]
func (s *Server) handleRCurrenciesList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		pathVariables := mux.Vars(r)
		datePathVariable := pathVariables["date"]
		codePathVariable := pathVariables["code"]
		dateParsed, _ := time.Parse("02.01.2006", datePathVariable)
		rcurrencies, err := s.store.RCurrency().FindByDateAndCode(dateParsed, codePathVariable)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Произошла ошибка при получении списка валют %s", err.Error()))
			s.customErrorMessage(w, r, http.StatusConflict, "Произошла ошибка при получении списка валют")
		}
		if len(rcurrencies) < 1 {
			s.logger.Info(fmt.Sprintf("Запрос на получение данных из стора прошел успешно, но данных за такой период нет"))
			s.respond(w, r, http.StatusNoContent, rcurrencies)
		} else {
			s.logger.Info(fmt.Sprintf("Запрос на получение данных из стора прошел успешно"))
			s.respond(w, r, http.StatusOK, rcurrencies)
		}
	}
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) customErrorMessage(w http.ResponseWriter, r *http.Request, code int, err string) {
	s.respond(w, r, code, map[string]string{"error": err})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// AsyncGetRequestToGetRates - async get request for get rates api of national bank func
func (s *Server) AsyncGetRequestToGetRates(date string) *HttpResponse {
	ch := make(chan *HttpResponse) // buffered
	url := "https://nationalbank.kz/rss/get_rates.cfm?fdate"
	if date != "" || len(date) < 1 {
		url = fmt.Sprintf("%s=%s", url, date)
	} else {
		panic("error")
	}
	var responses *HttpResponse
	go func(url string) {
		s.logger.Info(fmt.Sprintf("Получение списка %s \n", url))
		resp, err := http.Get(url)
		ch <- &HttpResponse{url, resp, err}
	}(url)

	for {
		select {
		case r := <-ch:
			fmt.Println()
			s.logger.Info(fmt.Sprintf("%s список получен\n", r.url))
			responses = r
			return responses
		case <-time.After(50 * time.Millisecond):
			fmt.Print(".")
		}
	}
	return responses
}
