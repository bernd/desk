package deskfile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

type CallableType string

const (
	TypeAlias    CallableType = "alias"
	TypeFunction CallableType = "function"
	TypeExport   CallableType = "export"
)

type Callable struct {
	Name string
	Doc  string
	Type CallableType
}

type Deskfile struct {
	Path        string
	Name        string
	Description string
	Callables   []Callable
}

func Parse(path string) (*Deskfile, error) {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)
	
	if ext == ".fish" {
		return parseFish(path, name)
	}
	return parseBash(path, name)
}

func parseBash(path, name string) (*Deskfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parser := syntax.NewParser(syntax.KeepComments(true))
	file, err := parser.Parse(strings.NewReader(string(data)), path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	desk := &Deskfile{
		Path:      path,
		Name:      name,
		Callables: []Callable{},
	}

	lines := strings.Split(string(data), "\n")
	descRe := regexp.MustCompile(`#\s*Description:\s*(.+)`)
	for _, line := range lines {
		if match := descRe.FindStringSubmatch(line); match != nil {
			desk.Description = match[1]
			break
		}
	}

	type nodeInfo struct {
		node syntax.Node
		line uint
	}

	var allNodes []nodeInfo
	syntax.Walk(file, func(node syntax.Node) bool {
		if node != nil {
			allNodes = append(allNodes, nodeInfo{node: node, line: node.Pos().Line()})
		}
		return true
	})

	commentMap := make(map[uint]string)
	for _, info := range allNodes {
		if comment, ok := info.node.(*syntax.Comment); ok {
			text := strings.TrimSpace(comment.Text)
			text = strings.TrimPrefix(text, "#")
			text = strings.TrimSpace(text)
			if !strings.HasPrefix(text, "Description:") {
				commentMap[info.line] = text
			}
		}
	}

	findDocComment := func(targetLine uint) string {
		for i := int(targetLine) - 1; i >= 1 && i >= int(targetLine)-5; i-- {
			if doc, ok := commentMap[uint(i)]; ok {
				return doc
			}
		}
		return ""
	}

	for _, stmt := range file.Stmts {
		line := stmt.Pos().Line()
		doc := findDocComment(line)

		switch cmd := stmt.Cmd.(type) {
		case *syntax.FuncDecl:
			desk.Callables = append(desk.Callables, Callable{
				Name: cmd.Name.Value,
				Doc:  doc,
				Type: TypeFunction,
			})

		case *syntax.CallExpr:
			if len(cmd.Args) > 0 {
				firstWord := getFirstLiteral(cmd.Args[0])
				
				if firstWord == "alias" && len(cmd.Args) > 1 {
					aliasArg := getFirstLiteral(cmd.Args[1])
					aliasName := strings.SplitN(aliasArg, "=", 2)[0]

					desk.Callables = append(desk.Callables, Callable{
						Name: aliasName,
						Doc:  doc,
						Type: TypeAlias,
					})
				}
			}

			for _, assign := range cmd.Assigns {
				if assign.Name != nil {
					desk.Callables = append(desk.Callables, Callable{
						Name: assign.Name.Value,
						Doc:  doc,
						Type: TypeExport,
					})
				}
			}

		case *syntax.DeclClause:
			for _, assign := range cmd.Args {
				if assign != nil && assign.Name != nil {
					desk.Callables = append(desk.Callables, Callable{
						Name: assign.Name.Value,
						Doc:  doc,
						Type: TypeExport,
					})
				}
			}
		}
	}

	return desk, nil
}

func parseFish(path, name string) (*Deskfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	desk := &Deskfile{
		Path:      path,
		Name:      name,
		Callables: []Callable{},
	}

	descRe := regexp.MustCompile(`#\s*Description:\s*(.+)`)
	funcRe := regexp.MustCompile(`^\s*function\s+(\w+)`)

	var lastComment string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if match := descRe.FindStringSubmatch(line); match != nil {
			desk.Description = match[1]
			continue
		}

		if strings.HasPrefix(line, "#") {
			comment := strings.TrimPrefix(line, "#")
			comment = strings.TrimSpace(comment)
			if !strings.HasPrefix(comment, "Description:") {
				lastComment = comment
			}
			continue
		}

		if match := funcRe.FindStringSubmatch(line); match != nil {
			funcName := match[1]
			doc := lastComment
			lastComment = ""

			desk.Callables = append(desk.Callables, Callable{
				Name: funcName,
				Doc:  doc,
				Type: TypeFunction,
			})
			continue
		}

		if line != "" && !strings.HasPrefix(line, "#") {
			lastComment = ""
		}
	}

	return desk, nil
}

func getFirstLiteral(word *syntax.Word) string {
	if word == nil || len(word.Parts) == 0 {
		return ""
	}
	
	if lit, ok := word.Parts[0].(*syntax.Lit); ok {
		return lit.Value
	}
	
	var result strings.Builder
	for _, part := range word.Parts {
		switch p := part.(type) {
		case *syntax.Lit:
			result.WriteString(p.Value)
		case *syntax.SglQuoted:
			result.WriteString(p.Value)
		case *syntax.DblQuoted:
			for _, qpart := range p.Parts {
				if lit, ok := qpart.(*syntax.Lit); ok {
					result.WriteString(lit.Value)
				}
			}
		}
	}
	return result.String()
}

func Find(desksDir, deskName, ext string) (string, error) {
	deskPath := filepath.Join(desksDir, deskName+ext)
	
	if _, err := os.Stat(deskPath); err == nil {
		return deskPath, nil
	}

	return "", fmt.Errorf("desk not found: %s", deskName)
}

func List(desksDir string) ([]*Deskfile, error) {
	entries, err := os.ReadDir(desksDir)
	if err != nil {
		return nil, err
	}

	deskMap := make(map[string]*Deskfile)
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".sh" && ext != ".fish" {
			continue
		}

		path := filepath.Join(desksDir, entry.Name())
		desk, err := Parse(path)
		if err != nil {
			continue
		}

		if _, exists := deskMap[desk.Name]; !exists {
			deskMap[desk.Name] = desk
		}
	}

	var desks []*Deskfile
	for _, desk := range deskMap {
		desks = append(desks, desk)
	}

	sort.Slice(desks, func(i, j int) bool {
		return desks[i].Name < desks[j].Name
	})

	return desks, nil
}
