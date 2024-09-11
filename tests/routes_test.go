package tests

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
	"my-s3-clone/router"
	"my-s3-clone/dto"
	"my-s3-clone/storage"
)

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	storage.Storage
}

// Mock implementation of the DeleteObject method
func (m *MockStorage) DeleteObject(bucketName, objectKey string) error {
	return nil // Simulate successful deletion
}

// Test for the /probe-bsign{suffix:.*} route
func TestProbeBSignRoute(t *testing.T) {
	r := router.SetupRouter()

	tests := []struct {
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{"GET", "/probe-bsign", http.StatusOK, "<Response></Response>"},
		{"HEAD", "/probe-bsign", http.StatusOK, ""},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != tt.expectedCode {
			t.Errorf("expected status %d but got %d", tt.expectedCode, rr.Code)
		}

		if tt.method == "GET" && rr.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q but got %q", tt.expectedBody, rr.Body.String())
		}
	}
}

// Test for the /{bucketName}/?delete= (POST batch delete)
func TestHandleDeleteBatch(t *testing.T) {
	// Mock storage
	mockStorage := &MockStorage{}

	// Initialize the router with mock storage
	r := router.SetupRouterWithStorage(mockStorage) // Assuming you have a way to inject storage into the router

	// Create a DeleteBatchRequest with objects to delete
	deleteReq := dto.DeleteBatchRequest{
		Objects: []dto.ObjectToDelete{
			{Key: "object1.txt"},
			{Key: "object2.txt"},
		},
	}

	// Convert the request to XML
	body, err := xml.Marshal(deleteReq)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}

	// Create the POST request with the XML body
	req, err := http.NewRequest("POST", "/bucketName/?delete=", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the router with the simulated request
	r.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response content
	var deleteResult dto.DeleteResult
	err = xml.Unmarshal(rr.Body.Bytes(), &deleteResult)
	if err != nil {
		t.Fatalf("Error unmarshaling response body: %v", err)
	}

	// Check that the deleted objects match those in the request
	expectedDeleted := []string{"object1.txt", "object2.txt"}
	if len(deleteResult.DeletedResult) != len(expectedDeleted) {
		t.Errorf("expected %d deleted objects, got %d", len(expectedDeleted), len(deleteResult.DeletedResult))
	}

	for i, obj := range deleteResult.DeletedResult {
		if obj.Key != expectedDeleted[i] {
			t.Errorf("expected deleted object %s, got %s", expectedDeleted[i], obj.Key)
		}
	}
}
