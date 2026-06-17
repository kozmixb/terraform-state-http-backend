package http

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestStateRoundTripPreservesPayload(t *testing.T) {
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	payload := `{"version":4,"serial":1,"outputs":{"quote":"it's ok"}}`

	update := requestContext(echo.POST, "/group/key", "group", "key", payload)
	if err := Update(update.context); err != nil {
		t.Fatal(err)
	}
	if update.response.Code != 200 {
		t.Fatalf("expected status 200, got %d", update.response.Code)
	}
	if strings.TrimSpace(update.response.Body.String()) != payload {
		t.Fatalf("expected update body %q, got %q", payload, update.response.Body.String())
	}

	show := requestContext(echo.GET, "/group/key", "group", "key", "")
	if err := Show(show.context); err != nil {
		t.Fatal(err)
	}
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

	show := requestContext(echo.GET, "/group/missing", "group", "missing", "")
	if err := Show(show.context); err != nil {
		t.Fatal(err)
	}
	if show.response.Code != 404 {
		t.Fatalf("expected status 404, got %d", show.response.Code)
	}
}

func TestLockConflictAndUnlock(t *testing.T) {
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	firstLock := `{"ID":"first"}`
	secondLock := `{"ID":"second"}`

	lock := requestContext(echo.PUT, "/group/key", "group", "key", firstLock)
	if err := Lock(lock.context); err != nil {
		t.Fatal(err)
	}
	if lock.response.Code != 200 {
		t.Fatalf("expected first lock status 200, got %d", lock.response.Code)
	}

	conflict := requestContext(echo.PUT, "/group/key", "group", "key", secondLock)
	if err := Lock(conflict.context); err != nil {
		t.Fatal(err)
	}
	if conflict.response.Code != 423 {
		t.Fatalf("expected lock conflict status 423, got %d", conflict.response.Code)
	}
	if strings.TrimSpace(conflict.response.Body.String()) != firstLock {
		t.Fatalf("expected current lock body %q, got %q", firstLock, conflict.response.Body.String())
	}

	wrongUnlock := requestContext(echo.DELETE, "/group/key", "group", "key", secondLock)
	if err := Unlock(wrongUnlock.context); err != nil {
		t.Fatal(err)
	}
	if wrongUnlock.response.Code != 409 {
		t.Fatalf("expected wrong unlock status 409, got %d", wrongUnlock.response.Code)
	}

	unlock := requestContext(echo.DELETE, "/group/key", "group", "key", firstLock)
	if err := Unlock(unlock.context); err != nil {
		t.Fatal(err)
	}
	if unlock.response.Code != 200 {
		t.Fatalf("expected unlock status 200, got %d", unlock.response.Code)
	}

	relock := requestContext(echo.PUT, "/group/key", "group", "key", secondLock)
	if err := Lock(relock.context); err != nil {
		t.Fatal(err)
	}
	if relock.response.Code != 200 {
		t.Fatalf("expected relock status 200, got %d", relock.response.Code)
	}
}

type handlerContext struct {
	context  echo.Context
	response *httptest.ResponseRecorder
}

func requestContext(method string, target string, group string, key string, body string) handlerContext {
	e := echo.New()
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := e.NewContext(request, response)
	context.SetParamNames("group", "key")
	context.SetParamValues(group, key)

	return handlerContext{
		context:  context,
		response: response,
	}
}
