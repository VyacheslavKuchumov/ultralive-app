package crmhttp

import (
	"VyacheslavKuchumov/test-backend/service/tracker"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteStoreErrorMapsStatuses(t *testing.T) {
	testCases := []struct {
		name       string
		err        error
		statusCode int
	}{
		{
			name:       "not found",
			err:        tracker.ErrNotFound,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "invalid reference",
			err:        tracker.ErrInvalidReference,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "unexpected error",
			err:        errors.New("unexpected failure"),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			WriteStoreError(rr, tc.err)

			if rr.Code != tc.statusCode {
				t.Fatalf("expected status %d, got %d", tc.statusCode, rr.Code)
			}
		})
	}
}
