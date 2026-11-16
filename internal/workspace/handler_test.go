package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thakur-vijay/adorbit-go/apperrors"
)

type mockService struct {
	mu sync.Mutex

	createWorkspace *Workspace
	createErr       error
	fetchWorkspaces []Workspace
	fetchErr        error
	createHook      func()

	createCalls int
	fetchCalls  int
	createUser  string
	createUsers []string
	fetchUser   string
	createDTO   CreateWorkspaceDTO
}

func (m *mockService) Create(_ context.Context, userID string, dto CreateWorkspaceDTO) (*Workspace, error) {
	m.mu.Lock()
	m.createCalls++
	m.createUser = userID
	m.createUsers = append(m.createUsers, userID)
	m.createDTO = dto
	workspace, err, hook := m.createWorkspace, m.createErr, m.createHook
	m.mu.Unlock()

	if hook != nil {
		hook()
	}
	return workspace, err
}

func (m *mockService) Fetch(_ context.Context, userID string) ([]Workspace, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fetchCalls++
	m.fetchUser = userID
	return m.fetchWorkspaces, m.fetchErr
}

type handlerResponse struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

func decodeHandlerResponse(t *testing.T, response *httptest.ResponseRecorder) handlerResponse {
	t.Helper()
	var body handlerResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response JSON %q: %v", response.Body.String(), err)
	}
	return body
}

func createRequestBody(t *testing.T, dto CreateWorkspaceDTO) string {
	t.Helper()
	body, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("marshal request DTO: %v", err)
	}
	return string(body)
}

