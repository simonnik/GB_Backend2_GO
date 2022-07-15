package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type SignUp struct {
	login string
	pass  string
}

func (c RedisClient) SignupCompleteHandler(w http.ResponseWriter, r *http.Request) {
	code, err := c.checkCode(r)
	if err != nil {
		err = fmt.Errorf("check code: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusNotFound)
		_, err = fmt.Fprintln(w, "Code not found")
		if err != nil {
			log.Printf("[ERR] while get content: %v", err)
		}
		return
	}
	err = c.Delete(code.Code)
	if err != nil {
		err = fmt.Errorf("delete code: %w", err)
		log.Printf("[ERR] %v", err)
	}
	_, err = fmt.Fprintln(w, "Success")
	if err != nil {
		log.Printf("[ERR] while get content: %v", err)
	}
}

func (c RedisClient) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var su SignUp
	err := json.NewDecoder(r.Body).Decode(&su)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	code := 100000 + rand.Intn(99999)
	cd, err := c.Create(code)
	if err != nil {
		err = fmt.Errorf("create Code: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("Sms-code: %d", cd.Code)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c RedisClient) checkCode(r *http.Request) (*Code, error) {
	var code Code
	err := json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		return nil, err
	}

	ok, err := c.Check(code.Code)
	if err != nil {
		return nil, fmt.Errorf("check sms code %q: %w", code.Code, err)
	}
	if !ok {
		return nil, fmt.Errorf("sms code %d not found", code.Code)
	}
	return &code, nil
}
