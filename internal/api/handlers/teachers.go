package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"simpleapi/internal/models"
)

// for now

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Rudra",
		LastName:  "Shivdev",
		Class:     "6A",
		Subject:   "math",
	}
	nextId++

	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Rudrina",
		LastName:  "Shivdev",
		Class:     "10B",
		Subject:   "computer",
	}

	nextId++

	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Tanjiro",
		LastName:  "Kamado",
		Class:     "all",
		Subject:   "Dance",
	}

	nextId++

	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Zenitsu",
		LastName:  "Agatsuma",
		Class:     "8C",
		Subject:   "Science",
	}

	nextId++

	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Inosuke",
		LastName:  "Hashibira",
		Class:     "5D",
		Subject:   "Sports",
	}

}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idstr := strings.TrimSuffix(path, "/")

	w.Header().Set("Content-Type", "application/json")

	//handle quary parametre
	if idstr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teacherList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		//handle path parametre
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		tearcher, exists := teachers[id]
		if !exists {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(tearcher)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}

}

func postTeachersHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for i, teacher := range newTeachers {
		nextId++
		newTeachers[i].ID = nextId

		teacher.ID = nextId
		teachers[nextId] = teacher
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(newTeachers),
		Data:   newTeachers,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func TeachersRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "string")
	fmt.Println("someone accessed: Teachers route")

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandler(w, r)
	case http.MethodPut:
		fmt.Fprintln(w, "accessed : Teachers. with: Put")
	case http.MethodPatch:
		fmt.Fprintln(w, "accessed : Teachers. with: Patch")
	case http.MethodDelete:
		fmt.Fprintln(w, "accessed : Teachers. with: Delete")
	default:
		fmt.Fprintln(w, "accessed : Teachers")

	}
}
