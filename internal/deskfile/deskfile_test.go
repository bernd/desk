package deskfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseBash(t *testing.T) {
	content := `# Description: test desk for testing
cd /tmp

# A test function
test_func() {
  echo "test"
}

# A test alias
alias t="test"

# An exported variable
export TEST_VAR="value"
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.sh")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	desk, err := Parse(testFile)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if desk.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", desk.Name)
	}

	if desk.Description != "test desk for testing" {
		t.Errorf("Expected description 'test desk for testing', got '%s'", desk.Description)
	}

	if len(desk.Callables) != 3 {
		t.Fatalf("Expected 3 callables, got %d", len(desk.Callables))
	}

	expectedCallables := map[string]struct {
		typ CallableType
		doc string
	}{
		"test_func": {TypeFunction, "A test function"},
		"t":         {TypeAlias, "A test alias"},
		"TEST_VAR":  {TypeExport, "An exported variable"},
	}

	for _, callable := range desk.Callables {
		expected, ok := expectedCallables[callable.Name]
		if !ok {
			t.Errorf("Unexpected callable: %s", callable.Name)
			continue
		}

		if callable.Type != expected.typ {
			t.Errorf("Callable %s: expected type %s, got %s", callable.Name, expected.typ, callable.Type)
		}

		if callable.Doc != expected.doc {
			t.Errorf("Callable %s: expected doc '%s', got '%s'", callable.Name, expected.doc, callable.Doc)
		}
	}
}

func TestParseFish(t *testing.T) {
	content := `# Description: fish test desk

# Say hello
function hello
    echo hello
end
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.fish")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	desk, err := Parse(testFile)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if desk.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", desk.Name)
	}

	if desk.Description != "fish test desk" {
		t.Errorf("Expected description 'fish test desk', got '%s'", desk.Description)
	}

	if len(desk.Callables) != 1 {
		t.Fatalf("Expected 1 callable, got %d", len(desk.Callables))
	}

	if desk.Callables[0].Name != "hello" {
		t.Errorf("Expected callable 'hello', got '%s'", desk.Callables[0].Name)
	}

	if desk.Callables[0].Doc != "Say hello" {
		t.Errorf("Expected doc 'Say hello', got '%s'", desk.Callables[0].Doc)
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()

	testDesks := map[string]string{
		"desk1.sh": `# Description: first desk
echo test`,
		"desk2.sh": `# Description: second desk
echo test`,
		"desk3.fish": `# Description: third desk
echo test`,
	}

	for name, content := range testDesks {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	desks, err := List(tmpDir)
	if err != nil {
		t.Fatalf("Failed to list desks: %v", err)
	}

	if len(desks) != 3 {
		t.Fatalf("Expected 3 desks, got %d", len(desks))
	}

	for i, desk := range desks {
		expectedName := filepath.Base(desk.Path)
		expectedName = expectedName[:len(expectedName)-len(filepath.Ext(expectedName))]
		if desk.Name != expectedName {
			t.Errorf("Desk %d: expected name '%s', got '%s'", i, expectedName, desk.Name)
		}
	}
}
