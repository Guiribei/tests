package exam

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCanonicalFixtures(t *testing.T) {
	for _, path := range []string{"../../option-one.json", "../../option-two.json", "../../option-three.json"} {
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

func TestCatalogReload(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "lote")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name string, content []byte) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	fixture, err := os.ReadFile("../../option-three.json")
	if err != nil {
		t.Fatal(err)
	}
	var catalog Catalog
	write("prova.json", fixture)
	if err := catalog.Reload(dir); err != nil {
		t.Fatal(err)
	}
	first, ok := catalog.Get("prova")
	if !ok {
		t.Fatal("prova da raiz ausente")
	}
	write("prova.json", []byte("{"))
	write("lote/prova.json", fixture)
	if err := catalog.Reload(dir); err == nil {
		t.Fatal("esperado erro no arquivo inválido")
	}
	if cached, _ := catalog.Get("prova"); cached != first {
		t.Fatal("versão válida anterior perdida")
	}
	if item, ok := catalog.Get("lote~prova"); !ok || item.SourceFile != "lote/prova.json" || len(catalog.All()) != 2 {
		t.Fatal("nova prova da subpasta ausente")
	}
	if err := os.Remove(filepath.Join(dir, "prova.json")); err != nil {
		t.Fatal(err)
	}
	if err := catalog.Reload(dir); err != nil {
		t.Fatal(err)
	}
	if _, ok := catalog.Get("prova"); ok || len(catalog.All()) != 1 {
		t.Fatal("prova removida ainda aparece no catálogo")
	}
}
