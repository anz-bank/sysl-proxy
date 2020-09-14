package sysl_proxy

import (
"bufio"
"bytes"
"path"
"regexp"
"strings"

"github.com/anz-bank/sysl/pkg/parse"

"github.com/joshcarp/gop/app"

"google.golang.org/protobuf/encoding/protojson"

"github.com/joshcarp/gop/gop"
)

const import_regex = `(?:#import.*)|(?:import )(?P<import>.*)`

type Retriever struct {
	primary      gop.Retriever
	secondary    gop.Retriever
	import_regex *regexp.Regexp
}

func New(primary, secondary gop.Retriever) Retriever {
	return Retriever{primary: primary, secondary: secondary, import_regex: regexp.MustCompile(import_regex)}
}

func (a Retriever) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	var res gop.Object
	var cached bool
	var err error
	res, cached, err = a.primary.Retrieve(repo, resource, version)
	if err == nil {
		return res, cached, err
	}
	res, cached, err = a.secondary.Retrieve(repo, resource, version)
	if err != nil {
		return res, cached, err
	}
	res.Content, err = a.ParseSysl(res)
	if err != nil {
		return res, cached, err
	}
	return res, cached, nil
}

func (a Retriever) retrieverHelper(repo, resource, version string) (gop.Object, bool, error) {
	var res gop.Object
	var cached bool
	var err error
	res, cached, err = a.primary.Retrieve(repo, resource, version)
	if err == nil {
		return res, cached, err
	}
	res, cached, err = a.secondary.Retrieve(repo, resource, version)
	if err != nil {
		return res, cached, err
	}
	return res, cached, nil
}

func (a Retriever) ParseSysl(obj gop.Object) ([]byte, error) {
	newImports := a.findImports(a.import_regex, obj.Content)
	var imports []byte
	for _, imp := range newImports {
		var repo, resource, version string
		var err error
		switch strings.Contains(imp, "//") {
		case true: // the version
			if !strings.Contains(imp, "@") {
				imp += "@master"
			}
			repo, resource, version, err = app.ProcessRequest(strings.TrimLeft(imp, "// "))
		case false:
			repo = obj.Repo
			version = obj.Version
			resource = path.Dir(obj.Repo) + "/" + imp
		}
		if err != nil { // either version or repo is missing
			return nil, err
		}
		newFile, _, err := a.retrieverHelper(repo, resource, version)
		if err != nil {
			continue
		}
		imports = append(imports, []byte("\n\n")...)
		imports = append(imports, newFile.Content...)
	}
	content := a.import_regex.ReplaceAll(obj.Content, imports)
	m, err := parse.NewParser().ParseString(string(content))
	if err != nil {
		return nil, err
	}
	content, err = protojson.Marshal(m)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (a Retriever) findImports(re *regexp.Regexp, file []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(file))
	var imports []string
	for scanner.Scan() {
		for _, match := range re.FindAllStringSubmatch(scanner.Text(), -1) {
			if match == nil {
				continue
			}
			for i, name := range re.SubexpNames() {
				if name == "import" && match[i] != "" {
					imports = append(imports, match[i])
				}
			}
		}
	}
	return imports
}
