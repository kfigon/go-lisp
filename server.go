package main

import (
	"encoding/json"
	"fmt"
	"go-lisp/config"
	"go-lisp/models"
	"net/http"
)

func NewServer(code string) (*http.ServeMux, error) {
	confStore, err := config.New(code)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/flag/{name}", func(w http.ResponseWriter, r *http.Request) {
		flagName := r.PathValue("name")
		fetchData(confStore, flagName, nil, w)
	})

	mux.HandleFunc("POST /api/flag/{name}", func(w http.ResponseWriter, r *http.Request) {
		flagName := r.PathValue("name")
		type Input struct {
			Args []any
		}
		var input Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]any{"error": "error parsing input json: " + err.Error()})
			return
		}
		args := []models.SExpression{}
		for _, a := range input.Args {
			switch v := a.(type) {
			case int:
				args = append(args, models.Number(v))
			case float64:
				args = append(args, models.Number(v))
			case bool:
				args = append(args, models.Bool(v))
			case string:
				args = append(args, models.String(v))
			default:
				writeJson(w, http.StatusBadRequest, map[string]any{"error": fmt.Sprintf("invalid input type: %T", v)})
			}
		}
		fetchData(confStore, flagName, args, w)
	})

	return mux, nil
}

func fetchData(confStore *config.ConfigurationStore, flagName string, params []models.SExpression, w http.ResponseWriter) {
	v, err := confStore.Get(flagName, params...)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("failed to get %s: %v", flagName, err.Error()),
		})
		return
	}

	switch v := v.(type) {
	case models.Bool:
		writeJson(w, http.StatusOK, map[string]any{"bool_value": bool(v)})
	case models.Number:
		writeJson(w, http.StatusOK, map[string]any{"num_value": int(v)})
	case models.String:
		writeJson(w, http.StatusOK, map[string]any{"str_value": string(v)})
	default:
		writeJson(w, http.StatusBadRequest, map[string]any{
			"error": fmt.Sprintf("invalid response type %T", v),
		})
	}
}

func writeJson(w http.ResponseWriter, status int, out map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(out)
}
