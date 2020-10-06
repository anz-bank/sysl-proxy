package sysl_proxy

import (
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/joshcarp/gop/gop"
	"google.golang.org/protobuf/encoding/protojson"
)

type Processor struct{
	retr gop.Retriever
	parser *parse.Parser
}

func NewProcessor(retr gop.Retriever)Processor{
	return Processor{retr:retr, parser: parse.NewParser()}
}

func (a Processor) Retrieve(resource string) ([]byte, bool, error) {
	m, err := a.parser.Parse(resource, a.retr)
	if err != nil {
		return nil, false, err
	}
	res, err := protojson.Marshal(m)
	if err != nil {
		return res, false, err
	}
	return res, false, nil
}

