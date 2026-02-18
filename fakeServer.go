package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

// statusHandler возвращает ответ с кодом, который передан
// в заголовке X-Status. Например:
//
//	X-Status = 200 -> вернет ответ с кодом 200
//	X-Status = 404 -> вернет ответ с кодом 404
//	X-Status = 503 -> вернет ответ с кодом 503
//
// Если заголовок отстутствует, возвращает ответ с кодом 200.
// Тело ответа пустое.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("X-Status")
	if header == "" {
		w.WriteHeader(http.StatusOK)
        return
	}
	code, err := strconv.Atoi(header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(code)
}

// echoHandler возвращает ответ с тем же телом
// и заголовком Content-Type, которые пришли в запросе
func echoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqHeaderCT := r.Header.Get("Content-Type")  
	w.Header().Set("Content-Type", reqHeaderCT)
	w.Write(reqBody)
}

// jsonHandler проверяет, что Content-Type = application/json,
// а в теле запроса пришел валидный JSON,
// после чего возвращает ответ с кодом 200.
// Если какая-то проверка не прошла — возвращает ответ с кодом 400.
// Тело ответа пустое.
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	reqHeaderCT := r.Header.Get("Content-Type")
	if reqHeaderCT != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	reqJSON, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !json.Valid(reqJSON){
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func startServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler)
	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("/json", jsonHandler)
	return httptest.NewServer(mux)
}

func main() {
	server := startServer()
	defer server.Close()
	client := server.Client()

	{
		uri := server.URL + "/status"
		req, _ := http.NewRequest(http.MethodGet, uri, nil)
		req.Header.Add("X-Status", "201")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Status)
		// 201 Created
	}

	{
		uri := server.URL + "/echo"
		reqBody := []byte("hello world")
		resp, err := client.Post(uri, "text/plain", bytes.NewReader(reqBody))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Println(resp.Status)
		fmt.Println(string(respBody))
		// 200 OK
		// hello world
	}

	{
		uri := server.URL + "/json"
		reqBody, _ := json.Marshal(map[string]bool{"ok": true})
		resp, err := client.Post(uri, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Status)
		// 200 OK
	}
}
