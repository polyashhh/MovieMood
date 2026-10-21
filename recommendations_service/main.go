package main
import(
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
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
type Parameters struct{
	Type        string `json:"type"`
	Genre       string `json:"genre"`
	Mood        string `json:"mood"`
	MaxDuration int    `json:"max_duration"`
	Company     string `json:"company"`
}
type MoviesResult struct{
	Movies []Movie
	Err    error
}
type ParametersResult struct{
	Parameters Parameters
	Err        error
}

func getMovies() ([]Movie, error) {
	response, err := http.Get("http://localhost:8081/movies")
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

func getParameters() (Parameters, error) {
	response, err := http.Get("http://localhost:8082/parameters")
	if err != nil {
		return Parameters{}, err
	}
	defer response.Body.Close()
	var parameters Parameters
	err = json.NewDecoder(response.Body).Decode(&parameters)
	if err != nil{
		return Parameters{}, err
	}
	return parameters, nil
}

func getData()([]Movie, Parameters, error) {
	moviesChannel := make(chan MoviesResult)
	parametersChannel := make(chan ParametersResult)
	go func(){
		movies, err := getMovies()
		moviesChannel <- MoviesResult{
			Movies: movies,
			Err:    err,
		}
	}()
	go func(){
		parameters, err := getParameters()
		parametersChannel <- ParametersResult{
			Parameters: parameters,
			Err:        err,
		}
	}()
	moviesResult := <-moviesChannel
	parametersResult := <-parametersChannel
	if moviesResult.Err != nil {
		return nil, Parameters{}, moviesResult.Err
	}
	if parametersResult.Err != nil {
		return nil, Parameters{}, parametersResult.Err
	}
	return moviesResult.Movies, parametersResult.Parameters, nil
}

func filterMovies(movies []Movie, parameters Parameters)[]Movie {
	recommendations := make([]Movie, 0)
	for _, movie := range movies {
		if parameters.Type != "" && movie.Type != parameters.Type {
			continue
		}
		if parameters.Genre != "" && movie.Genre != parameters.Genre {
			continue
		}
		if parameters.Mood != "" && movie.Mood != parameters.Mood {
			continue
		}
		if parameters.MaxDuration > 0 &&
			movie.Duration > parameters.MaxDuration {
			continue
		}
		if parameters.Company != "" &&
			movie.Company != parameters.Company {
			continue
		}
		recommendations = append(recommendations, movie)
	}
	return recommendations
}

func getRecommendations(w http.ResponseWriter, r *http.Request){
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	movies, parameters, err := getData()
	if err != nil {
		http.Error(
			w,
			"Не удалось получить данные",
			http.StatusInternalServerError,
		)
		return
	}
	recommendations := filterMovies(movies, parameters)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Duration <
			recommendations[j].Duration
	})
	json.NewEncoder(w).Encode(recommendations)
}

func getRandomRecommendation(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	movies, parameters, err := getData()
	if err != nil {
		http.Error(
			w,
			"Не удалось получить данные",
			http.StatusInternalServerError,
		)
		return
	}
	recommendations := filterMovies(movies, parameters)
	if len(recommendations) == 0 {
		http.Error(
			w,
			"Подходящие варианты не найдены",
			http.StatusNotFound,
		)
		return
	}
	randomIndex := rand.Intn(len(recommendations))
	json.NewEncoder(w).Encode(
		recommendations[randomIndex],
	)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main(){
	http.HandleFunc(
		"/recommendations",
		getRecommendations,
	)
	http.HandleFunc(
		"/recommendations/random",
		getRandomRecommendation,
	)
	fmt.Println(
		"Recommendations service запущен: http://localhost:8083",
	)
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}