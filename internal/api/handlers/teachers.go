package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"simpleapi/internal/models"
	"simpleapi/internal/repositories/sqlconnect"
)

// for now

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

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
	db_name := os.Getenv("DB_NAME")

	db, err := sqlconnect.ConnectDB(db_name)
	if err != nil {
		http.Error(w, "error connecting to server", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email, class, subject) VALUES(?, ?, ?, ?, ?)")

	if err != nil {
		http.Error(w,"error in praparing sql query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for i, teacher := range newTeachers{
		res, err := stmt.Exec(teacher.FirstName,teacher.LastName,teacher.Email,teacher.Class,teacher.Subject)
		if err != nil {
			http.Error(w, "error inserting values in the database", http.StatusInternalServerError)
			return
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "error getting last inserted id", http.StatusInternalServerError)
			return
		}
		newTeachers[i].ID = int(lastId)
		
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
