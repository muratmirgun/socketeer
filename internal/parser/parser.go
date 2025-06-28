package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"encoding/json"

	"bufio"

	"github.com/muratmirgun/socketeer/internal/spec"
	"gopkg.in/yaml.v3"
)

// Parse scans Go source files for WebSocket annotations and returns a list of Socket specs.
func Parse(dir string) ([]*spec.Socket, error) {
	var sockets []*spec.Socket
	fset := token.NewFileSet()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		println("Visiting:", path)

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// Build a list of all function positions
		funcs := []*ast.FuncDecl{}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				funcs = append(funcs, fn)
			}
		}

		// Map: function name -> all annotation blocks
		funcAnnots := map[string][][]string{}

		// For each comment group, find the first function that follows it
		for _, cg := range file.Comments {
			if len(cg.List) == 0 {
				continue
			}
			block := extractAnnotationBlock(cg.List)
			if len(block) == 0 {
				continue
			}
			cgEnd := cg.End()
			for _, fn := range funcs {
				if fn.Pos() > cgEnd {
					funcAnnots[fn.Name.Name] = append(funcAnnots[fn.Name.Name], block)
					break
				}
			}
		}

		// For each function, merge all annotation blocks and parse as a single socket
		for fnName, blocks := range funcAnnots {
			var merged []string
			for _, b := range blocks {
				merged = append(merged, b...)
			}
			println("Function:", fnName)
			println("Merged annotation block:", strings.Join(merged, " | "))
			if isWebSocketBlock(merged) {
				socket := parseSocketBlock(merged)
				sockets = append(sockets, socket)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sockets, nil
}

// extractAnnotationBlock extracts all consecutive annotation lines from a comment group, including multi-blocks.
func extractAnnotationBlock(comments []*ast.Comment) []string {
	var block []string
	for _, c := range comments {
		line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if strings.HasPrefix(line, "@") || (len(block) > 0 && (strings.HasPrefix(line, "{") || strings.HasPrefix(line, "}"))) {
			block = append(block, line)
		}
	}
	return block
}

// isWebSocketBlock checks if the annotation block starts with @WebSocket.
func isWebSocketBlock(block []string) bool {
	for _, line := range block {
		if strings.HasPrefix(line, "@WebSocket") {
			return true
		}
	}
	return false
}

// StructInfo holds struct field info for JSON example generation
type StructInfo struct {
	Fields map[string]interface{}
}

// CollectStructs walks all Go files under dir and returns a map of struct name (with/without package) to StructInfo
func CollectStructs(dir string) map[string]StructInfo {
	structs := map[string]StructInfo{}
	fset := token.NewFileSet()
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}
		pkg := file.Name.Name
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, specNode := range gd.Specs {
				ts, ok := specNode.(*ast.TypeSpec)
				st, ok2 := ts.Type.(*ast.StructType)
				if !ok || !ok2 {
					continue
				}
				fields := map[string]interface{}{}
				for _, f := range st.Fields.List {
					name := ""
					if len(f.Names) > 0 {
						name = f.Names[0].Name
					}
					jsonName := name
					if f.Tag != nil {
						tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
						if j, ok := tag.Lookup("json"); ok && j != "-" {
							jsonName = strings.Split(j, ",")[0]
						}
					}
					if jsonName == "" || jsonName == "-" {
						continue
					}
					// Basit örnek değerler
					var example interface{} = ""
					switch ft := f.Type.(type) {
					case *ast.Ident:
						t := ft.Name
						if t == "string" {
							example = "string"
						} else if t == "int" || t == "int64" || t == "int32" {
							example = 0
						} else if t == "bool" {
							example = false
						} else {
							example = t
						}
					case *ast.ArrayType:
						example = []interface{}{}
					case *ast.MapType:
						example = map[string]interface{}{}
					}
					fields[jsonName] = example
				}
				structs[ts.Name.Name] = StructInfo{Fields: fields}
				structs[pkg+"."+ts.Name.Name] = StructInfo{Fields: fields}
			}
		}
		return nil
	})
	return structs
}

