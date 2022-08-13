package integration

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rg-km/final-project-engineering-12/backend/config"
	"github.com/rg-km/final-project-engineering-12/backend/model"
	"github.com/rg-km/final-project-engineering-12/backend/test/setup"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("User Course API", func() {

	var (
		server *gin.Engine
		token  string
		ok     bool
		idUser int
	)

	BeforeEach(func() {
		configuration := config.New("../../.env.test")

		_, err := setup.SuiteSetup(configuration)
		if err != nil {
			panic(err)
		}

		router := setup.ModuleSetup(configuration)
		server = router

		var user = model.UserRegisterResponse{
			Name:           "akuntest",
			Username:       "akuntest",
			Email:          "akuntest@gmail.com",
			Password:       "123456ll",
			Role:           1,
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

		responseRegister := writer.Result()

		bodyRegister, _ := io.ReadAll(responseRegister.Body)
		var responseBodyRegister map[string]interface{}
		_ = json.Unmarshal(bodyRegister, &responseBodyRegister)

		idUser = int(responseBodyRegister["data"].(map[string]interface{})["id"].(float64))

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

	Describe("Create User Course", func() {
		When("Data is empty", func() {
			It("should return Data User Course", func() {
				// Create User Course
				requestBody := strings.NewReader(`{"user_id": 10,"course_id": 1}`)
				request := httptest.NewRequest(http.MethodPost, "/api/usercourse", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				Body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(Body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusCreated))
				Expect(responseBody["status"]).To(Equal("User Course Create Succesfully"))
				Expect(responseBody["data"].(map[string]interface{})["user_id"]).To(Equal(float64(10)))
				Expect(responseBody["data"].(map[string]interface{})["course_id"]).To(Equal(float64(1)))
			})
		})
	})

	Describe("Get User Course By UserId and CourseId", func() {
		When("Data is empty", func() {
			It("should return User Course", func() {
				// Create User Course 1
				requestBody := strings.NewReader(`{"user_id": 10,"course_id": 1}`)
				request := httptest.NewRequest(http.MethodPost, "/api/usercourse", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)
				// Get User Course By UserId and CourseId
				request = httptest.NewRequest(http.MethodGet, "/api/usercourse/10/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				Body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(Body, &responseBody)

				log.Println(responseBody["status"])
				usercourse := responseBody["data"].(map[string]interface{})

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody["status"]).To(Equal("Get User Course Successfull"))
				Expect(usercourse["user_id"]).To(Equal(float64(10)))
				Expect(usercourse["course_id"]).To(Equal(float64(1)))
			})
		})
	})

	Describe("Delete User Course By UserId and CourseId", func() {
		When("Data is empty", func() {
			It("should return User Course", func() {
				// Create User Course 1
				requestBody := strings.NewReader(`{"user_id": 10,"course_id": 1}`)
				request := httptest.NewRequest(http.MethodPost, "/api/usercourse", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)
				// Delete User Course By UserId and CourseId
				request = httptest.NewRequest(http.MethodDelete, "/api/usercourse/10/1", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				Body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(Body, &responseBody)

				Expect(int(responseBody["code"].(float64))).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Find All Course By User Logged In", func() {
		When("Data is exists", func() {
			It("should return all course when user is logged in", func() {
				// Create Course
				requestBody := strings.NewReader(`{"name": "Teknik Komputer Jaringan","class": "TKJ-3","tools": "Router, RJ-45","about": "Pada pelajaran kali ini akan lebih difokuskan pada pembuatan tower","description": "Siswa mampu membuat tower sendiri"}`)
				request := httptest.NewRequest(http.MethodPost, "/api/courses", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer := httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response := writer.Result()

				Body, _ := io.ReadAll(response.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(Body, &responseBody)

				// Create User Course
				courseId := int(responseBody["data"].(map[string]interface{})["id"].(float64))
				bodyCourse := fmt.Sprintf(`{"user_id": %v,"course_id": %v}`, idUser, courseId)
				requestBody = strings.NewReader(bodyCourse)
				request = httptest.NewRequest(http.MethodPost, "/api/usercourse", requestBody)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				// Find All Course By User Logged In
				request = httptest.NewRequest(http.MethodGet, "/api/usercourse/courses", nil)
				request.Header.Add("Content-Type", "application/json")
				request.Header.Set("Authorization", token)

				writer = httptest.NewRecorder()
				server.ServeHTTP(writer, request)

				response = writer.Result()

				Body, _ = io.ReadAll(response.Body)
				var responseBody1 map[string]interface{}
				_ = json.Unmarshal(Body, &responseBody1)

				items := responseBody1["data"].([]interface{})
				itemResponse := items[0].(map[string]interface{})

				Expect(int(responseBody1["code"].(float64))).To(Equal(http.StatusOK))
				Expect(responseBody1["status"]).To(Equal("Get All User Course Successfull"))
				Expect(int(itemResponse["id_course"].(float64))).To(Equal(int(responseBody["data"].(map[string]interface{})["id"].(float64))))
				Expect(itemResponse["course_name"]).To(Equal(responseBody["data"].(map[string]interface{})["name"]))
				Expect(itemResponse["course_code"]).To(Equal(responseBody["data"].(map[string]interface{})["code_course"]))
				Expect(itemResponse["course_class"]).To(Equal(responseBody["data"].(map[string]interface{})["class"]))
			})
		})
	})
})
