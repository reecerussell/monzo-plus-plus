package routing

import (
	"fmt"
	"net/http"
	"testing"
)

var testRoutes = map[string]string{
	"/api/plugin/{name}/": "/api/plugin/budget/he",
}

var testHandler = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func TestMatch(t *testing.T) {
	r := NewRouter()
	var rt *Route
	for t, _ := range testRoutes {
		rt = r.GetFunc(t, testHandler)
	}

	for _, p := range testRoutes {
		if fr := r.match(p, http.MethodGet); fr != rt {
			t.Errorf("expected a handler")
			return
		}
	}
}

func TestVars(t *testing.T) {
	r := &Route{Template: "/api/user"}
	id := "hello"
	vars := r.getVars(fmt.Sprintf("/api/user/%s", id))

	if v, ok := vars["id"]; !ok || v != id {
		t.Errorf("expected '%s' but got '%s'", id, v)
		return
	}
}
