package application_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator/internal/application"
	"calculator/pkg/calculation"
)

type testExpression struct {
	Expression string `json:"expression"`
}

func TestApplication_RunServer(t *testing.T) {
	tt := []struct {
		expression string
		statusCode int
		result     float64
	}{
		{"2+2*2", http.StatusOK, 6},
		{"2*(2+2)", http.StatusOK, 8},
		{"9/5+3", http.StatusOK, 4.8},
		{"7*(1/(2+8))", http.StatusOK, 0.7},
	}

	for _, tc := range tt {
		t.Run(tc.expression, func(t *testing.T) {
			data := testExpression{Expression: tc.expression}
			jsonData, _ := json.Marshal(data)
			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(application.CalcHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.statusCode)
			}

			body, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			if err != nil {
				t.Fatal(err)
			}

		})
	}

	uncorrect_tt := []struct {
		expression string
		statusCode int
		error      string
	}{
		{"grechatop", http.StatusBadRequest, calculation.ErrInvalidExpression.Error()},
		{"2+qwerty", http.StatusBadRequest, calculation.ErrInvalidExpression.Error()},
		{"", http.StatusBadRequest, calculation.ErrInvalidExpression.Error()},
		{"5/0", http.StatusBadRequest, calculation.ErrDivisionByZero.Error()},
		{"1+((12)", http.StatusBadRequest, calculation.ErrInvalidExpression.Error()},
	}

	for _, tc := range uncorrect_tt {
		t.Run(tc.expression, func(t *testing.T) {
			data := testExpression{Expression: tc.expression}
			jsonData, _ := json.Marshal(data)
			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(application.CalcHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.statusCode)
			}

			body, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			if err != nil {
				t.Fatal(err)
			}

			if responseBody["error"].(string) != tc.error {
				t.Errorf("handler returned wrong result: got %v want %v", responseBody["error"], tc.error)
			}
		})
	}
}
