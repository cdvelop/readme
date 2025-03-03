package readme

// Readme represents a README file
type Readme struct {
	path     string // Path to the README file default is README.md change if needed or test
	config   *Config
	sections []Section
}

// Config configuration for the Readme handler
type Config struct {
	// WriteFile writes data to a file fileMode represents the file mode and permission bits  eg. 0644, 0755
	WriteFile func(path string, data []byte, fileMode uint32) error
	// ReadFile reads the content of a file
	ReadFile func(path string) ([]byte, error)
	// IsNotExist reports whether the error is known to report that a file or directory does not exist
	IsNotExist func(err error) bool
}

func New(rc *Config) *Readme {
	return &Readme{
		path:     "README.md",
		config:   rc,
		sections: make([]Section, 0),
	}
}

// set readme file path default is README.md for test eg. temp/README.md
func (r *Readme) SetReadmePath(newPath string) {
	r.path = newPath
}

// AddSection adds a section to the cache if the title doesn't exist
func (r *Readme) AddSection(newSections ...Section) {
	for _, newSection := range newSections {
		exists := false
		for i := range r.sections {
			if r.sections[i].Title == newSection.Title {
				r.sections[i].Content = newSection.Content
				exists = true
				break
			}
		}
		if !exists {
			r.sections = append(r.sections, newSection)
		}
	}

}

// UpdateAllSectionsInReadmeFile writes all cached sections to the README file
func (r *Readme) UpdateAllSectionsInReadmeFile() error {
	for _, section := range r.sections {
		if err := r.UpdateSectionInReadmeFile(section); err != nil {
			return err
		}
	}
	return nil
}

// GetSections returns all cached sections
func (r *Readme) GetSections() []Section {
	return r.sections
}
