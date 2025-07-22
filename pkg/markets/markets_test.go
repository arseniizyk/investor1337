package markets

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

type mockPaginatedResponse struct {
	Items  []int
	Next   string
	Called int
}

func (m *mockPaginatedResponse) LenData() int {
	return len(m.Items)
}

func (m *mockPaginatedResponse) Cursor() string {
	return m.Next
}

type mockResponse struct {
	Data string `json:"data"`
}

func TestDoJSONRequest(t *testing.T) {
	type testCase struct {
		name         string
		handler      http.HandlerFunc
		expectedData mockResponse
		expectedErr  error
	}

	tests := []testCase{
		{
			name: "OK",
			handler: func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(mockResponse{Data: "success"})
			},
			expectedData: mockResponse{Data: "success"},
		},
		{
			name: "HTTP error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "fail", http.StatusInternalServerError)
			},
			expectedErr: ErrBadStatusCode,
		},
		{
			name: "Bad request (400)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "bad request", http.StatusBadRequest)
			},
			expectedErr: ErrNoOffers,
		},
		{
			name: "Bad JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `{"invalid":`)
			},
			expectedErr: ErrDecodeJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			req, _ := http.NewRequest("GET", server.URL, nil)
			logger := zaptest.NewLogger(t)

			res, err := DoJSONRequest[mockResponse](context.Background(), server.Client(), req, logger)

			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected err %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil && res != tt.expectedData {
				t.Fatalf("expected result %+v, got %+v", tt.expectedData, res)
			}
		})
	}
}
