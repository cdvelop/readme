package readme_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cdvelop/readme"
)

func TestCachedSections(t *testing.T) {
	tmpDir := t.TempDir()
	readmePath := filepath.Join(tmpDir, "README.md")

	r := readme.New(&readme.Config{
		WriteFile: func(path string, data []byte, fileMode uint32) error {
			return os.WriteFile(path, data, os.FileMode(fileMode))
		},
		ReadFile:   os.ReadFile,
		IsNotExist: os.IsNotExist,
	})

	r.SetReadmePath(readmePath)

	// Define test sections
	sections := []readme.Section{
		{
			Title:   "Languages",
			Content: "- en\n- es\n",
		},
		{
			Title:   "Features",
			Content: "- Feature 1\n- Feature 2\n",
		},
	}

	// Add sections to cache
	r.AddSection(sections...)

	// Verify sections were cached
	cached := r.GetSections()
	if len(cached) != len(sections) {
		t.Errorf("expected %d sections, got %d", len(sections), len(cached))
	}

	// Test writing all sections
	err := r.UpdateAllSectionsInReadmeFile()
	if err != nil {
		t.Fatal(err)
	}

	// Read the file and verify content
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}

	fileContent := string(content)

	// Verify all sections are present
	for _, section := range sections {
		if !strings.Contains(fileContent, section.Title) {
			t.Errorf("section title %q not found in file", section.Title)
		}
		if !strings.Contains(fileContent, section.Content) {
			t.Errorf("section content %q not found in file", section.Content)
		}
	}

	// Test updating existing sections
	updatedSection := readme.Section{
		Title:   "Languages",
		Content: "- en\n- es\n- fr\n",
	}
	r.AddSection(updatedSection)

	err = r.UpdateAllSectionsInReadmeFile()
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}

	updatedContent := string(content)
	if !strings.Contains(updatedContent, "- fr") {
		t.Error("updated content not found in file")
	}
}
