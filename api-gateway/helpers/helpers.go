package helpers

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)

	if err != nil {
		http.Error(w, "Cannot read body", 440)
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, data any, status int, header ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(header) > 0 {
		for key, value := range header[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)

	if err != nil {
		http.Error(w, "Cannot write json", 440)
		return err
	}

	return nil
}
