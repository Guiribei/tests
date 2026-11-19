package web

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Guiribei/tests/internal/exam"
	"github.com/Guiribei/tests/internal/scoring"
	"github.com/Guiribei/tests/internal/store"
)

//go:embed templates/*.html static/*
var assets embed.FS

type Config struct {
	Catalog       *exam.Catalog
	Store         *store.Store
	Logger        *slog.Logger
	RepositoryURL string
	SecureCookies bool
}

type Server struct {
	catalog       *exam.Catalog
	store         *store.Store
	logger        *slog.Logger
	repositoryURL string
	secureCookies bool
	templates     map[string]*template.Template
	handler       http.Handler
}

type pageData struct {
	Title         string
	RepositoryURL string
	Exams         []examCard
	Attempt       store.Attempt
	Exam          *exam.Exam
	Question      questionView
	Answers       map[int]string
	Answered      int
	Total         int
	Progress      float64
	Previous      int
	Next          int
	Score         scoring.Score
	Cohort        scoring.CohortScore
	Results       []questionResult
	Disciplines   []disciplineResult
	ErrorMessage  string
}

type examCard struct {
	Number            int
	DisplayTitle      string
	Subtitle          string
	TotalQuestions    int
	Slug              string
	Title             string
	RoleCode          string
	Type              string
	Notice            string
	SourceFile        string
	Hash              string
	GeneralQuestions  int
	SpecificQuestions int
	TotalMinutes      int
}

type questionView struct {
	exam.Question
	Alternatives []alternativeView
}

type alternativeView struct {
	Letter   string
	Text     string
	Selected bool
}

type questionResult struct {
	Number        int
	Discipline    string
	Paper         string
	Statement     string
	BaseText      string
	Alternatives  []resultAlternative
	UserAnswer    string
	CorrectAnswer string
	Correct       bool
	Blank         bool
}

type resultAlternative struct {
	Letter  string
	Text    string
	Correct bool
	Chosen  bool
}

type disciplineResult struct {
	Name    string
	Correct int
	Total   int
}