func TestCreateHandler_HTTPBehavior(t *testing.T) {
	userID := uuid.NewString()
	fixture := newJWTTestFixture(t)
	auth := fixture.middleware(t)

	validDTO := testWorkspaceDTO()
	validWorkspace := testWorkspace(userID, validDTO.Name)
	nullFieldsWorkspace := testWorkspace(userID, validDTO.Name)
	nullFieldsWorkspace.Industry = nil
	nullFieldsWorkspace.Website = nil
	nullFieldsWorkspace.Description = nil
	populatedDTO := testWorkspaceDTO()
	populatedWorkspace := testWorkspace(userID, populatedDTO.Name)

	tests := []struct {
		name           string
		body           string
		authenticated  bool
		serviceErr     error
		serviceResult  *Workspace
		expectedStatus int
		expectedMsg    string
		expectedCalls  int
		expectedDTO    *CreateWorkspaceDTO
		expectedResult *Workspace
	}{
		{
			name:           "valid request forwards principal and every DTO field",
			body:           createRequestBody(t, validDTO),
			authenticated:  true,
			serviceResult:  validWorkspace,
			expectedStatus: http.StatusCreated,
			expectedMsg:    "Workspace created successfully",
			expectedCalls:  1,
			expectedDTO:    &validDTO,
			expectedResult: validWorkspace,
		},
		{
			name:           "missing principal is unauthorized",
			body:           createRequestBody(t, validDTO),
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "authentication failed in handler",
		},
		{
			name:           "invalid JSON is rejected",
			body:           `{"name":`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request payload",
		},
		{
			name:           "empty JSON object is rejected",
			body:           `{}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request payload",
		},
		{
			name:           "empty body is rejected",
			body:           "",
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request payload",
		},
		{
			name:           "missing workspace name is rejected",
			body:           `{"type":"publisher","country":"US","timezone":"UTC","currency":"USD"}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "name of type string is required",
		},
		{
			name:           "invalid workspace type is rejected",
			body:           `{"name":"test","type":"agency","country":"US","timezone":"UTC","currency":"USD"}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "type must be one of: advertiser, publisher",
		},
		{
			name:           "invalid country length is rejected",
			body:           `{"name":"test","type":"publisher","country":"USA","timezone":"UTC","currency":"USD"}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "country is invalid",
		},
		{
			name:           "invalid currency length is rejected",
			body:           `{"name":"test","type":"publisher","country":"US","timezone":"UTC","currency":"US"}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "currency is invalid",
		},
		{
			name:           "invalid website URL is rejected",
			body:           `{"name":"test","type":"publisher","country":"US","timezone":"UTC","currency":"USD","website":"example.test"}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "website is invalid",
		},
		{
			name:           "conflict from service is preserved",
			body:           createRequestBody(t, validDTO),
			authenticated:  true,
			serviceErr:     apperrors.NewConflict("workspace already exists", nil),
			expectedStatus: http.StatusConflict,
			expectedMsg:    "workspace already exists",
			expectedCalls:  1,
		},
		{
			name:           "internal service error is classified",
			body:           createRequestBody(t, validDTO),
			authenticated:  true,
			serviceErr:     errors.New("unexpected failure"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "internal server error",
			expectedCalls:  1,
		},
		{
			name:           "unavailable service error is preserved",
			body:           createRequestBody(t, validDTO),
			authenticated:  true,
			serviceErr:     apperrors.NewServiceUnavailable("database unavailable", nil),
			expectedStatus: http.StatusServiceUnavailable,
			expectedMsg:    "database unavailable",
			expectedCalls:  1,
		},
		{
			name:           "nil workspace follows item response contract",
			body:           createRequestBody(t, validDTO),
			authenticated:  true,
			expectedStatus: http.StatusCreated,
			expectedMsg:    "Workspace created successfully",
			expectedCalls:  1,
			expectedDTO:    &validDTO,
		},
		{
			name:           "nullable workspace fields serialize as null",
			body:           createRequestBody(t, validDTO),
			authenticated:  true,
			serviceResult:  nullFieldsWorkspace,
			expectedStatus: http.StatusCreated,
			expectedMsg:    "Workspace created successfully",
			expectedCalls:  1,
			expectedDTO:    &validDTO,
			expectedResult: nullFieldsWorkspace,
		},
		{
			name:           "unknown JSON fields are ignored",
			body:           strings.TrimSuffix(createRequestBody(t, populatedDTO), "}") + `,"extra":"ignored"}`,
			authenticated:  true,
			serviceResult:  populatedWorkspace,
			expectedStatus: http.StatusCreated,
			expectedMsg:    "Workspace created successfully",
			expectedCalls:  1,
			expectedDTO:    &populatedDTO,
			expectedResult: populatedWorkspace,
		},
		{
			name:           "malformed JSON is rejected",
			body:           `{"name":"test",}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request payload",
		},
		{
			name:           "incorrect JSON field type is rejected",
			body:           `{"name":12,"type":"publisher","country":"US","timezone":"UTC","currency":"USD"}`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request payload",
		},
		{
			name:           "trailing JSON data is rejected",
			body:           createRequestBody(t, validDTO) + ` trailing`,
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid request payload",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &mockService{
				createWorkspace: test.serviceResult,
				createErr:       test.serviceErr,
			}
			handler := NewHandler(service)
			request := httptest.NewRequest(http.MethodPost, "/workspaces", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			if test.authenticated {
				request.Header.Set("Authorization", "Bearer "+fixture.token(t, userID, time.Now().Add(time.Hour)))
				auth.RequireAuth(http.HandlerFunc(handler.Create)).ServeHTTP(recorder, request)
			} else {
				handler.Create(recorder, request)
			}

			if recorder.Code != test.expectedStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.expectedStatus, recorder.Body.String())
			}
			body := decodeHandlerResponse(t, recorder)
			if body.StatusCode != test.expectedStatus {
				t.Errorf("response status_code = %d, want %d", body.StatusCode, test.expectedStatus)
			}
			if test.expectedMsg != "" && body.Message != test.expectedMsg {
				t.Errorf("response message = %q, want %q", body.Message, test.expectedMsg)
			}

			service.mu.Lock()
			defer service.mu.Unlock()
			if service.createCalls != test.expectedCalls {
				t.Errorf("Create calls = %d, want %d", service.createCalls, test.expectedCalls)
			}
			if test.expectedCalls == 1 {
				if service.createUser != userID {
					t.Errorf("service UserID = %q, want %q", service.createUser, userID)
				}
				if test.expectedDTO != nil && !reflect.DeepEqual(service.createDTO, *test.expectedDTO) {
					t.Errorf("service DTO = %#v, want %#v", service.createDTO, *test.expectedDTO)
				}
			}
		})
	}
}

func TestFetchHandler_HTTPBehavior(t *testing.T) {
	userID := uuid.NewString()
	fixture := newJWTTestFixture(t)
	auth := fixture.middleware(t)

	first := *testWorkspace(userID, "First workspace")
	second := *testWorkspace(userID, "Second workspace")
	second.Industry = nil
	second.Website = nil
	second.Description = nil
	ordered := []Workspace{first, second}

	tests := []struct {
		name           string
		authenticated  bool
		result         []Workspace
		serviceErr     error
		expectedStatus int
		expectedMsg    string
		expectedJSON   string
		expectedResult []Workspace
	}{
		{
			name:           "multiple workspaces preserve order and nullable fields",
			authenticated:  true,
			result:         ordered,
			expectedStatus: http.StatusOK,
			expectedMsg:    "Workspaces fetched successfully",
			expectedResult: ordered,
		},
		{
			name:           "empty non-nil slice serializes as array",
			authenticated:  true,
			result:         []Workspace{},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Workspaces fetched successfully",
			expectedJSON:   "[]",
			expectedResult: []Workspace{},
		},
		{
			name:           "nil slice serializes as null",
			authenticated:  true,
			result:         nil,
			expectedStatus: http.StatusOK,
			expectedMsg:    "Workspaces fetched successfully",
			expectedJSON:   "null",
		},
		{
			name:           "single workspace is returned",
			authenticated:  true,
			result:         []Workspace{second},
			expectedStatus: http.StatusOK,
			expectedMsg:    "Workspaces fetched successfully",
			expectedResult: []Workspace{second},
		},
		{
			name:           "internal service error is classified",
			authenticated:  true,
			serviceErr:     errors.New("unexpected failure"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "internal server error",
		},
		{
			name:           "unavailable service error is preserved",
			authenticated:  true,
			serviceErr:     apperrors.NewServiceUnavailable("database unavailable", nil),
			expectedStatus: http.StatusServiceUnavailable,
			expectedMsg:    "database unavailable",
		},
		{
			name:           "missing principal is unauthorized",
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "authentication failed in handler",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &mockService{
				fetchWorkspaces: test.result,
				fetchErr:        test.serviceErr,
			}
			handler := NewHandler(service)
			request := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
			recorder := httptest.NewRecorder()

			if test.authenticated {
				request.Header.Set("Authorization", "Bearer "+fixture.token(t, userID, time.Now().Add(time.Hour)))
				auth.RequireAuth(http.HandlerFunc(handler.Fetch)).ServeHTTP(recorder, request)
			} else {
				handler.Fetch(recorder, request)
			}

			if recorder.Code != test.expectedStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.expectedStatus, recorder.Body.String())
			}
			body := decodeHandlerResponse(t, recorder)
			if body.StatusCode != test.expectedStatus {
				t.Errorf("response status_code = %d, want %d", body.StatusCode, test.expectedStatus)
			}
			if test.expectedMsg != "" && body.Message != test.expectedMsg {
				t.Errorf("response message = %q, want %q", body.Message, test.expectedMsg)
			}
			if test.expectedJSON != "" && string(body.Data) != test.expectedJSON {
				t.Errorf("response data = %s, want %s", body.Data, test.expectedJSON)
			}
			if test.expectedResult != nil {
				var got []Workspace
				if err := json.Unmarshal(body.Data, &got); err != nil {
					t.Fatalf("decode workspace array: %v", err)
				}
				if !reflect.DeepEqual(got, test.expectedResult) {
					t.Errorf("response workspaces = %#v, want %#v", got, test.expectedResult)
				}
			}

			service.mu.Lock()
			defer service.mu.Unlock()
			wantCalls := 0
			if test.authenticated {
				wantCalls = 1
			}
			if service.fetchCalls != wantCalls {
				t.Errorf("Fetch calls = %d, want %d", service.fetchCalls, wantCalls)
			}
			if wantCalls == 1 && service.fetchUser != userID {
				t.Errorf("service UserID = %q, want %q", service.fetchUser, userID)
			}
		})
	}
}
