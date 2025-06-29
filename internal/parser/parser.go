package parser

import (
	"bufio"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

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
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
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
					// Try to extract Example comment
					example := ""
					if f.Doc != nil {
						for _, c := range f.Doc.List {
							line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
							if strings.HasPrefix(line, "Example:") {
								ex := strings.TrimSpace(strings.TrimPrefix(line, "Example:"))
								if len(ex) > 0 {
									example = ex
								}
							}
						}
					}
					if example == "" && f.Comment != nil {
						for _, c := range f.Comment.List {
							line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
							if strings.HasPrefix(line, "Example:") {
								ex := strings.TrimSpace(strings.TrimPrefix(line, "Example:"))
								if len(ex) > 0 {
									example = ex
								}
							}
						}
					}
					if example != "" {
						// Remove leading and trailing double quotes if present
						if strings.HasPrefix(example, "\"") && strings.HasSuffix(example, "\"") && len(example) > 1 {
							example = strings.TrimPrefix(example, "\"")
							example = strings.TrimSuffix(example, "\"")
						}
						// Try to parse as int, float, or JSON, fallback to string
						var v interface{} = example
						if i, err := strconv.ParseInt(example, 10, 64); err == nil {
							v = i
						} else if f, err := strconv.ParseFloat(example, 64); err == nil {
							v = f
						} else if (strings.HasPrefix(example, "{") && strings.HasSuffix(example, "}")) || (strings.HasPrefix(example, "[") && strings.HasSuffix(example, "]")) {
							var j interface{}
							if err := json.Unmarshal([]byte(example), &j); err == nil {
								v = j
							}
						}
						fields[jsonName] = v
					} else {
						// Fallback: old logic
						switch ft := f.Type.(type) {
						case *ast.Ident:
							t := ft.Name
							if t == "string" {
								fields[jsonName] = "string"
							} else if t == "int" || t == "int64" || t == "int32" {
								fields[jsonName] = 0
							} else if t == "bool" {
								fields[jsonName] = false
							} else {
								fields[jsonName] = t
							}
						case *ast.ArrayType:
							fields[jsonName] = []interface{}{}
						case *ast.MapType:
							fields[jsonName] = map[string]interface{}{}
						}
					}
				}
				structs[ts.Name.Name] = StructInfo{Fields: fields}
				structs[pkg+"."+ts.Name.Name] = StructInfo{Fields: fields}
			}
		}
		return nil
	})
	return structs
}

// parseSocketBlock parses a block of annotations into a Socket struct (supports grouped @Send/@Receive).
func parseSocketBlock(block []string) *spec.Socket {
	socket := &spec.Socket{}
	structMap := CollectStructs("./")
	messageGroups := make(map[string]*spec.GroupedMessage)

	var currentMsg *spec.Message
	var currentType string
	var currentDirection string

	for i := 0; i < len(block); i++ {
		line := block[i]
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
			// Add previous message to grouped structure if exists
			if currentMsg != nil && currentType != "" && currentDirection != "" {
				if currentDirection == "send" {
					messageGroups[currentType].Send = currentMsg
					if messageGroups[currentType].Description == "" {
						messageGroups[currentType].Description = currentMsg.Description
					}
				} else if currentDirection == "receive" {
					messageGroups[currentType].Receive = currentMsg
					if messageGroups[currentType].Description == "" {
						messageGroups[currentType].Description = currentMsg.Description
					}
				}
			}
			currentType = ""
			if len(fields) > 1 {
				currentType = fields[1]
			}
			if messageGroups[currentType] == nil {
				messageGroups[currentType] = &spec.GroupedMessage{Type: currentType}
			}
		case "@Send":
			// Add previous message to grouped structure if exists
			if currentMsg != nil && currentType != "" && currentDirection != "" {
				if currentDirection == "send" {
					messageGroups[currentType].Send = currentMsg
					if messageGroups[currentType].Description == "" {
						messageGroups[currentType].Description = currentMsg.Description
					}
				} else if currentDirection == "receive" {
					messageGroups[currentType].Receive = currentMsg
					if messageGroups[currentType].Description == "" {
						messageGroups[currentType].Description = currentMsg.Description
					}
				}
			}
			currentDirection = "send"
			currentMsg = &spec.Message{Type: currentType, Direction: "send"}
		case "@Receive":
			// Add previous message to grouped structure if exists
			if currentMsg != nil && currentType != "" && currentDirection != "" {
				if currentDirection == "send" {
					messageGroups[currentType].Send = currentMsg
					if messageGroups[currentType].Description == "" {
						messageGroups[currentType].Description = currentMsg.Description
					}
				} else if currentDirection == "receive" {
					messageGroups[currentType].Receive = currentMsg
					if messageGroups[currentType].Description == "" {
						messageGroups[currentType].Description = currentMsg.Description
					}
				}
			}
			currentDirection = "receive"
			currentMsg = &spec.Message{Type: currentType, Direction: "receive"}
		case "@Payload":
			if currentMsg != nil {
				payloadArg := strings.TrimSpace(strings.TrimPrefix(line, "@Payload"))
				if payloadArg != "" {
					// Check if it's inline JSON (starts with { or [)
					if strings.HasPrefix(payloadArg, "{") || strings.HasPrefix(payloadArg, "[") {
						// Parse inline JSON
						var jsonPayload interface{}
						if err := json.Unmarshal([]byte(payloadArg), &jsonPayload); err == nil {
							if b, err := json.Marshal(jsonPayload); err == nil {
								currentMsg.Payload = string(b)
							} else {
								currentMsg.Payload = jsonPayload
							}
						} else {
							currentMsg.Payload = payloadArg // fallback to raw string
						}
					} else {
						// Treat as struct name
						structName := payloadArg
						found := false
						if s, ok := structMap[structName]; ok {
							if s.Fields != nil {
								if b, err := json.Marshal(s.Fields); err == nil {
									currentMsg.Payload = string(b)
								} else {
									currentMsg.Payload = s.Fields // fallback
								}
							}
							found = true
						}
						if !found {
							for k, s := range structMap {
								if strings.HasSuffix(k, "."+structName) {
									if s.Fields != nil {
										if b, err := json.Marshal(s.Fields); err == nil {
											currentMsg.Payload = string(b)
										} else {
											currentMsg.Payload = s.Fields
										}
									}
									found = true
									break
								}
							}
						}
						// If not found, do not attempt fuzzy/contains match. Leave payload empty or log warning.
					}
				}
			}
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
	// Add the last message if exists
	if currentMsg != nil && currentType != "" && currentDirection != "" {
		if currentDirection == "send" {
			messageGroups[currentType].Send = currentMsg
			if messageGroups[currentType].Description == "" {
				messageGroups[currentType].Description = currentMsg.Description
			}
		} else if currentDirection == "receive" {
			messageGroups[currentType].Receive = currentMsg
			if messageGroups[currentType].Description == "" {
				messageGroups[currentType].Description = currentMsg.Description
			}
		}
	}
	// Convert grouped messages to slice and add to socket
	for _, groupedMsg := range messageGroups {
		socket.GroupedMessages = append(socket.GroupedMessages, *groupedMsg)
	}
	// For backward compatibility, also populate the old Messages field
	for _, groupedMsg := range socket.GroupedMessages {
		if groupedMsg.Send != nil {
			socket.Messages = append(socket.Messages, *groupedMsg.Send)
		}
		if groupedMsg.Receive != nil {
			socket.Messages = append(socket.Messages, *groupedMsg.Receive)
		}
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
