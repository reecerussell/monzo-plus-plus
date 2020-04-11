package routing

import (
	"fmt"
	"net/http"
	"testing"
)

var testRoutes = map[string]string{
	"/api/plugin/{name}/": "/api/plugin/budget/preferences",
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
	r := &Route{Template: "/api/plugin/{name}/"}
	name := "budget"
	vars := r.getVars(fmt.Sprintf("/api/plugin/%s/preferences", name))

	if v, ok := vars["name"]; !ok || v != name {
		t.Errorf("expected '%s' but got '%s'", name, v)
		return
	}
}
