package workspace

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func workspaceRouteHandler(t *testing.T, service *mockService) (http.Handler, *jwtTestFixture) {
	t.Helper()
	db := workspaceTestDB(t, true)
	fixture := newJWTTestFixture(t)
	mux := http.NewServeMux()
	RegisterRoutes(mux, NewHandler(service), fixture.middleware(t), newIdempotencyMiddleware(t, db))
	return mux, fixture
}

func callWorkspaceRoute(t *testing.T, handler http.Handler, method, token, key, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, "/workspaces", strings.NewReader(body))
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	if key != "" {
		request.Header.Set("Idempotency-Key", key)
	}
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func routeToken(t *testing.T, fixture *jwtTestFixture, userID string) string {
	t.Helper()
	return fixture.token(t, userID, time.Now().Add(time.Hour))
}

func TestWorkspaceRoutes_AuthenticationAndIdempotencyKeyValidation(t *testing.T) {
	service := &mockService{}
	handler, fixture := workspaceRouteHandler(t, service)
	userID := uuid.NewString()
	expiredToken := fixture.token(t, userID, time.Now().Add(-time.Hour))
	validToken := routeToken(t, fixture, userID)

	tests := []struct {
		name           string
		method         string
		token          string
		key            string
		setKey         bool
		expectedStatus int
	}{
		{name: "POST without JWT", method: http.MethodPost, expectedStatus: http.StatusUnauthorized},
		{name: "GET without JWT", method: http.MethodGet, expectedStatus: http.StatusUnauthorized},
		{name: "POST with expired JWT", method: http.MethodPost, token: expiredToken, key: uuid.NewString(), expectedStatus: http.StatusUnauthorized},
		{name: "GET with expired JWT", method: http.MethodGet, token: expiredToken, expectedStatus: http.StatusUnauthorized},
		{name: "POST with malformed JWT", method: http.MethodPost, token: "not.a.jwt", key: uuid.NewString(), expectedStatus: http.StatusUnauthorized},
		{name: "GET with malformed JWT", method: http.MethodGet, token: "not.a.jwt", expectedStatus: http.StatusUnauthorized},
		{name: "POST without idempotency key", method: http.MethodPost, token: validToken, expectedStatus: http.StatusBadRequest},
		{name: "POST with empty idempotency key", method: http.MethodPost, token: validToken, setKey: true, expectedStatus: http.StatusBadRequest},
		{name: "POST with malformed idempotency key", method: http.MethodPost, token: validToken, key: "not-a-uuid", expectedStatus: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/workspaces", strings.NewReader(createRequestBody(t, testWorkspaceDTO())))
			if test.token != "" {
				request.Header.Set("Authorization", "Bearer "+test.token)
			}
			if test.key != "" || test.setKey {
				request.Header.Set("Idempotency-Key", test.key)
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			if recorder.Code != test.expectedStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.expectedStatus, recorder.Body.String())
			}
		})
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if service.createCalls != 0 || service.fetchCalls != 0 {
		t.Errorf("rejected requests reached service: Create calls=%d Fetch calls=%d", service.createCalls, service.fetchCalls)
	}
}

func TestWorkspaceRoutes_AuthenticatedRequestsAndMethods(t *testing.T) {
	userID := uuid.NewString()
	service := &mockService{
		createWorkspace: testWorkspace(userID, "Route workspace"),
		fetchWorkspaces: []Workspace{*testWorkspace(userID, "Route workspace")},
	}
	handler, fixture := workspaceRouteHandler(t, service)
	token := routeToken(t, fixture, userID)
	dto := testWorkspaceDTO()
	key := uuid.NewString()

	created := callWorkspaceRoute(t, handler, http.MethodPost, token, key, createRequestBody(t, dto))
	if created.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d; body = %s", created.Code, http.StatusCreated, created.Body.String())
	}

	fetched := callWorkspaceRoute(t, handler, http.MethodGet, token, "", "")
	if fetched.Code != http.StatusOK {
		t.Fatalf("GET without Idempotency-Key status = %d, want %d; body = %s", fetched.Code, http.StatusOK, fetched.Body.String())
	}

	unsupported := callWorkspaceRoute(t, handler, http.MethodPatch, token, "", "")
	if unsupported.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PATCH status = %d, want %d", unsupported.Code, http.StatusMethodNotAllowed)
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if service.createCalls != 1 || service.fetchCalls != 1 {
		t.Errorf("service calls = (Create %d, Fetch %d), want (1, 1)", service.createCalls, service.fetchCalls)
	}
	if service.createUser != userID || service.fetchUser != userID {
		t.Errorf("service users = (Create %q, Fetch %q), want %q", service.createUser, service.fetchUser, userID)
	}
}

