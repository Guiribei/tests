package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/store"
)

func TestOfflinePacketUsesAttemptSnapshot(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "provas.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	e, err := exam.LoadFile("../../option-one.json")
	if err != nil {
		t.Fatal(err)
	}
	session := randomID()
	a, err := db.StartOrResume(t.Context(), randomID(), session, e)
	if err != nil {
		t.Fatal(err)
	}
	// The current catalog is empty. This download must still serve the snapshot.
	s, err := New(Config{Catalog: &exam.Catalog{}, Store: db})
	if err != nil {
		t.Fatal(err)
	}
	get := func(cookie string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, "/tentativas/"+a.ID+"/pacote", nil)
		r.AddCookie(&http.Cookie{Name: "prova_session", Value: cookie})
		w := httptest.NewRecorder()
		s.ServeHTTP(w, r)
		return w
	}
	if w := get(randomID()); w.Code != http.StatusNotFound {
		t.Fatalf("foreign packet: %d", w.Code)
	}
	w := get(session)
	if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("packet: %d %v", w.Code, w.Header())
	}
	var packet struct {
		store.Sheet
		Title     string
		Role      string
		Questions []offlineQuestion
	}
	decoder := json.NewDecoder(w.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&packet); err != nil {
		t.Fatal(err)
	}
	if len(packet.Questions) != 180 || packet.Title != a.ExamTitle || len(packet.Answers) != 0 {
		t.Fatalf("unexpected packet for %s", a.ID)
	}
	q := packet.Questions[0]
	if q.Statement != e.Questions[0].Statement || len(q.Alternatives) != 5 || q.BaseText != e.Questions[0].BaseText {
		t.Fatal("incomplete first question")
	}
}
