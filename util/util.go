package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func WriteJSONResponse(w http.ResponseWriter, v interface{}) {
	if enc, err := json.Marshal(v); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
	} else {
		w.Write(enc)
	}
}

func ShowNewLines(i string) string {
	return strings.Replace(i, "\n", "\\n\n", -1)
}

func TrimRightSpace(i string) string {
	return strings.TrimRightFunc(i, unicode.IsSpace)
}
