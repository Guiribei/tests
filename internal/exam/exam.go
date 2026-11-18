package exam

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Metadata struct {
	Title             string      `json:"titulo"`
	Type              string      `json:"tipo"`
	RoleCode          string      `json:"codigo_cargo"`
	Notice            string      `json:"edital"`
	NoticeURL         string      `json:"edital_url"`
	ReferenceDate     string      `json:"data_de_referencia_legislativa"`
	AlternativeOrder  []string    `json:"ordem_das_alternativas"`
	Instruction       string      `json:"instrucao"`
	Papers            []PaperMeta `json:"provas"`
	DistributionNotes string      `json:"criterio_de_distribuicao"`
}

type PaperMeta struct {
	Name            string `json:"nome"`
	TotalQuestions  int    `json:"total_questoes"`
	DurationMinutes int    `json:"duracao_minutos"`
	Weight          int    `json:"peso"`
}

type BaseText struct {
	ID        string `json:"id"`
	Questions []int  `json:"questoes"`
	Text      string `json:"texto"`
}

type Section struct {
	Discipline string         `json:"disciplina"`
	Questions  []QuestionJSON `json:"questoes"`
}

type QuestionJSON struct {
	Number       int      `json:"numero"`
	Statement    string   `json:"enunciado"`
	Alternatives []string `json:"alternativas"`
}

type Document struct {
	Metadata  Metadata          `json:"metadados"`
	BaseTexts []BaseText        `json:"textos_base"`
	General   []Section         `json:"conhecimentos_gerais"`
	Specific  []Section         `json:"conhecimentos_especificos"`
	AnswerKey map[string]string `json:"gabarito"`
}

type Question struct {
	Number       int
	Statement    string
	Alternatives []string
	Discipline   string
	Paper        string
	Weight       int
	BaseTextID   string
	BaseText     string
}

type Exam struct {
	Slug       string
	SourceFile string
	Hash       string
	Document   Document
	Questions  []Question
}

func LoadFile(path string) (*Exam, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ler %s: %w", path, err)
	}
	var document Document
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("JSON inválido em %s: %w", path, err)
	}
	exam := &Exam{
		Slug:       slugify(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))),
		SourceFile: filepath.Base(path),
		Document:   document,
	}
	sum := sha256.Sum256(data)
	exam.Hash = hex.EncodeToString(sum[:8])
	if err := exam.normalizeAndValidate(); err != nil {
		return nil, fmt.Errorf("prova %s: %w", path, err)
	}
	return exam, nil
}

func FromSnapshot(data []byte, slug string) (*Exam, error) {
	var document Document
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	exam := &Exam{Slug: slug, Document: document}
	if err := exam.normalizeAndValidate(); err != nil {
		return nil, err
	}
	return exam, nil
}

func (e *Exam) Snapshot() ([]byte, error) {
	return json.Marshal(e.Document)
}

func (e *Exam) Question(number int) (Question, bool) {
	index := sort.Search(len(e.Questions), func(i int) bool { return e.Questions[i].Number >= number })
	if index >= len(e.Questions) || e.Questions[index].Number != number {
		return Question{}, false
	}
	return e.Questions[index], true
}

func (e *Exam) normalizeAndValidate() error {
	if strings.TrimSpace(e.Document.Metadata.Title) == "" {
		return errors.New("metadados.titulo é obrigatório")
	}
	if len(e.Document.General) == 0 || len(e.Document.Specific) == 0 {
		return errors.New("as duas provas, geral e específica, são obrigatórias")
	}
	if len(e.Document.Metadata.Papers) != 2 {
		return errors.New("metadados.provas deve conter exatamente duas provas")
	}

	baseByQuestion := make(map[int]BaseText)
	for _, base := range e.Document.BaseTexts {
		for _, number := range base.Questions {
			baseByQuestion[number] = base
		}
	}

	seen := make(map[int]bool)
	appendSections := func(sections []Section, paper string, weight int) error {
		for _, section := range sections {
			if strings.TrimSpace(section.Discipline) == "" {
				return errors.New("disciplina vazia")
			}
			for _, raw := range section.Questions {
				if raw.Number < 1 || seen[raw.Number] {
					return fmt.Errorf("número de questão inválido ou duplicado: %d", raw.Number)
				}
				if strings.TrimSpace(raw.Statement) == "" || len(raw.Alternatives) != 5 {
					return fmt.Errorf("questão %d deve ter enunciado e cinco alternativas", raw.Number)
				}
				seen[raw.Number] = true
				question := Question{
					Number: raw.Number, Statement: raw.Statement, Alternatives: raw.Alternatives,
					Discipline: section.Discipline, Paper: paper, Weight: weight,
				}
				if base, ok := baseByQuestion[raw.Number]; ok {
					question.BaseTextID, question.BaseText = base.ID, base.Text
				}
				e.Questions = append(e.Questions, question)
			}
		}
		return nil
	}
	if err := appendSections(e.Document.General, e.Document.Metadata.Papers[0].Name, e.Document.Metadata.Papers[0].Weight); err != nil {
		return err
	}
	if err := appendSections(e.Document.Specific, e.Document.Metadata.Papers[1].Name, e.Document.Metadata.Papers[1].Weight); err != nil {
		return err
	}
	sort.Slice(e.Questions, func(i, j int) bool { return e.Questions[i].Number < e.Questions[j].Number })

	if len(e.Questions) != 180 {
		return fmt.Errorf("esperadas 180 questões, encontradas %d", len(e.Questions))
	}
	for expected, question := range e.Questions {
		if question.Number != expected+1 {
			return fmt.Errorf("numeração deve ser contínua de 1 a 180; faltou %d", expected+1)
		}
		answer := strings.ToUpper(e.Document.AnswerKey[fmt.Sprint(question.Number)])
		if len(answer) != 1 || answer < "A" || answer > "E" {
			return fmt.Errorf("gabarito inválido para questão %d", question.Number)
		}
		e.Document.AnswerKey[fmt.Sprint(question.Number)] = answer
	}
	return nil
}

func slugify(value string) string {
	value = strings.ToLower(value)
	var result strings.Builder
	dash := false
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			result.WriteRune(r)
			dash = false
		} else if !dash && result.Len() > 0 {
			result.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(result.String(), "-")
}
