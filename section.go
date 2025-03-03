package readme

import (
	"bytes"
	"regexp"
	"strings"
)

// Section represents a markdown section in the README
type Section struct {
	Title   string // eg. "Supported Languages"
	Content string // eg. "- en\n"
}

func (r Readme) UpdateSectionInReadmeFile(section Section) error {
	content, err := r.readFileContent(r.path)
	if err != nil {
		if !r.config.IsNotExist(err) {
			return err
		}
		return r.createNewReadme(section)
	}

	return r.updateExistingReadme(content, section)
}

func (r Readme) readFileContent(path string) ([]byte, error) {
	return r.config.ReadFile(path)
}

func (r Readme) createNewReadme(section Section) error {
	content := formatSection(section)
	return r.config.WriteFile(r.path, []byte(content), 0644)
}

func (r Readme) updateExistingReadme(content []byte, section Section) error {
	currentContent := string(content)
	existingSection := findSection(currentContent, section.Title)

	if existingSection == "" {
		return r.appendSection(currentContent, section)
	}

	if needsUpdate(existingSection, section) {
		return r.updateSection(content, section)
	}

	return nil
}

func findSection(content, title string) string {
	pattern := `(?s)## ` + title + `\n\n.*?(?:\n#|$)`
	re := regexp.MustCompile(pattern)
	return re.FindString(content)
}

func (r Readme) appendSection(content string, section Section) error {
	newContent := content + "\n" + formatSection(section)
	return r.config.WriteFile(r.path, []byte(newContent), 0644)
}

func needsUpdate(existingSection string, section Section) bool {
	return strings.TrimSpace(existingSection) != strings.TrimSpace(formatSection(section))
}

func (r Readme) updateSection(content []byte, section Section) error {
	pattern := `(?s)## ` + section.Title + `\n\n.*?(?:\n#|$)`
	re := regexp.MustCompile(pattern)
	formattedSection := formatSection(section)

	var buf bytes.Buffer
	buf.Write(re.ReplaceAll(content, []byte(formattedSection)))
	return r.config.WriteFile(r.path, buf.Bytes(), 0644)
}

func formatSection(section Section) string {
	return "## " + section.Title + "\n\n" + section.Content
}

// CreateBulletList creates a markdown bullet list from a list of items
// eg. ["item1", "item2"] -> "- item1\n- item2\n"
func CreateBulletList(items []string) string {
	var buf strings.Builder
	for _, item := range items {
		buf.WriteString("- " + item + "\n")
	}
	return buf.String()
}
