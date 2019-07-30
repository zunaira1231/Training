package functional

import (
	"net/http"
	"net/http/httptest"
	"restfulAPI2/handler"
	"strings"
	"testing"
)

func TestGetSamples(t *testing.T) {
	testServer := setupServer()
	//returns a new incoming server Request, suitable for passing to an http.Handler for testing.
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/samples", nil)
	if err != nil {
		t.Fatal(err)
	}

	//nitialized ResponseRecorder
	rec := httptest.NewRecorder()

	testServer.ServeHTTP(rec, req)

got := strings.TrimSpace(rec.Body.String())

want := `[{"id":1,"title":"Do dishes","note":"","due_date":"2000-01-01T00:00:00Z"},{"id":2,"title":"Do homework","note":"","due_date":"2000-01-01T00:00:00Z"},{"id":2,"title":"Twitter","note":"","due_date":"2000-01-01T00:00:00Z"}]`

if got != want {
t.Fatalf("Want: %v, Got: %v", want, got)
}
}

//set uprouting
func setupServer() *http.ServeMux {
	return handler.SetUpRouting()
}

/*
// ...
func TestGetSamples(t *testing.T) {
    testServer := setupServer(nil)
//...
func setupServer(postgres *db.Postgres) *http.ServeMux {
    return handler.SetUpRouting(postgres)
}
// ...
func TestGetAllTodo(t *testing.T) {
    postgres := &db.Postgres{testdb.Setup()}
    testServer := setupServer(postgres)

    todo := &schema.Todo{
        Title:   "My Task1",
        DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
    }

    _, err := postgres.Insert(todo)
    if err != nil {
        t.Fatal(err)
    }

    req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/todo", nil)
    if err != nil {
        t.Fatal(err)
    }

    rec := httptest.NewRecorder()
    testServer.ServeHTTP(rec, req)

    got := strings.TrimSpace(rec.Body.String())

    want := `[{"id":1,"title":"My Task1","note":"","due_date":"2000-01-01T00:00:00+09:00"}]`

    if got != want {
        t.Fatalf("Want: %v, Got: %v", want, got)
    }
}

func TestSaveTodo(t *testing.T) {
    postgres := &db.Postgres{testdb.Setup()}
    testServer := setupServer(postgres)

    body := []byte(`{"id":1,"title":"My Task1","note":"","due_date":"2000-01-01T00:00:00+09:00"}`)

    req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/todo", bytes.NewReader(body))
    if err != nil {
        t.Fatal(err)
    }

    rec := httptest.NewRecorder()
    testServer.ServeHTTP(rec, req)

    got := strings.TrimSpace(rec.Body.String())
    want := "1"

    if got != want {
        t.Fatalf("Want: %v, Got: %v", want, got)
    }

    gotTodo, err := postgres.GetAll()
    if err != nil {
        t.Fatal(err)
    }

    wantTodo := []schema.Todo{
        {
            Title:   "My Task1",
            DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
        },
    }

    if !reflect.DeepEqual(got, want) {
        t.Fatalf("Want: %v, Got: %v\n", wantTodo, gotTodo)
    }
}

func TestDeleteTodo(t *testing.T) {
    postgres := &db.Postgres{testdb.Setup()}
    testServer := setupServer(postgres)

    todo := &schema.Todo{
        Title:   "My Task1",
        DueDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
    }

    id, err := postgres.Insert(todo)
    if err != nil {
        t.Fatal(err)
    }

    body := []byte(fmt.Sprintf(`{"id":%d}`, id))

    req, err := http.NewRequest(http.MethodDelete, "http://localhost:9999/todo", bytes.NewReader(body))
    if err != nil {
        t.Fatal(err)
    }

    rec := httptest.NewRecorder()
    testServer.ServeHTTP(rec, req)

    got := rec.Body.String()

    want := ""

    if got != want {
        t.Fatalf("Want: %v, Got: %v", want, got)
    }

    gotTodo, err := postgres.GetAll()
    if err != nil {
        t.Fatal(err)
    }

    if len(gotTodo) > 0 {
        t.Fatalf("Should return the empty slice, Got: %v\n", gotTodo)
    }
}
 */
