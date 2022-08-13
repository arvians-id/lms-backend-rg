package integration

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rg-km/final-project-engineering-12/backend/config"
	"github.com/rg-km/final-project-engineering-12/backend/model"
	"github.com/rg-km/final-project-engineering-12/backend/test/setup"
)

var _ = Describe("Answer API", func() {

	var (
		server         *gin.Engine
		token          string
		ok             bool
		answer         model.CreateAnswerRequest
		createquestion model.CreateQuestionRequest
	)

	BeforeEach(func() {
		configuration := config.New("../../.env.test")

		_, err := setup.SuiteSetup(configuration)
		if err != nil {
			panic(err)
		}

		router := setup.ModuleSetup(configuration)
		server = router

		answer = model.CreateAnswerRequest{
			QuestionId:  1,
			UserId:      1,
			Description: "test",
		}

		var user = model.UserRegisterResponse{
			Name:           "akuntest",
			Username:       "akuntest",
			Email:          "akuntest@gmail.com",
			Password:       "123456ll",
			Role:           2,
			Phone:          "085156789011",
			Gender:         1,
			DisabilityType: 1,
			Birthdate:      "2002-04-01",
		}

		login := model.GetUserLogin{
			Email:    "akuntest@gmail.com",
			Password: "123456ll",
		}

		// Register User
		userData, _ := json.Marshal(user)
		requestBody := strings.NewReader(string(userData))
		request := httptest.NewRequest(http.MethodPost, "/api/users", requestBody)
		request.Header.Add("Content-Type", "application/json")

		writer := httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		//Login User
		userData, _ = json.Marshal(login)
		requestBody = strings.NewReader(string(userData))
		request = httptest.NewRequest(http.MethodPost, "/api/users/login", requestBody)
		request.Header.Add("Content-Type", "application/json")

		writer = httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		responseLogin := writer.Result()

		body, _ := io.ReadAll(responseLogin.Body)
		var responseBodyLogin map[string]interface{}
		_ = json.Unmarshal(body, &responseBodyLogin)

		log.Println(responseBodyLogin["status"])
		token, ok = responseBodyLogin["token"].(string)
		if !ok {
			panic("Can't get token")
		} else {
			log.Println("Token: ", token)
		}
		createquestion = model.CreateQuestionRequest{
			UserId:      1,
			Title:       "Algoritma Naive Bayes",
			CourseId:    1,
			Tags:        "#NaiveBayes",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		}
		questionData, _ := json.Marshal(createquestion)
		requestBody = strings.NewReader(string(questionData))
		request = httptest.NewRequest(http.MethodPost, "/api/questions/create", requestBody)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Set("Authorization", token)

		writer = httptest.NewRecorder()
		server.ServeHTTP(writer, request)

		response := writer.Result()

		body, _ = io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		_ = json.Unmarshal(body, &responseBody)

		Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
		Expect(responseBody["status"]).To(Equal("question successfully created"))
		Expect(responseBody["data"].(map[string]interface{})["user_id"]).To(Equal(float64(1)))
		Expect(responseBody["data"].(map[string]interface{})["course_id"]).To(Equal(float64(1)))
		Expect(responseBody["data"].(map[string]interface{})["title"]).To(Equal("Algoritma Naive Bayes"))
		Expect(responseBody["data"].(map[string]interface{})["tags"]).To(Equal("#NaiveBayes"))
		Expect(responseBody["data"].(map[string]interface{})["description"]).To(Equal("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."))
	})

	AfterEach(func() {
		configuration := config.New("../../.env.test")
		db, err := setup.SuiteSetup(configuration)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = setup.TearDownTest(db)
		if err != nil {
			panic(err)
		}
	})
	Describe("Get Answer", func() {
		When("Answer is empty and question is empty", func() {
			It("should return data", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/answers/all", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Create Answer", func() {
		When("Answer is empty and question is empty", func() {
			It("should return error", func() {
				answerData, _ := json.Marshal(answer)
				requestBody := strings.NewReader(string(answerData))
				request := httptest.NewRequest(http.MethodPost, "/api/answers/create", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["data"]).Should(BeNil())
			})
		})
	})

	Describe("Update Answer", func() {
		When("Answer is empty and question is empty", func() {
			It("should return error", func() {
				var answerupdate = model.UpdateAnswerRequest{
					QuestionId:  1,
					UserId:      1,
					Description: "kentang jalan",
				}

				answerData, _ := json.Marshal(answerupdate)
				requestBody := strings.NewReader(string(answerData))
				request := httptest.NewRequest(http.MethodPut, "/api/answers/update/1", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["data"]).Should(BeNil())
			})
		})
	})
	Describe("Delete Answer", func() {
		When("Answer is empty and question is empty", func() {
			It("should return error", func() {
				request := httptest.NewRequest(http.MethodDelete, "/api/answers/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusInternalServerError))
				Expect(responseBody["data"]).Should(BeNil())
			})
		})
	})
	Describe("GetAnswers", func() {
		When("Answer is empty and question is empty", func() {
			It("should return data", func() {
				request := httptest.NewRequest(http.MethodGet, "/api/answers/by-user/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})
})
