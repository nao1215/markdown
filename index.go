package markdown

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
)

// IndexOption are options for generating an index.
type IndexOption func(*Index) error

// WithTitle sets the title of the index.
func WithTitle(title string) IndexOption {
	return func(i *Index) error {
		i.title = title
		return nil
	}
}

// WithDescription sets the description of the index.
func WithDescription(description []string) IndexOption {
	return func(i *Index) error {
		i.description = description
		return nil
	}
}

// WithWriter sets the writer to write the index to.
func WithWriter(w io.Writer) IndexOption {
	return func(i *Index) error {
		i.w = w
		return nil
	}
}

// Index represents an Index of all markdown files in a directory.
type Index struct {
	// w is the writer to write the index to.
	// Default is file that locatted in the targetDir/index.md.
	w io.Writer
	// targetDir is the directory to index.
	targetDir string
	// dir is a list of directories and their files.
	dir []*dir
	// title is the title of the index markdown.
	title string
	// description is the description of the index markdown.
	description []string
}

// newIndex creates a new index.
func newIndex(targetDir string, opts ...IndexOption) (*Index, error) {
	index := &Index{
		targetDir: targetDir,
	}

	for _, opt := range opts {
		if err := opt(index); err != nil {
			return nil, err
		}
	}
	return index, nil
}

// walk walks the target directory and creates an index.
func (i *Index) walk() error {
	err := godirwalk.Walk(i.targetDir, &godirwalk.Options{
		Callback: func(path string, dirent *godirwalk.Dirent) error {
			if dirent.IsDir() {
				if !i.hasDir(path) {
					i.appendDir(path)
				}
				return nil
			}

			if !strings.Contains(path, ".md") {
				return godirwalk.SkipThis
			}
			i.appendFile(path)
			return nil
		},
		Unsorted: false,
	})
	return err
}

// write writes the index to the provided io.Writer.
func (i *Index) write() (err error) {
	if i.w == nil {
		const readUserOnly = 0600
		if i.w, err = os.OpenFile(filepath.Clean(filepath.Join(i.targetDir, "index.md")), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, readUserOnly); err != nil {
			return err
		}
	}

	markdown := NewMarkdown(i.w)
	if i.title != "" {
		markdown.H2(i.title)
	}
	if len(i.description) != 0 {
		for _, d := range i.description {
			markdown.PlainText(d).LF()
		}
	}

	for _, d := range i.dir {
		if len(d.files) == 0 {
			continue
		}
		if len(d.files) == 1 && d.files[0] == filepath.Join(i.targetDir, "index.md") {
			continue
		}

		subTitle := filepath.Base(d.path)
		if subTitle == "." {
			subTitle = "top"
		}
		markdown.H3(subTitle)

		for _, f := range d.files {
			if h1 := firstH1orH2(f); h1 != "" {
				markdown.BulletList(Link(h1, strings.Replace(f, i.targetDir+string(filepath.Separator), "", 1)))
				continue
			}
			markdown.BulletList(Link(filepath.Base(f), strings.Replace(f, i.targetDir+string(filepath.Separator), "", 1)))
		}
		markdown.LF()
	}
	return markdown.Build()
}

// hasDir checks if the index has a directory.
func (i *Index) hasDir(dir string) bool {
	for _, d := range i.dir {
		if d.path == dir {
			return true
		}
	}
	return false
}

// appendDir appends a directory to the index.
func (i *Index) appendDir(path string) {
	i.dir = append(i.dir, &dir{
		path: path,
	})
}

// appendFile appends a file to the index.
func (i *Index) appendFile(filePath string) {
	dirPath := filepath.Dir(filePath)

	for _, d := range i.dir {
		if d.path == dirPath {
			d.files = append(d.files, filePath)
			return
		}
	}
}

// dir represents a directory and its files.
type dir struct {
	// path is the path to the directory.
	path string
	// files is a list of files in the directory.
	files []string
}

// GenerateIndex generates an index of all markdown files in the target directory.
// The index is written to the provided io.Writer.
func GenerateIndex(targetDir string, opts ...IndexOption) error {
	i, err := newIndex(targetDir, opts...)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInitMarkdownIndex, err.Error())
	}

	if err := i.walk(); err != nil {
		return fmt.Errorf("%w: %s", ErrCreateMarkdownIndex, err.Error())
	}

	if err := i.write(); err != nil {
		return fmt.Errorf("%w: %s", ErrWriteMarkdownIndex, err.Error())
	}
	return nil
}

// firstH1orH2 returns the first H1 or H2 of the markdown file.
// For example, if the file has the following content:
//
//	# H1
//	## H2
//
// above pattern, it returns "H1".
// If the file doesn't have H1(H2) or error happen, it returns an empty string.
func firstH1orH2(markdownFilePath string) string {
	file, err := os.Open(filepath.Clean(markdownFilePath))
	if err != nil {
		return ""
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close file: %s", err.Error())
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[1:])
		}
		if strings.HasPrefix(line, "## ") {
			return strings.TrimSpace(line[2:])
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}
