package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/bmasaiti/go-projects/my-vault/internal/domain"
)

const requestJSON = `{
  "Id": 34,
  "name": "server-secret",
  "kv_map": {
    "username": "Tamuka",
    "password": "pass123"
  }
}`

const getResponseJSON = `{
  "message": "Retrieved secret object",
  "id": "01970123-643a-7d9c-abbf-684a3abe129c",
  "name": "server-secret",
  "kv_map": {
    "password": "pass123",
    "username": "Tamuka"
  }
}`

const responseJSON = `{
  "message": "Successfully created secret with secretId: 019700cc-d778-7b08-8660-a8965842faf3 and name: server-secret",
  "id": "019700cc-d778-7b08-8660-a8965842faf3",
  "name": "server-secret"
}`

type mockSecretRepo struct{}

func (m *mockSecretRepo) ListAllSecrets() ([]domain.Secret, error) {
	return []domain.Secret{
		{
			Id:   "019705dc-8e40-7ead-b13d-bd41c3f7f476",
			Name: "server-secret",
			KVMap: map[string]string{
				"username": "Tamuka",
				"password": "pass123",
			},
		},
	}, nil
}
// secrets := map[string]Secret{
// 		"01970123-643a-7d9c-abbf-684a3abe129c": {
// 			Id:   "01970123-643a-7d9c-abbf-684a3abe129c",
// 			Name: "server-secret-1",
// 			KVMap: map[string]string{
// 				"username": "Tamuka",
// 				"password": "pass123",
// 			},
// 		},
// 		"01970123-743b-8e0d-bbcc-12345ef67890": {
// 			Id:   "01970123-743b-8e0d-bbcc-12345ef67890",
// 			Name: "server-secret-2",
// 			KVMap: map[string]string{
// 				"api_key": "abcd1234",
// 				"env":     "production",
// 			},
// 		},
// 	}

// func TestHandleGetSecrets(t *testing.T){

// 	t.Run("Returns a secret object", func(t *testing.T){
// 		request := httptest.NewRequest(http.MethodGet,"/v1/secrets/01970123-643a-7d9c-abbf-684a3abe129c", nil)
// 		//secretId:= strings.TrimPrefix(req.URL.Path , "v1/secrets/secretId")
// 		response := httptest.NewRecorder()

// 		HandleGetSecretById(response,request)

// 		got := response.Body.String()

// 		want:= getResponseJSON

// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 		})

// }

// func TestHandleListAllSecrets(t *testing.T){

// 	t.Run("Returns a secret object", func(t *testing.T){
// 		request := httptest.NewRequest(http.MethodGet,"/v1/secrets/", nil)
// 		response := httptest.NewRecorder()

// 		HandleListSecrets(response,request)

// 		got := response.Body.String()

// 		want:= getResponseJSON // return a list here of secrets.

// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 		})

// }

func TestHandleListAllSecrets(t *testing.T) {
	handler := &SecretHandler{
		DB: &mockSecretRepo{},
	}
	t.Run("Returns a secret object", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/v1/secrets/", nil)
		response := httptest.NewRecorder()

		HandleListSecrets(response, request)

		body := response.Body.Bytes()
		log.Default().Fatal(string(body))

		var got []domain.Secret

		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		want := []domain.Secret{
			{
				Id:   "019705dc-8e40-7ead-b13d-bd41c3f7f476",
				Name: "server-secret",
				KVMap: map[string]string{
					"username": "Tamuka",
					"password": "pass123",
				},
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %+v, got %+v", want, got)
		}
	})
}

// func TestHandlePostSecret(t *testing.T) {

// 	const requestJSON = `"id": 34,
// 		"name": "server-secret",
// 		"kv_map": {
// 			"username": "Tamuka",
// 			"password": "pass123"
// 		}

// 	}`

// 	method := http.MethodPost
// 	endpoint := "/v1/secrets"

// 	request := httptest.NewRequest(method, endpoint, strings.NewReader(requestJSON))
// 	response := httptest.NewRecorder()
// 	defer request.Body.Close()

// 	HandlePostSecret(response, request)

// 	resp := response.Result()
// 	body, err := io.ReadAll(resp.Body)
// 	log.Println("JSON decode error:", err)
// 	log.Println("Body:", string(body))

// 	if err != nil {
// 		log.Println("JSON decode error:", err)
// 		t.Fatal(err)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("want: %03d, got: %03d", http.StatusOK, resp.StatusCode)
// 	}
// 	want := fmt.Sprintf(`[%s] %s OK!`, method, endpoint)
// 	got := string(body)

// 	if want != got {
// 		t.Fatalf("want: %+q, got: %+q", want, got)
// 	}
// }

// func TestHandleGetSecretById(t *testing.T) {
// 	mockRepo := &MockSecretsRepository{
// 		GetSecretByIdFunc: func(id string) (Secret, error) {
// 			if id == "123" {
// 				return Secret{
// 					Id:   "123",
// 					Name: "TestSecret",
// 					KVMap: map[string]string{
// 						"key1": "value1",
// 					},
// 				}, nil
// 			}
// 			return Secret{}, errors.New("not found")
// 		},
// 	}

// 	handler := &SecretHandler{Repo: mockRepo}

// 	req := httptest.NewRequest("GET", "/v1/secrets/123", nil)
// 	// Simulate router param: "secret_id" = "123"
// 	req = mux.SetURLVars(req, map[string]string{"secret_id": "123"})

// 	rr := httptest.NewRecorder()

// 	handler.HandleGetSecretById(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	expected := `{"Id":"123","Name":"TestSecret","KVMap":{"key1":"value1"}}`
// 	assert.JSONEq(t, expected, rr.Body.String())
// }

func TestHandlePostSecret(t *testing.T) {
	const requestJSON = `{
		"id": 34,
		"name": "server-secret",
		"kv_map": {
			"username": "Tamuka",
			"password": "pass123"
		}
	}`

	// respData := struct {
	// 	Message string `json:"message"`
	// 	ID      string `json:"id"`
	// 	Name    string `json:"name"`
	// }
	method := http.MethodPost
	endpoint := "/v1/secrets"

	request := httptest.NewRequest(method, endpoint, strings.NewReader(requestJSON))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	HandlePostSecret(response, request)

	resp := response.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Failed to read response body:", err)
	}

	log.Println("Body:", string(body))

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}

	want := fmt.Sprintf(`[%s] %s OK!`, method, endpoint)
	got := string(body)

	if want != got {
		t.Fatalf("Expected body %q, got %q", want, got)
	}
}