// parseSocketBlock parses a block of annotations into a Socket struct (stub for now).
func parseSocketBlock(block []string) *spec.Socket {
	socket := &spec.Socket{}
	var currentMsg *spec.Message
	var inPayload, inExample bool
	var payloadLines, exampleLines []string

	// struct haritasını bir kez topla
	structMap := CollectStructs("./")

	for i := 0; i < len(block); i++ {
		line := block[i]
		if inPayload {
			if strings.HasPrefix(line, "@") && !strings.HasPrefix(line, "@Payload") {
				inPayload = false
				payload := parseJSONBlock(payloadLines)
				if currentMsg != nil {
					currentMsg.Payload = payload
				}
				payloadLines = nil
				// Continue to parse this line as a new annotation
			} else {
				payloadLines = append(payloadLines, strings.TrimPrefix(line, "// "))
				continue
			}
		}
		if inExample {
			if strings.HasPrefix(line, "@") && !strings.HasPrefix(line, "@Example") {
				inExample = false
				example := parseJSONBlock(exampleLines)
				if currentMsg != nil {
					currentMsg.Example = example
				}
				exampleLines = nil
				// Continue to parse this line as a new annotation
			} else {
				exampleLines = append(exampleLines, strings.TrimPrefix(line, "// "))
				continue
			}
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "@WebSocket":
			if len(fields) > 1 {
				socket.Name = fields[1]
			}
		case "@Group":
			if len(fields) > 1 {
				socket.Group = strings.Join(fields[1:], " ")
			}
		case "@URL":
			if len(fields) > 1 {
				socket.URL = fields[1]
			}
		case "@Description":
			desc := strings.TrimPrefix(line, "@Description")
			desc = strings.TrimSpace(desc)
			if currentMsg != nil {
				currentMsg.Description = desc
			} else {
				socket.Description = desc
			}
		case "@Tags":
			if currentMsg != nil && len(fields) > 1 {
				tags := strings.Split(strings.Join(fields[1:], " "), ",")
				for i, t := range tags {
					tags[i] = strings.TrimSpace(t)
				}
				currentMsg.Tags = tags
			} else if currentMsg == nil && len(fields) > 1 {
				tags := strings.Split(strings.Join(fields[1:], " "), ",")
				for i, t := range tags {
					tags[i] = strings.TrimSpace(t)
				}
				socket.Tags = tags
			}
		case "@ConnectionParam":
			if len(fields) >= 5 {
				param := spec.ConnectionParam{
					Name:     fields[1],
					In:       fields[2],
					Type:     fields[3],
					Required: fields[4] == "required",
				}
				if len(fields) > 5 {
					param.Description = strings.Join(fields[5:], " ")
				}
				socket.ConnectionParams = append(socket.ConnectionParams, param)
			}
		case "@Message":
			if currentMsg != nil {
				socket.Messages = append(socket.Messages, *currentMsg)
			}
			currentMsg = &spec.Message{Type: ""}
			if len(fields) > 1 {
				currentMsg.Type = fields[1]
			}
		case "@Direction":
			if currentMsg != nil && len(fields) > 1 {
				currentMsg.Direction = fields[1]
			}
		case "@Payload":
			// Eğer struct ismi verilmişse, onu kullan
			if len(fields) == 2 {
				structName := fields[1]
				if s, ok := structMap[structName]; ok {
					if currentMsg != nil {
						currentMsg.Payload = s.Fields
					}
					inPayload = false
					continue
				}
			}
			// Eski JSON block mantığı
			inPayload = true
			payloadLines = nil
		case "@Example":
			inExample = true
			exampleLines = nil
		case "@Error":
			if currentMsg != nil && len(fields) > 1 {
				err := spec.Error{Code: fields[1]}
				if len(fields) > 2 {
					err.Description = strings.Join(fields[2:], " ")
				}
				currentMsg.Errors = append(currentMsg.Errors, err)
			}
		case "@Deprecated":
			if currentMsg != nil {
				currentMsg.Deprecated = true
			}
		}
	}
	if currentMsg != nil {
		socket.Messages = append(socket.Messages, *currentMsg)
	}
	return socket
}

// parseJSONBlock joins lines and parses JSON, returns map or array or string.
func parseJSONBlock(lines []string) interface{} {
	joined := strings.Join(lines, "\n")
	var v interface{}
	if err := json.Unmarshal([]byte(joined), &v); err != nil {
		return joined // fallback: raw string
	}
	return v
}

// ParseInfoAnnotations scans Go files in srcDir for top-level API info annotations.
func ParseInfoAnnotations(srcDir string) spec.Info {
	info := spec.Info{}
	filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()
		lines := []string{}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		// 1. Dosya başındaki consecutive annotation'ları tara
		for _, line := range lines {
			l := strings.TrimSpace(strings.TrimPrefix(line, "//"))
			if !strings.HasPrefix(l, "@") {
				break
			}
			parseInfoLine(&info, l)
		}
		// 2. main fonksiyonu üstündeki consecutive annotation block'u tara
		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], "func main(") {
				// Yukarıya doğru consecutive //@ annotation'ları topla
				for j := i - 1; j >= 0; j-- {
					l := strings.TrimSpace(strings.TrimPrefix(lines[j], "//"))
					if strings.HasPrefix(l, "@") {
						parseInfoLine(&info, l)
					} else if l == "" {
						continue
					} else {
						break
					}
				}
				break
			}
		}
		return nil
	})
	return info
}

// parseInfoLine yardımcı fonksiyonu
func parseInfoLine(info *spec.Info, line string) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return
	}
	switch fields[0] {
	case "@title":
		info.Title = strings.Join(fields[1:], " ")
	case "@version":
		info.Version = strings.Join(fields[1:], " ")
	case "@description":
		info.Description = strings.Join(fields[1:], " ")
	case "@contact.name":
		info.Contact.Name = strings.Join(fields[1:], " ")
	case "@contact.email":
		info.Contact.Email = strings.Join(fields[1:], " ")
	case "@license.name":
		info.License.Name = strings.Join(fields[1:], " ")
	case "@license.url":
		info.License.URL = strings.Join(fields[1:], " ")
	}
}

// ParseAndWriteSpec parses Go files in srcDir and writes the spec to outFile (YAML).
func ParseAndWriteSpec(srcDir, outFile string) error {
	sockets, err := Parse(srcDir)
	if err != nil {
		return err
	}
	info := ParseInfoAnnotations(srcDir)
	if info.Title == "" {
		info.Title = "WebSocket API"
	}
	if info.Version == "" {
		info.Version = "1.0.0"
	}
	if info.Description == "" {
		info.Description = "Generated by wsdoc"
	}
	s := spec.Spec{
		Info:    info,
		Sockets: []spec.Socket{},
	}
	for _, sock := range sockets {
		s.Sockets = append(s.Sockets, *sock)
	}

	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(&s); err != nil {
		return err
	}
	return nil
}