func TestWorkspaceRoutes_IdempotencyReplayConflictAndUserScope(t *testing.T) {
	userA, userB := uuid.NewString(), uuid.NewString()
	service := &mockService{createWorkspace: testWorkspace(userA, "Idempotent workspace")}
	handler, fixture := workspaceRouteHandler(t, service)
	tokenA := routeToken(t, fixture, userA)
	tokenB := routeToken(t, fixture, userB)
	key := uuid.NewString()
	body := createRequestBody(t, testWorkspaceDTO())

	first := callWorkspaceRoute(t, handler, http.MethodPost, tokenA, key, body)
	if first.Code != http.StatusCreated {
		t.Fatalf("first POST status = %d, want %d; body = %s", first.Code, http.StatusCreated, first.Body.String())
	}
	firstBody := append([]byte(nil), first.Body.Bytes()...)

	replay := callWorkspaceRoute(t, handler, http.MethodPost, tokenA, key, body)
	if replay.Code != first.Code {
		t.Errorf("same-key replay status = %d, want %d", replay.Code, first.Code)
	}
	if !bytes.Equal(replay.Body.Bytes(), firstBody) {
		t.Errorf("same-key replay body = %s, want exact stored response %s", replay.Body.String(), firstBody)
	}

	differentPayload := createRequestBody(t, func() CreateWorkspaceDTO {
		dto := testWorkspaceDTO()
		dto.Name = "Different workspace"
		return dto
	}())
	conflict := callWorkspaceRoute(t, handler, http.MethodPost, tokenA, key, differentPayload)
	if conflict.Code != http.StatusConflict {
		t.Errorf("same key with different payload status = %d, want %d; body = %s", conflict.Code, http.StatusConflict, conflict.Body.String())
	}

	separateScope := callWorkspaceRoute(t, handler, http.MethodPost, tokenB, key, body)
	if separateScope.Code != http.StatusCreated {
		t.Errorf("same key in a different user scope status = %d, want %d; body = %s", separateScope.Code, http.StatusCreated, separateScope.Body.String())
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if service.createCalls != 2 {
		t.Errorf("Create calls = %d, want 2 (one per user scope)", service.createCalls)
	}
	if !reflect.DeepEqual(service.createUsers, []string{userA, userB}) {
		t.Errorf("Create users = %v, want [%s %s]", service.createUsers, userA, userB)
	}
}

func TestWorkspaceRoutes_ConcurrentSameKeyRunsCreateOnce(t *testing.T) {
	userID := uuid.NewString()
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	service := &mockService{
		createWorkspace: testWorkspace(userID, "Concurrent workspace"),
		createHook: func() {
			once.Do(func() {
				close(entered)
				<-release
			})
		},
	}
	handler, fixture := workspaceRouteHandler(t, service)
	token := routeToken(t, fixture, userID)
	key := uuid.NewString()
	body := createRequestBody(t, testWorkspaceDTO())

	const requestCount = 8
	type routeResult struct {
		status int
		body   string
	}
	results := make(chan routeResult, requestCount)
	start := make(chan struct{})
	for i := 0; i < requestCount; i++ {
		go func() {
			<-start
			request := httptest.NewRequest(http.MethodPost, "/workspaces", strings.NewReader(body))
			request.Header.Set("Authorization", "Bearer "+token)
			request.Header.Set("Idempotency-Key", key)
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			results <- routeResult{status: recorder.Code, body: recorder.Body.String()}
		}()
	}
	close(start)

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("no request reached the Create service")
	}

	for i := 0; i < requestCount-1; i++ {
		select {
		case result := <-results:
			if result.status != http.StatusConflict {
				t.Errorf("concurrent duplicate status = %d, want %d; body = %s", result.status, http.StatusConflict, result.body)
			}
		case <-time.After(5 * time.Second):
			close(release)
			t.Fatal("concurrent duplicate requests did not complete")
		}
	}
	close(release)

	select {
	case result := <-results:
		if result.status != http.StatusCreated {
			t.Errorf("owner request status = %d, want %d; body = %s", result.status, http.StatusCreated, result.body)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("owner request did not complete")
	}

	service.mu.Lock()
	defer service.mu.Unlock()
	if service.createCalls != 1 {
		t.Errorf("concurrent Create calls = %d, want exactly 1", service.createCalls)
	}
}
