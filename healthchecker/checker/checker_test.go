package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealCheckSuccess(t *testing.T) {
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	defer server.Close()

	result := CheckEndpoint(server.URL, 2*time.Second)

	if !result.Healthy {
		t.Errorf("Expected healthy endpoint, got unhealthy")
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 , got %d", result.StatusCode)
	}

	if result.ResponseTime == 0 {
		t.Errorf("Expected non zero response time")
	}

}
