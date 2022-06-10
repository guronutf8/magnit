package api

import (
	"encoding/json"
	"magnit/internal/DB"
	"magnit/internal/Url"
	"magnit/internal/conf"
	"net/http"
	"strconv"
)

type Server struct {
	server http.Server
	mux    *http.ServeMux
	DB     *DB.DB
}

func Init(cfg conf.Config) *Server {
	// настройки сервиса
	mux := http.NewServeMux()
	server := http.Server{Addr: ":" + strconv.Itoa(cfg.APIPort), Handler: mux}
	s := &Server{server: server, mux: mux}
	s.initRoutes()

	s.DB = DB.Init(cfg)

	var _ = server.ListenAndServe()
	return s
}

/**
Все доступные routes
*/
func (s *Server) initRoutes() {
	// Иконку не отдаем, делаем тот путь, что бы браузер нне дергал лишнее
	s.mux.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(404)
	})

	// Отдаем ТЗ текстом
	s.mux.HandleFunc("/", handlerIndex)
	// Основной метод сервиса
	s.mux.HandleFunc("/v1/short/url", s.handlerUrl)
}

/**
Отдаем ТЗ текстом
*/
func handlerIndex(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("Сервис сокращения URL" +
		"\nОписание API:" +
		"\n- Метод для сокращения полного URL в короткий" +
		"\n- POST /v1/short/url" +
		"\n- Тело запроса: {“originalURL”: “http://foo.bar/baz”}    " +
		"\n- Ответ: 200 OK" +
		"\n- Тело ответа: {“shortURL”: “http://localhost/c2nn10mn”}" +
		"\n- Метод для преобразования короткого URL в полный" +
		"\n- GET /v1/short/url?shortURL=http://localhost/c2nn10mn    " +
		"\n- Ответ: 301 Moved Permanently" +
		"\nТребования к сокращенному URL’у:" +
		"\n- Минимально возможная длина URL (без учета домена)" +
		"\n- Используемые символы: a-z, A-Z, 0-9" +
		"\nТребования:" +
		"\n- Код пишем на Golang" +
		"\n- В качестве db использовать postgresql" +
		"\n- Сервис должен запускать в docker контейнере" +
		"\nОпционально:" +
		"\n- Генерация openapi спеки на основании кода"))

}

/**
Основной метод сервиса
*/
func (s *Server) handlerUrl(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		// Если это GET, короткую ссылку форматируем в число, что бы в БД найти Оригинальный адрес по primary index ID
		originalURLs := request.URL.Query()["shortURL"]
		if len(originalURLs) < 1 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		url := Url.Url{Short: originalURLs[0]}
		// Get отдает правильный Status(301 если нашли оригинальный путь, и 404 если не нашли)
		originalUrl, status := url.Get(s.DB)
		writer.WriteHeader(status)
		if originalUrl != "" {
			// в Location пишем оригинальный путь
			writer.Header().Set("Location", originalUrl)
		}
	} else if request.Method == "POST" {
		// пришел запрос с оригинальной ссылкой, отдадим ему в ответ короткую
		decoder := json.NewDecoder(request.Body)
		var t struct{ OriginalURL string }
		err := decoder.Decode(&t)

		if err != nil { // оригинальной ссылки не было(
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// заполняем структуру для ссылки, сохраняем, и получаем которую, ее и отдадим
		url := Url.Url{Long: t.OriginalURL}
		short, status := url.Save(s.DB)

		writer.WriteHeader(status) // ставим правильный статус
		if short == "" {
			return
		}

		shortJSON := struct {
			ShortURL string
		}{ShortURL: short}
		// формируем json и отдаем
		bytes, err := json.MarshalIndent(shortJSON, "", "	")
		if err == nil {
			writer.Write(bytes)
		}
	}
}
func (s *Server) String() string {
	return "I am web service on addr: " + s.server.Addr
}