func New(config Config) (*Server, error) {
	if config.Catalog == nil || config.Store == nil {
		return nil, errors.New("catálogo e banco são obrigatórios")
	}
	server := &Server{
		catalog: config.Catalog, store: config.Store, logger: config.Logger,
		repositoryURL: config.RepositoryURL, secureCookies: config.SecureCookies,
	}
	if server.logger == nil {
		server.logger = slog.Default()
	}
	if err := server.parseTemplates(); err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	staticFS, err := fs.Sub(assets, "static")
	if err != nil {
		return nil, err
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	mux.HandleFunc("GET /{$}", server.home)
	mux.HandleFunc("POST /provas/{slug}/iniciar", server.startAttempt)
	mux.HandleFunc("GET /tentativas/{id}/questoes/{number}", server.question)
	mux.HandleFunc("POST /tentativas/{id}/respostas/{number}", server.saveAnswer)
	mux.HandleFunc("POST /tentativas/{id}/finalizar", server.finishAttempt)
	mux.HandleFunc("GET /tentativas/{id}/resultado", server.result)
	mux.HandleFunc("GET /sw.js", server.serviceWorker)
	mux.HandleFunc("GET /modo-offline", server.offlinePage)
	mux.HandleFunc("GET /tentativas/{id}/pacote", server.packet)
	mux.HandleFunc("POST /tentativas/{id}/sincronizar", server.syncAnswers)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	server.handler = server.recoverer(server.securityHeaders(mux))
	return server, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.handler.ServeHTTP(w, r) }

func (s *Server) parseTemplates() error {
	functions := template.FuncMap{
		"letters": func() []string { return []string{"A", "B", "C", "D", "E"} },
		"percent": func(value float64) string { return fmt.Sprintf("%.1f", value) },
		"score":   func(value float64) string { return fmt.Sprintf("%.2f", value) },
		"float":   func(value int) float64 { return float64(value) },
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mul": func(a, b float64) float64 { return a * b },
		"duration": func(start time.Time, finished *time.Time) string {
			end := time.Now().UTC()
			if finished != nil {
				end = *finished
			}
			d := end.Sub(start).Round(time.Minute)
			if d < time.Minute {
				return "menos de 1 min"
			}
			return fmt.Sprintf("%dh%02d", int(d.Hours()), int(d.Minutes())%60)
		},
	}
	s.templates = make(map[string]*template.Template)
	for _, page := range []string{"home", "question", "result", "error"} {
		t, err := template.New("base.html").Funcs(functions).ParseFS(assets, "templates/base.html", "templates/"+page+".html")
		if err != nil {
			return fmt.Errorf("template %s: %w", page, err)
		}
		s.templates[page] = t
	}
	return nil
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, name string, status int, data pageData) {
	t, ok := s.templates[name]
	if !ok {
		http.Error(w, "template indisponível", http.StatusInternalServerError)
		return
	}
	data.RepositoryURL = s.repositoryURL
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	templateName := "base"
	if r.Header.Get("HX-Request") == "true" && name != "error" {
		templateName = "content"
	}
	if err := t.ExecuteTemplate(w, templateName, data); err != nil {
		s.logger.Error("falha ao renderizar template", "template", name, "error", err)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.sessionID(w, r)
	data := pageData{Title: "Índice de provas"}
	for index, item := range s.catalog.All() {
		card := examCard{
			Number: index + 1, TotalQuestions: len(item.Questions),
			Slug: item.Slug, Title: item.Document.Metadata.Title, RoleCode: item.Document.Metadata.RoleCode,
			Type: item.Document.Metadata.Type, Notice: item.Document.Metadata.Notice,
			SourceFile: item.SourceFile, Hash: item.Hash,
		}
		// Keep the complete source title in the record and expose its natural
		// title/subtitle split for scanning the catalog.
		title := strings.ReplaceAll(card.Title, " – ", " — ")
		parts := strings.Split(title, " — ")
		if len(parts) > 1 && strings.EqualFold(parts[0], card.Type) {
			parts = parts[1:]
		}
		card.DisplayTitle = strings.TrimPrefix(parts[0], "Simulado ")
		card.Subtitle = strings.Join(parts[1:], " · ")
		card.GeneralQuestions = len(flatten(item.Document.General))
		card.SpecificQuestions = len(flatten(item.Document.Specific))
		for _, paper := range item.Document.Metadata.Papers {
			card.TotalMinutes += paper.DurationMinutes
		}
		data.Exams = append(data.Exams, card)
	}
	s.render(w, r, "home", http.StatusOK, data)
}

func flatten(sections []exam.Section) []exam.QuestionJSON {
	var questions []exam.QuestionJSON
	for _, section := range sections {
		questions = append(questions, section.Questions...)
	}
	return questions
}

func (s *Server) startAttempt(w http.ResponseWriter, r *http.Request) {
	sessionID := s.sessionID(w, r)
	selected, ok := s.catalog.Get(r.PathValue("slug"))
	if !ok {
		s.notFound(w, r)
		return
	}
	attempt, err := s.store.StartOrResume(r.Context(), randomID(), sessionID, selected)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	location := fmt.Sprintf("/tentativas/%s/questoes/1?retomar=1", attempt.ID)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", location)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, location, http.StatusSeeOther)
}

func (s *Server) question(w http.ResponseWriter, r *http.Request) {
	sessionID := s.sessionID(w, r)
	attempt, currentExam, answers, ok := s.loadAttempt(w, r, sessionID)
	if !ok {
		return
	}
	if attempt.Status == "finished" {
		http.Redirect(w, r, fmt.Sprintf("/tentativas/%s/resultado", attempt.ID), http.StatusSeeOther)
		return
	}
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	question, exists := currentExam.Question(number)
	if !exists {
		s.notFound(w, r)
		return
	}
	view := questionView{Question: question}
	for index, text := range question.Alternatives {
		letter := string(rune('A' + index))
		view.Alternatives = append(view.Alternatives, alternativeView{Letter: letter, Text: text, Selected: answers[number] == letter})
	}
	data := pageData{
		Title: fmt.Sprintf("Questão %d", number), Attempt: attempt, Exam: currentExam,
		Question: view, Answers: answers, Answered: len(answers), Total: len(currentExam.Questions),
		Progress: float64(len(answers)) / float64(len(currentExam.Questions)) * 100,
	}
	if number > 1 {
		data.Previous = number - 1
	}
	if number < len(currentExam.Questions) {
		data.Next = number + 1
	}
	s.render(w, r, "question", http.StatusOK, data)
}

func (s *Server) saveAnswer(w http.ResponseWriter, r *http.Request) {
	sessionID := s.sessionID(w, r)
	id := r.PathValue("id")
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil {
		s.notFound(w, r)
		return
	}
	choice := strings.ToUpper(r.FormValue("choice"))
	if _, err := s.store.ApplyAnswers(r.Context(), sessionID, id, nil, map[int]string{number: choice}, false); err != nil {
		s.answerError(w, r, err)
		return
	}
	if r.Header.Get("HX-Request") != "true" {
		http.Redirect(w, r, fmt.Sprintf("/tentativas/%s/questoes/%d", id, number), http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<span class="save-ok" role="status">Salvo agora</span>`))
}

func (s *Server) finishAttempt(w http.ResponseWriter, r *http.Request) {
	sessionID := s.sessionID(w, r)
	id := r.PathValue("id")
	if _, err := s.store.ApplyAnswers(r.Context(), sessionID, id, nil, nil, true); err != nil {
		s.answerError(w, r, err)
		return
	}
	location := fmt.Sprintf("/tentativas/%s/resultado", id)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", location)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, location, http.StatusSeeOther)
}

func (s *Server) result(w http.ResponseWriter, r *http.Request) {
	sessionID := s.sessionID(w, r)
	attempt, currentExam, answers, ok := s.loadAttempt(w, r, sessionID)
	if !ok {
		return
	}
	if attempt.Status != "finished" {
		http.Redirect(w, r, fmt.Sprintf("/tentativas/%s/questoes/1", attempt.ID), http.StatusSeeOther)
		return
	}
	score := scoring.Calculate(currentExam, answers)
	cohort, err := s.store.Cohort(r.Context(), currentExam.Slug)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	data := pageData{
		Title: "Seu resultado", Attempt: attempt, Exam: currentExam, Answers: answers,
		Answered: len(answers), Total: len(currentExam.Questions), Score: score,
		Cohort: scoring.Standardize(score, cohort),
	}
	disciplineIndex := make(map[string]int)
	for _, question := range currentExam.Questions {
		key := fmt.Sprint(question.Number)
		correctAnswer := currentExam.Document.AnswerKey[key]
		userAnswer := answers[question.Number]
		result := questionResult{
			Number: question.Number, Discipline: question.Discipline, Paper: question.Paper,
			Statement: question.Statement, BaseText: question.BaseText, UserAnswer: userAnswer,
			CorrectAnswer: correctAnswer, Correct: userAnswer == correctAnswer, Blank: userAnswer == "",
		}
		for index, text := range question.Alternatives {
			letter := string(rune('A' + index))
			result.Alternatives = append(result.Alternatives, resultAlternative{
				Letter: letter, Text: text, Correct: letter == correctAnswer, Chosen: letter == userAnswer,
			})
		}
		data.Results = append(data.Results, result)
		index, exists := disciplineIndex[question.Discipline]
		if !exists {
			index = len(data.Disciplines)
			disciplineIndex[question.Discipline] = index
			data.Disciplines = append(data.Disciplines, disciplineResult{Name: question.Discipline})
		}
		data.Disciplines[index].Total++
		if result.Correct {
			data.Disciplines[index].Correct++
		}
	}
	s.render(w, r, "result", http.StatusOK, data)
}

func (s *Server) loadAttempt(w http.ResponseWriter, r *http.Request, sessionID string) (store.Attempt, *exam.Exam, map[int]string, bool) {
	attempt, err := s.store.GetAttempt(r.Context(), sessionID, r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		s.notFound(w, r)
		return store.Attempt{}, nil, nil, false
	}
	if err != nil {
		s.internalError(w, r, err)
		return store.Attempt{}, nil, nil, false
	}
	currentExam, err := exam.FromSnapshot(attempt.ExamJSON, attempt.ExamSlug)
	if err != nil {
		s.internalError(w, r, err)
		return store.Attempt{}, nil, nil, false
	}
	answers, err := s.store.Answers(r.Context(), attempt.ID)
	if err != nil {
		s.internalError(w, r, err)
		return store.Attempt{}, nil, nil, false
	}
	return attempt, currentExam, answers, true
}

func (s *Server) sessionID(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie("prova_session"); err == nil && len(cookie.Value) == 64 {
		return cookie.Value
	}
	value := randomID()
	http.SetCookie(w, &http.Cookie{
		Name: "prova_session", Value: value, Path: "/", MaxAge: 60 * 60 * 24 * 365,
		HttpOnly: true, Secure: s.secureCookies, SameSite: http.SameSiteLaxMode,
	})
	return value
}

func randomID() string {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buffer)
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "error", http.StatusNotFound, pageData{Title: "Página não encontrada", ErrorMessage: "Não encontramos esta página ou prova. Volte ao índice e escolha uma prova disponível."})
}

func (s *Server) internalError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Error("erro interno", "method", r.Method, "path", r.URL.Path, "error", err)
	s.render(w, r, "error", http.StatusInternalServerError, pageData{Title: "Não foi possível concluir a ação", ErrorMessage: "Tente novamente em instantes ou volte ao índice de provas."})
}

func (s *Server) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				s.internalError(w, r, fmt.Errorf("panic: %v", recovered))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
