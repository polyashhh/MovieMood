package main
import(
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
var movies []Movie

func loadMovies() error {
	data, err := os.ReadFile("movies.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &movies)
	if err != nil {
		return err
	}
	return nil
}

func saveMovies() error {
	data, err := json.MarshalIndent(movies, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile("movies.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func addMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(
			w,
			"Ошибка при чтении данных",
			http.StatusBadRequest,
		)
		return
	}
	movie.ID = len(movies) + 1
	movies = append(movies, movie)
	err = saveMovies()
	if err != nil {
		http.Error(
			w,
			"Ошибка при сохранении фильма",
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func moviesHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		getMovies(w, r)
		return
	}
	if r.Method == http.MethodPost {
		addMovie(w, r)
		return
	}
	http.Error(
		w,
		"Метод не поддерживается",
		http.StatusMethodNotAllowed,
	)
}

func main() {
	err := loadMovies()
	if err != nil {
		fmt.Println("Ошибка загрузки каталога:", err)
		return
	}
	http.HandleFunc("/movies", moviesHandler)
	fmt.Println(
		"Catalog service запущен: http://localhost:8081",
	)
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}