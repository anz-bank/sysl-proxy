package sysl_proxy

import (
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/joshcarp/gop/gop"
	"google.golang.org/protobuf/encoding/protojson"
	"regexp"
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

func (a Retriever) Retrieve(resource string) ([]byte, bool, error) {
	var res []byte
	var cached bool
	var err error
	var m *sysl.Module
	res, cached, err = a.primary.Retrieve(resource)
	if err == nil {
		return res, cached, err
	}
	m, err = parse.NewParser().Parse(resource, a.secondary)
	res, err = protojson.Marshal(m)
	if err != nil {
		return res, false, err
	}
	return res, cached, nil
}
