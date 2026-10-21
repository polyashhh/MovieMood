package main
import(
	"encoding/json"
	"fmt"
	"net/http"
)
type Movie struct{
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Genre     string `json:"genre"`
	Mood      string `json:"mood"`
	Duration  int    `json:"duration"`
	Company   string `json:"company"`
	AgeRating int    `json:"age_rating"`
}
type Result struct{
	Count   int     `json:"count"`
	Movies  []Movie `json:"movies"`
	Message string  `json:"message"`
}

func enableCORS(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getRecommendations() ([]Movie, error) {
	response, err := http.Get(
		"http://localhost:8083/recommendations",
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var movies []Movie
	err = json.NewDecoder(response.Body).Decode(&movies)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func getResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movies, err := getRecommendations()
	if err != nil{
		http.Error(
			w,
			"Не удалось получить рекомендации",
			http.StatusInternalServerError,
		)
		return
	}
	result := Result{
		Count:  len(movies),
		Movies: movies,
	}
	if len(movies) == 0 {
		result.Message = "Подходящие варианты не найдены"
	} else {
		result.Message = "Подходящие варианты найдены"
	}
	json.NewEncoder(w).Encode(result)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {
		getResults(w, r)
		return
	}
	http.Error(
		w,
		"Метод не поддерживается",
		http.StatusMethodNotAllowed,
	)
}

func main(){
	http.HandleFunc("/results",resultsHandler)
	fmt.Println(
		"Results service запущен: http://localhost:8084",
	)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}