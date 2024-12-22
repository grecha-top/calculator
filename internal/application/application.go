package application

import (
	"bufio"
	"calculator/pkg/calculation"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		//вычисляем выражение
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		var statusCode int
		var errorMessage string

		if errors.Is(err, calculation.ErrInvalidExpression) {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		} else if errors.Is(err, calculation.ErrDivisionByZero) {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
}

func (a *Application) RunServer() error {
	http.HandleFunc("/", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
