const findButton = document.getElementById("findButton");
const randomButton = document.getElementById("randomButton");
const message = document.getElementById("message");
const resultsSection = document.getElementById("resultsSection");
const resultsContainer = document.getElementById("resultsContainer");

function getParameters(){
    return {
        type: document.getElementById("type").value,
        genre: document.getElementById("genre").value,
        mood: document.getElementById("mood").value,
        max_duration: Number(
            document.getElementById("duration").value
        ),
        company: document.getElementById("company").value
    };
}

async function saveParameters(){
    const response = await fetch(
        "http://localhost:8082/parameters",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(getParameters())
        }
    );
    return response.ok;
}

function createMovieCard(movie) {
    return `
        <div class="result-card">
            <h3>${movie.title}</h3>
            <p><strong>Тип:</strong> ${movie.type}</p>
            <p><strong>Жанр:</strong> ${movie.genre}</p>
            <p><strong>Настроение:</strong> ${movie.mood}</p>
            <p><strong>Длительность:</strong> ${movie.duration} мин.</p>
            <p><strong>Просмотр:</strong> ${movie.company}</p>
            <p><strong>Возраст:</strong> ${movie.age_rating}+</p>
        </div>
    `;
}

findButton.addEventListener("click", async function () {
    message.textContent = "Ищем подходящие варианты...";
    resultsSection.classList.add("hidden");
    resultsContainer.innerHTML = "";
    try {
        const saved = await saveParameters();
        if (!saved) {
            message.textContent =
                "Не удалось сохранить параметры";
            return;
        }
        const response = await fetch(
            "http://localhost:8084/results"
        );
        if (!response.ok) {
            message.textContent =
                "Не удалось получить результаты";
            return;
        }
        const result = await response.json();
        message.textContent = result.message;
        resultsSection.classList.remove("hidden");
        if (!result.movies || result.movies.length === 0) {
            resultsContainer.innerHTML =
                "<p>Попробуйте изменить параметры поиска</p>";
            return;
        }
        result.movies.forEach(function (movie) {
            resultsContainer.innerHTML +=
                createMovieCard(movie);
        });
    } catch (error) {
        message.textContent =
            "Не удалось подключиться к сервисам MovieMood";
    }
});

randomButton.addEventListener("click", async function () {
    message.textContent = "Выбираем случайный вариант...";
    resultsSection.classList.add("hidden");
    resultsContainer.innerHTML = "";
    try {
        const saved = await saveParameters();
        if (!saved){
            message.textContent =
                "Не удалось сохранить параметры";
            return;
        }
        const response = await fetch(
            "http://localhost:8083/recommendations/random"
        );
        if (response.status === 404) {
            message.textContent =
                "Подходящие варианты не найдены";
            return;
        }
        if (!response.ok) {
            message.textContent =
                "Не удалось получить рекомендацию";
            return;
        }
        const movie = await response.json();
        message.textContent =
            "Вот что MovieMood выбрал для тебя:";
        resultsSection.classList.remove("hidden");
        resultsContainer.innerHTML =
            createMovieCard(movie);
    } catch (error) {
        message.textContent =
            "Не удалось подключиться к сервисам MovieMood";
    }
});