package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
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

// Mock implementation of SecretsRepository
type MockSecretsRepository struct {
	secrets []domain.Secret
	err     error
}

func (m *MockSecretsRepository) PutNewSecret(secret domain.Secret) error {
	return m.err
}

func (m *MockSecretsRepository) GetScretsById(Id string) (domain.Secret, error) {

	return m.secrets[0],m.err  //cheating
	
}

func (m *MockSecretsRepository) DeleteSecretByID(Id string) (string, error) {
	return "", m.err
}

func (m *MockSecretsRepository) ListAllSecrets() ([]domain.Secret, error) {
	return m.secrets, m.err
}

func TestHandlePostSecret(t *testing.T) {
	t.Run("returns 200 OK with JSON on success", func(t *testing.T) {
		mockRepo := &MockSecretsRepository{
			secrets: []domain.Secret{
				{Id: "1", Name: "test-secret", KVMap: map[string]string{"key": "value"}},
			},
			err: nil,
		}
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("POST", "/secrets", bytes.NewBuffer([]byte(`{"name": "test-secret", "key": "value"}`)))
		w := httptest.NewRecorder()

		handler.HandlePostSecret(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}
	})

	t.Run("returns 500 on database error", func(t *testing.T) {
		mockRepo := &MockSecretsRepository{
			err: errors.New("database connection failed"),
		}
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("POST", "/secrets", bytes.NewBuffer([]byte(`{"name": "test-secret", "key": "value"}`)))
		w := httptest.NewRecorder()

		handler.HandlePostSecret(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestHandleGetSecretById(t *testing.T) {
		mockRepo := &MockSecretsRepository{
				secrets: []domain.Secret{
					{
						Id:   "01970123-643a-7d9c-abbf-684a3abe129c",  
						Name: "server-secret",                        
						KVMap: map[string]string{                      
							"password": "pass123",                     
							"username": "Tamuka",                     
						},
					},
				},
				err: nil,
			}
	t.Run("returns 200 OK with json secret object", func(t *testing.T) {
		
		
	
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("GET", "/secrets/01970123-643a-7d9c-abbf-684a3abe129c", nil)
		w := httptest.NewRecorder()

		handler.HandleGetSecretById(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		 if !strings.Contains(w.Body.String(), "01970123-643a-7d9c-abbf-684a3abe129c") {
			t.Errorf("Expected secret ID in response body,got %s", w.Body.String())
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}
	})
	t.Run("returns 500 on database error", func(t *testing.T) {
		// mockRepo := &MockSecretsRepository{
		// 	err: errors.New("database connection failed"),
		// }
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("GET", "/secrets/01970123-643a-7d9c-abbf-684a3abe129c", nil)
		w := httptest.NewRecorder()

		handler.HandleGetSecretById(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}

		// if !strings.Contains(w.Body.String(), "Unexpected internal server error") {
		// 	t.Errorf("expected error message in response body")
		// }
	})

}
func TestHandleDeleteSecretById(t *testing.T) {
	t.Run("returns 200 OK with JSON on success", func(t *testing.T) {
		mockRepo := &MockSecretsRepository{
			secrets: []domain.Secret{
				{Id: "1", Name: "test-secret", KVMap: map[string]string{"key": "value"}},
			},
			err: nil,
		}
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("DELETE", "/secrets/1", nil)
		w := httptest.NewRecorder()

		handler.HandleDeleteSecretById(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
			t.Errorf("expected Content-Type text/plain, got %s", contentType)
		}
	})

	t.Run("returns 500 on database error", func(t *testing.T) {
		mockRepo := &MockSecretsRepository{
			err: errors.New("database connection failed"),
		}
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("DELETE", "/secrets/1", nil)
		w := httptest.NewRecorder()

		handler.HandleDeleteSecretById(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}

		// if !strings.Contains(w.Body.String(), "Unexpected internal server error") {
		// 	t.Errorf("expected error message in response body")
		// }
	})
}

func TestHandleListSecretsById(t *testing.T) {
	t.Run("returns 200 OK with JSON on success", func(t *testing.T) {
		mockRepo := &MockSecretsRepository{
			secrets: []domain.Secret{
				{Id: "1", Name: "test-secret", KVMap: map[string]string{"key": "value"}},
			},
			err: nil,
		}
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("GET", "/secrets", nil)
		w := httptest.NewRecorder()

		handler.HandleListSecrets(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}
	})

	t.Run("returns 500 on database error", func(t *testing.T) {
		mockRepo := &MockSecretsRepository{
			err: errors.New("database connection failed"),
		}
		handler := &SecretHandler{DB: mockRepo}

		req := httptest.NewRequest("GET", "/secrets", nil)
		w := httptest.NewRecorder()

		handler.HandleListSecrets(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), "Unexpected internal server error") {
			t.Errorf("expected error message in response body")
		}
	})
}
