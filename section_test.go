package readme_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cdvelop/readme"
)

func TestUpdateReadmeLangSection(t *testing.T) {
	// Create a temporary README.md for testing
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

	tests := []struct {
		name           string
		initialContent string
		languages      []string
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:           "Should create README with Spanish and English languages",
			initialContent: "",
			languages:      []string{"es", "en"},
			wantContains: []string{
				"## Supported Languages",
				"- es",
				"- en",
			},
		},
		{
			name:           "Should preserve project info and add languages section",
			initialContent: "# My Project\n\nMultilanguage project for data management\n",
			languages:      []string{"es", "en", "fr"},
			wantContains: []string{
				"# My Project",
				"Multilanguage project",
				"## Supported Languages",
				"- es",
				"- en",
				"- fr",
			},
		},
		{
			name:           "Should update outdated languages list",
			initialContent: "## Supported Languages\n\n- fr\n- it\n",
			languages:      []string{"es", "en"},
			wantContains: []string{
				"## Supported Languages",
				"- es",
				"- en",
			},
			wantNotContain: []string{
				"- fr",
				"- it",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create initial content if any
			if tt.initialContent != "" {
				err := os.WriteFile(readmePath, []byte(tt.initialContent), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Create section with language list
			section := readme.Section{
				Title:   "Supported Languages",
				Content: readme.CreateBulletList(tt.languages),
			}

			// Update README
			err := r.UpdateSectionInReadmeFile(section)
			if err != nil {
				t.Fatal(err)
			}

			// Read updated content
			content, err := os.ReadFile(readmePath)
			if err != nil {
				t.Fatal(err)
			}
			updatedContent := string(content)

			// Check expected content is present
			for _, want := range tt.wantContains {
				if !strings.Contains(updatedContent, want) {
					t.Errorf("\nTest: %s\nContent:\n%v\nShould contain: %v",
						tt.name, updatedContent, want)
				}
			}

			// Check unwanted content is not present
			for _, unwanted := range tt.wantNotContain {
				if strings.Contains(updatedContent, unwanted) {
					t.Errorf("\nTest: %s\nContent:\n%v\nShould not contain: %v",
						tt.name, updatedContent, unwanted)
				}
			}

			// Cleanup
			if err := os.Remove(readmePath); err != nil {
				t.Fatal(err)
			}
		})
	}
}
