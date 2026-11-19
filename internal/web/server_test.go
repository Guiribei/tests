package web

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/store"
)

func TestUnknownPageOffersRecovery(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "provas.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s, err := New(Config{Catalog: &exam.Catalog{}, Store: db})
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/pagina-inexistente", nil))
	if w.Code != http.StatusNotFound || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("expected an HTML 404 page, got %d %q", w.Code, w.Header().Get("Content-Type"))
	}
	for _, content := range []string{"Página não encontrada", "Ver provas disponíveis", `href="/"`} {
		if !strings.Contains(w.Body.String(), content) {
			t.Errorf("error page is missing %q", content)
		}
	}
}
