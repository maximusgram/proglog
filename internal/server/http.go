package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// A produce request contains the record that the caller
// of our API wants appended to the log
type ProduceRequest struct {
	Record Record `json:"record"`
}

// A produce response tells the caller what offset
// the log stored the records under
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// A consume request specifies which records the caller
// of our API wants to read
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

// A consumer response to send back those to
// the caller, which they asked for in the request
type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	// unmarshalling the request into a struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("%s \n", err.Error())
		fmt.Println(req)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// using that struct to produce to the log
	// and getting the offset that the log stored the record under
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// marshalling and writing the result to the response
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	// unmarshalling the request into a struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// using that struct to get record
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// marshalling and writing the result to the response
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
