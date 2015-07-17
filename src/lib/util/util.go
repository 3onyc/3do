package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, v interface{}) {
	if enc, err := json.Marshal(v); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
	} else {
		w.Write(enc)
	}
}
