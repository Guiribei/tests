package exam

import "testing"

func TestLoadCanonicalFixtures(t *testing.T) {
	for _, path := range []string{"../../option-one.json", "../../option-two.json"} {
		exam, err := LoadFile(path)
		if err != nil {
			t.Fatalf("LoadFile(%s): %v", path, err)
		}
		if len(exam.Questions) != 180 {
			t.Fatalf("%s: got %d questions", path, len(exam.Questions))
		}
		if len(exam.Document.AnswerKey) != 180 {
			t.Fatalf("%s: got %d answers", path, len(exam.Document.AnswerKey))
		}
	}
}
