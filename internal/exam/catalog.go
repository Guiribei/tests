package exam

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Catalog struct {
	mu     sync.RWMutex
	bySlug map[string]*Exam
	all    []*Exam
}

func (c *Catalog) All() []*Exam {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.all
}

func (c *Catalog) Get(slug string) (*Exam, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.bySlug[slug]
	return item, ok
}

// Reload publishes a complete scan in one step. If an individual file is
// invalid, keep its last usable version alongside the other valid exams.
func (c *Catalog) Reload(dir string) error {
	bySlug := make(map[string]*Exam)
	var all []*Exam
	var problems []error
	err := fs.WalkDir(os.DirFS(dir), ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && name != "." && strings.HasPrefix(entry.Name(), ".") {
			return fs.SkipDir
		}
		if !entry.Type().IsRegular() || !strings.EqualFold(path.Ext(name), ".json") {
			return nil
		}
		parts := strings.Split(strings.TrimSuffix(name, path.Ext(name)), "/")
		for i, part := range parts {
			parts[i] = slugify(part)
		}
		slug := strings.Join(parts, "~")
		if other := bySlug[slug]; other != nil {
			problems = append(problems, fmt.Errorf("%s e %s têm o mesmo identificador %q", name, other.SourceFile, slug))
			return nil
		}
		item, err := LoadFile(filepath.Join(dir, filepath.FromSlash(name)))
		if err != nil {
			problems = append(problems, err)
			item, _ = c.Get(slug)
			if item == nil || item.SourceFile != name {
				return nil
			}
		} else {
			item.Slug, item.SourceFile = slug, name
		}
		bySlug[slug] = item
		all = append(all, item)
		return nil
	})
	if err != nil {
		return fmt.Errorf("listar provas em %s: %w", dir, err)
	}
	sort.SliceStable(all, func(i, j int) bool {
		return all[i].Document.Metadata.Title < all[j].Document.Metadata.Title
	})
	c.mu.Lock()
	c.bySlug, c.all = bySlug, all
	c.mu.Unlock()
	return errors.Join(problems...)
}
