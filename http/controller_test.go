package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStateRoundTripPreservesPayload(t *testing.T) {
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	payload := `{"version":4,"serial":1,"outputs":{"quote":"it's ok"}}`

	update := requestContext(http.MethodPost, "/group/key", "group", "key", payload)
	Update(update.response, update.request)
	if update.response.Code != 200 {
		t.Fatalf("expected status 200, got %d", update.response.Code)
	}
	if strings.TrimSpace(update.response.Body.String()) != payload {
		t.Fatalf("expected update body %q, got %q", payload, update.response.Body.String())
	}

	show := requestContext(http.MethodGet, "/group/key", "group", "key", "")
	Show(show.response, show.request)
	if show.response.Code != 200 {
		t.Fatalf("expected status 200, got %d", show.response.Code)
	}
	if strings.TrimSpace(show.response.Body.String()) != payload {
		t.Fatalf("expected show body %q, got %q", payload, show.response.Body.String())
	}
}

func TestShowMissingStateReturnsNotFound(t *testing.T) {
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	show := requestContext(http.MethodGet, "/group/missing", "group", "missing", "")
	Show(show.response, show.request)
	if show.response.Code != 404 {
		t.Fatalf("expected status 404, got %d", show.response.Code)
	}
}

func TestLockConflictAndUnlock(t *testing.T) {
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	firstLock := `{"ID":"first"}`
	secondLock := `{"ID":"second"}`

	lock := requestContext(http.MethodPut, "/group/key", "group", "key", firstLock)
	Lock(lock.response, lock.request)
	if lock.response.Code != 200 {
		t.Fatalf("expected first lock status 200, got %d", lock.response.Code)
	}

	conflict := requestContext(http.MethodPut, "/group/key", "group", "key", secondLock)
	Lock(conflict.response, conflict.request)
	if conflict.response.Code != 423 {
		t.Fatalf("expected lock conflict status 423, got %d", conflict.response.Code)
	}
	if strings.TrimSpace(conflict.response.Body.String()) != firstLock {
		t.Fatalf("expected current lock body %q, got %q", firstLock, conflict.response.Body.String())
	}

	wrongUnlock := requestContext(http.MethodDelete, "/group/key", "group", "key", secondLock)
	Unlock(wrongUnlock.response, wrongUnlock.request)
	if wrongUnlock.response.Code != 409 {
		t.Fatalf("expected wrong unlock status 409, got %d", wrongUnlock.response.Code)
	}

	unlock := requestContext(http.MethodDelete, "/group/key", "group", "key", firstLock)
	Unlock(unlock.response, unlock.request)
	if unlock.response.Code != 200 {
		t.Fatalf("expected unlock status 200, got %d", unlock.response.Code)
	}

	relock := requestContext(http.MethodPut, "/group/key", "group", "key", secondLock)
	Lock(relock.response, relock.request)
	if relock.response.Code != 200 {
		t.Fatalf("expected relock status 200, got %d", relock.response.Code)
	}
}

type handlerContext struct {
	request  *http.Request
	response *httptest.ResponseRecorder
}

func requestContext(method string, target string, group string, key string, body string) handlerContext {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.SetPathValue("group", group)
	request.SetPathValue("key", key)
	response := httptest.NewRecorder()

	return handlerContext{
		request:  request,
		response: response,
	}
}
