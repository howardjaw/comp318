package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type toDo struct {
	Id          int
	Description string
}

type toDoList struct {
	mtx    sync.RWMutex
	list   map[int]string
	nextID int
}

/**
Create HTTP Handler
**/

func New() http.Handler {
	todos := &toDoList{
		list:   make(map[int]string),
		nextID: 1,
	}

	mux := http.NewServeMux()

	// Collection routes
	mux.HandleFunc("GET /todos", todos.getAll)
	mux.HandleFunc("POST /todos", todos.createToDo)
	mux.HandleFunc("OPTIONS /todos", todos.optionsAllToDos)

	// Individual resources
	mux.HandleFunc("GET /todos/", todos.getToDo)
	mux.HandleFunc("DELETE /todos/", todos.deleteToDo)

	// Long running example using context
	mux.HandleFunc("GET /slow", todos.slowHandler)

	return mux
}

func (t *toDoList) getAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	t.mtx.RLock()

	todos := make([]toDo, 0, len(t.list))

	for id, description := range t.list {
		todos = append(todos, toDo{
			Id:          id,
			Description: description,
		})
	}

	t.mtx.RUnlock()

	encoded, err := json.Marshal(todos)

	if err != nil {
		http.Error(
			w,
			"could not encode todos",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func (t *toDoList) createToDo(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(
			w,
			"unable to read body",
			http.StatusInternalServerError,
		)
		return
	}

	description := string(body)

	if description == "" {
		http.Error(
			w,
			"description cannot be empty",
			http.StatusBadRequest,
		)
		return
	}

	t.mtx.Lock()

	id := t.nextID
	t.nextID++

	t.list[id] = description

	t.mtx.Unlock()

	newTodo := toDo{
		Id:          id,
		Description: description,
	}

	encoded, err := json.Marshal(newTodo)

	if err != nil {
		http.Error(
			w,
			"could not encode todo",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(encoded)
}

func (t *toDoList) getToDo(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id, err := getIDFromPath(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"invalid todo id",
			http.StatusBadRequest,
		)
		return
	}

	t.mtx.RLock()
	description, exists := t.list[id]
	t.mtx.RUnlock()

	if !exists {
		http.Error(
			w,
			"todo not found",
			http.StatusNotFound,
		)
		return
	}

	todo := toDo{
		Id:          id,
		Description: description,
	}

	json.NewEncoder(w).Encode(todo)
}

func (t *toDoList) deleteToDo(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id, err := getIDFromPath(r.URL.Path)

	if err != nil {
		http.Error(
			w,
			"invalid todo id",
			http.StatusBadRequest,
		)
		return
	}

	t.mtx.Lock()

	_, exists := t.list[id]

	if !exists {
		t.mtx.Unlock()
		http.Error(
			w,
			"todo not found",
			http.StatusNotFound,
		)
		return
	}

	delete(t.list, id)

	t.mtx.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (t *toDoList) optionsAllToDos(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Allow", "GET, POST, OPTIONS")
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, POST, OPTIONS",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Content-Type",
	)

	w.WriteHeader(http.StatusOK)
}

func (t *toDoList) slowHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	timeoutCtx, cancel := context.WithTimeout(
		ctx,
		100*time.Millisecond,
	)

	defer cancel()

	result, err := slowFunction(timeoutCtx)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusRequestTimeout,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, result)
}

func slowFunction(ctx context.Context) (string, error) {
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("request canceled")
		default:
			//some arbitrary work
		}
		time.Sleep(10 * time.Millisecond) //simulating the computations some arbitrary work take
	}
	return "finished", nil
}

func getIDFromPath(path string) (int, error) {
	idString := strings.TrimPrefix(path, "/todos/")
	return strconv.Atoi(idString)
}

func main() {
	handler := New()

	var srv http.Server

	srv.Addr = "localhost:5318"
	srv.Handler = handler

	fmt.Println("server running at: http://localhost:5318")

	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}
}
