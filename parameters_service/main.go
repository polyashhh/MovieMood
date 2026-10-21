package main
import(
	"encoding/json"
	"fmt"
	"net/http"
)
type Parameters struct{
	Type        string `json:"type"`
	Genre       string `json:"genre"`
	Mood        string `json:"mood"`
	MaxDuration int    `json:"max_duration"`
	Company     string `json:"company"`
}
var currentParameters Parameters

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, POST, OPTIONS",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Content-Type",
	)
}

func saveParameters(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	var parameters Parameters
	err := json.NewDecoder(r.Body).Decode(&parameters)
	if err != nil {
		http.Error(
			w,
			"Ошибка при чтении параметров",
			http.StatusBadRequest,
		)
		return
	}
	currentParameters = parameters
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(currentParameters)
}

func getParameters(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentParameters)
}

func parametersHandler(
	w http.ResponseWriter,
	r *http.Request,
){
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {
		getParameters(w, r)
		return
	}
	if r.Method == http.MethodPost {
		saveParameters(w, r)
		return
	}
	http.Error(
		w,
		"Метод не поддерживается",
		http.StatusMethodNotAllowed,
	)
}

func main() {
	http.HandleFunc(
		"/parameters",
		parametersHandler,
	)
	fmt.Println(
		"Parameters service запущен: http://localhost:8082",
	)
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println(
			"Ошибка запуска сервера:",
			err,
		)
	}
}