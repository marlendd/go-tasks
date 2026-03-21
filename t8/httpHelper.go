package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Handy предоставляет удобный интерфейс
// для выполнения HTTP-запросов
type Handy struct{
	req *http.Request
	client *http.Client
	url string
	err error
}

// NewHandy создает новый экземпляр Handy
func NewHandy() *Handy {
	request, err := http.NewRequest("", "", nil)
	if err != nil {
		panic(err)
	}
	return &Handy{request, http.DefaultClient, "", nil}
}

// URL устанавливает URL, на который пойдет запрос
func (h *Handy) URL(uri string) *Handy {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	h.req.URL = parsedURL
	h.url = uri
	return h
}

// Client устанавливает HTTP-клиента
// вместо умолчательного http.DefaultClient
func (h *Handy) Client(client *http.Client) *Handy {
	h.client = client
	return h
}

// Header устанавливает значение заголовка
func (h *Handy) Header(key, value string) *Handy {
	h.req.Header.Add(key, value)
	return h
}

// Param устанавливает значение URL-параметра
func (h *Handy) Param(key, value string) *Handy {
	q := h.req.URL.Query()
    q.Add(key, value)
    h.req.URL.RawQuery = q.Encode()
	return h
}

// Form устанавливает данные, которые будут закодированы
// как application/x-www-form-urlencoded и отправлены в теле запроса
// с соответствующим content-type
func (h *Handy) Form(form map[string]string) *Handy {
	h.req.Header = make(http.Header)
	data := url.Values{}
	for key, val := range form {
		data.Add(key, val)
	}
	bodyStr := data.Encode()
	h.req.Body = io.NopCloser(bytes.NewReader([]byte(bodyStr)))
	h.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return h
}

// JSON устанавливает данные, которые будут закодированы
// как application/json и отправлены в теле запроса
// с соответствующим content-type
func (h *Handy) JSON(v any) *Handy {
	h.req.Header = make(http.Header)
	for key := range h.req.Header {
		h.req.Header.Del(key)
	}
	bodyBytes, err := json.Marshal(v)
	if err != nil {
		h.err = err
	}
	h.req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	h.req.Header.Set("Content-Type", "application/json") 
	return h
}

// Get выполняет GET-запрос с настроенными ранее параметрами
func (h *Handy) Get() *HandyResponse {
	if h.err != nil {
		return &HandyResponse{0, nil, h.err, nil}
	}
	h.req.Method = http.MethodGet
	resp, err := h.client.Do(h.req)
	if resp == nil {
		return &HandyResponse{0, nil, err, nil}
	}
	return &HandyResponse{resp.StatusCode, resp, err, nil}
}

// Post выполняет POST-запрос с настроенными ранее параметрами
func (h *Handy) Post() *HandyResponse {
	if h.err != nil {
		return &HandyResponse{0, nil, h.err, nil}
	}
	h.req.Method = http.MethodPost
	resp, err := h.client.Do(h.req)
	if resp == nil {
		return &HandyResponse{0, nil, err, nil}
	}
	return &HandyResponse{resp.StatusCode, resp, err, nil}
}

// HandyResponse представляет ответ на HTTP-запрос
type HandyResponse struct {
	StatusCode int
	Resp *http.Response
	Error error
	bodyBytes  []byte
}

// OK возвращает true, если во время выполнения запроса
// не произошло ошибок, а код HTTP-статуса ответа равен 200
func (r *HandyResponse) OK() bool {
	if (r.Error != nil || r.Resp == nil){
		return false
	}
	return r.StatusCode == http.StatusOK
}

// Bytes возвращает тело ответа как срез байт
func (r *HandyResponse) Bytes() []byte {
	if r.bodyBytes != nil {
        return r.bodyBytes
    }
	if r.Resp == nil || r.Resp.Body == nil {
        return nil
    }
	defer r.Resp.Body.Close()
	bytesResp, err := io.ReadAll(r.Resp.Body)
	r.Error = err
	r.bodyBytes = bytesResp
	return bytesResp
}

// String возвращает тело ответа как строку
func (r *HandyResponse) String() string {
	return string(r.Bytes())
}

// JSON декодирует тело ответа из JSON и сохраняет
// результат по адресу, на который указывает v
func (r *HandyResponse) JSON(v any) {
	// работает аналогично json.Unmarshal()
	// если при декодировании произошла ошибка,
	// она должна быть доступна через r.Err()
	defer r.Resp.Body.Close()
	bytesResp := r.Bytes()
	if len(bytesResp) == 0 {
		return
	}
	err := json.Unmarshal(bytesResp, &v)
	r.Error = err
}

// Err возвращает ошибку, которая возникла при выполнении запроса
// или обработке ответа
func (r *HandyResponse) Err() error {
	return r.Error
}


func main() {
	{
		// примеры запросов

		// GET-запрос с параметрами
		NewHandy().URL("https://httpbingo.org/get").Param("id", "42").Get()

		// HTTP-заголовки
		NewHandy().
			URL("https://httpbingo.org/get").
			Header("Accept", "text/html").
			Header("Authorization", "Bearer 1234567890").
			Get()

		// POST формы
		params := map[string]string{
			"brand":    "lg",
			"category": "tv",
		}
		NewHandy().URL("https://httpbingo.org/post").Form(params).Post()

		// POST JSON-документа
		NewHandy().URL("https://httpbingo.org/post").JSON(params).Post()
	}

	{
		// пример обработки ответа

		// отправляем GET-запрос с параметрами
		resp := NewHandy().URL("https://httpbingo.org/get").Param("id", "42").Get()
		if !resp.OK() {
			panic(resp.String())
		}

		// декодируем ответ в JSON
		var data map[string]any
		resp.JSON(&data)

		fmt.Println(data["url"])
		// "https://httpbingo.org/get"
		fmt.Println(data["args"])
		// map[id:[42]]
	}
}
