package sysl_proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/joshcarp/gop"
	gop3 "github.com/joshcarp/gop/gop"
	"github.com/joshcarp/gop/gop/cli"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/gop_gcs"
	"github.com/joshcarp/gop/gop/modules"
	"github.com/joshcarp/gop/gop/retriever/retriever_github"
	"github.com/joshcarp/gop/gop/retriever/retriever_wrapper"
	"github.com/spf13/afero"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var syslpbjsonfs = afero.NewMemMapFs()
var jsonfs = afero.NewMemMapFs()

const pbjsonaccept = `application/sysl.pb.json`

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var b, res []byte
	var cached bool

	defer func() { gop.HandleErr(w, err) }()
	reqestedResource := r.URL.Query().Get("resource")

	/* Make sure we're actually requesting a resource that is allowed */
	switch _, resource, _, _ := gop3.ProcessRequest(reqestedResource); path.Ext(resource) {
	case ".sysl", ".json", ".yaml", ".yml", ".proto":
	default:
		err = gop3.BadRequestError
		return
	}

	var accept = r.Header.Get("Accept")
	/* Make a new Gopper */
	s, _ := NewGopper(
		os.Getenv("CacheLocation"),
		os.Getenv("CacheLocationSyslJson"),
		os.Getenv("FsType"),
		accept,
	)

	res, cached, err = s.Retrieve(reqestedResource)
	if err != nil || res == nil {
		return
	}
	if !cached {
		if err := s.Cache(reqestedResource, res); err != nil {
			return
		}
	}

	switch accept {
	case "text/plain":
		b = res
	default:
		b, err = json.Marshal(gop3.Object{Content: res, Resource: reqestedResource})
		if err != nil {
			return
		}
	}
	if _, err := w.Write(b); err != nil {
		log.Println(err)
	}

}

func MemoryFs(accept string) afero.Fs {
	switch accept {
	case pbjsonaccept:
		return syslpbjsonfs
	default:
		return jsonfs
	}
}

func MemoryLoc(accept string) string {
	switch accept {
	case pbjsonaccept:
		return "sysl_pb_json"
	default:
		return "sysl"
	}
}

/* NewGopper returns a GopperService for a config; This Gopper can use an os filesystem, memory filesystem or a gcs bucket*/
func NewGopper(cachelocation, cachelocationsysljson, fsType, accept string) (*gop.GopperService, error) {
	r := gop.GopperService{}
	switch fsType {
	case "os":
		r.Gopper = gop_filesystem.New(afero.NewOsFs(), MemoryLoc(accept))
	case "mem", "memory", "":
		r.Gopper = gop_filesystem.New(MemoryFs(accept), "/")
	case "gcs":
		if accept == pbjsonaccept {
			cachelocation = cachelocationsysljson
		}
		gcs := gop_gcs.New(cachelocation)
		r.Gopper = &gcs
	}
	gh := retriever_github.New(
		cli.TokensFromString(
			"github.com:"+Secret("GH_TOKEN")))
	proxyURL, err := url.Parse(Secret("HTTP_PROXY"))
	if err != nil {
		return nil, err
	}
	gh.Client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	switch accept {
	case pbjsonaccept:
		r.Retriever =
			retriever_wrapper.New(
				NewProcessor(
					modules.New(
						gh,
						"sysl_modules/sysl_modules.yaml")))
	default:
		r.Retriever =
			retriever_wrapper.New(
				modules.New(
					gh,
					"sysl_modules/sysl_modules.yaml"))

	}
	return &r, nil
}

func Secret(name string) string {
	fmt.Println("Accessing Secret")
	secretClinet, _ := secretmanager.NewClient(context.Background())
	s, err := secretClinet.AccessSecretVersion(context.Background(), &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", os.Getenv("PROJECT_NUM"), name),
	})
	if err != nil {
		fmt.Println("Error accessing secret")
		return ""
	}
	fmt.Println("Secret retrieved")
	return string(s.Payload.Data)
}
