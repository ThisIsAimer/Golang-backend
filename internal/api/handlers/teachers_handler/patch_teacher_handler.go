package teachers

import (
	"encoding/json"
	"net/http"
	"strconv"

	teacherdb "simpleapi/internal/repositories/sql/teachersdb"
)

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idstr := r.PathValue("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var updates map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "error parsing json body", http.StatusBadRequest)
		return
	}

	err = teacherdb.PatchTeacherDBHandler(w, id, updates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// used for multi update
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]any

	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "invalid request payload json not in format", http.StatusBadRequest)
		return
	}
	err = teacherdb.PatchTeachersDBHandler(w, updates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
