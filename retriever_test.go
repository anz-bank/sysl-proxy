package sysl_proxy

import (
	"fmt"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/retriever/retriever_github"
	"github.com/spf13/afero"
	"testing"
)

func TestRetriever(t *testing.T){
	r := New(gop_filesystem.New(afero.NewMemMapFs(), "/"), retriever_github.New())
	a, b, c := r.Retrieve("github.com/joshcarp/sysl-1", "sysl.sysl", "main")
	fmt.Println(a, b, c)
}