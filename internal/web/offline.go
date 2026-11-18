package web

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/store"
)

func (s *Server) serviceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFileFS(w, r, assets, "static/sw.js")
}

// The cached navigation shell contains no session or attempt data.
func (s *Server) offlinePage(w http.ResponseWriter, r *http.Request) {
	data := pageData{Title: "Provas no aparelho", Exam: &exam.Exam{}}
	for _, letter := range []string{"A", "B", "C", "D", "E"} {
		data.Question.Alternatives = append(data.Question.Alternatives, alternativeView{Letter: letter})
	}
	s.render(w, r, "question", http.StatusOK, data)
}

type offlineQuestion struct {
	Number       int      `json:"number"`
	Statement    string   `json:"statement"`
	Alternatives []string `json:"alternatives"`
	Discipline   string   `json:"discipline"`
	Paper        string   `json:"paper"`
	BaseID       string   `json:"baseID,omitempty"`
	BaseText     string   `json:"baseText,omitempty"`
}

func (s *Server) packet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	e, sheet, err := s.store.Packet(r.Context(), s.sessionID(w, r), r.PathValue("id"))
	if err != nil {
		s.answerError(w, r, err)
		return
	}
	questions := make([]offlineQuestion, 0, len(e.Questions))
	for _, q := range e.Questions {
		questions = append(questions, offlineQuestion{q.Number, q.Statement, q.Alternatives, q.Discipline, q.Paper, q.BaseTextID, q.BaseText})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		store.Sheet
		Title     string            `json:"title"`
		Role      string            `json:"role"`
		Questions []offlineQuestion `json:"questions"`
	}{sheet, e.Document.Metadata.Title, e.Document.Metadata.RoleCode, questions})
}

func (s *Server) syncAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaType != "application/json" {
		http.Error(w, "Esperado JSON", http.StatusUnsupportedMediaType)
		return
	}
	var batch struct {
		Revision *int64         `json:"revision"`
		Answers  map[int]string `json:"answers"`
		Finish   bool           `json:"finish"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&batch); err != nil || batch.Revision == nil || *batch.Revision < 0 || batch.Answers == nil || dec.Decode(&struct{}{}) != io.EOF {
		http.Error(w, "Lote de respostas inválido", http.StatusBadRequest)
		return
	}
	sheet, err := s.store.ApplyAnswers(r.Context(), s.sessionID(w, r), r.PathValue("id"), batch.Revision, batch.Answers, batch.Finish)
	status := http.StatusOK
	if errors.Is(err, store.ErrRevision) {
		status = http.StatusConflict
	} else if err != nil {
		s.answerError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(sheet)
}

func (s *Server) answerError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		s.notFound(w, r)
	case errors.Is(err, store.ErrFinished):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, store.ErrInvalidAnswer):
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	default:
		s.internalError(w, r, err)
	}
}
