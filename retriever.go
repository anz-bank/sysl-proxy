package sysl_proxy

import (
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/joshcarp/gop/gop"
	"google.golang.org/protobuf/encoding/protojson"
)

type Retriever struct {
	primary      gop.Retriever
	secondary    gop.Retriever
}

func New(primary, secondary gop.Retriever) Retriever {
	return Retriever{primary: primary, secondary: secondary}
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
	p := parse.NewParser()
	p.SetVersioned()
	m, err = p.Parse(resource, a.secondary)
	res, err = protojson.Marshal(m)
	if err != nil {
		return res, false, err
	}
	return res, cached, nil
}
