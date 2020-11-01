<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://user-images.githubusercontent.com/32605850/97817997-df110f80-1cf3-11eb-9fae-2db765d09563.png" alt="Project logo"></a>
</p>


<h3 align="center">Sysl Proxy</h3>

<div align="center">

  [![Status](https://img.shields.io/badge/status-active-success.svg)]() 
  [![GitHub Issues](https://img.shields.io/github/issues/joshcarp/sysl-proxy)](https://github.com/joshcarp/sysl-proxy/issues)
  [![GitHub Pull Requests](https://img.shields.io/github/issues-pr/joshcarp/sysl-proxy)](https://github.com/joshcarp/sysl-proxy/pulls)
  [![License](https://img.shields.io/badge/license-apache2-blue.svg)](/LICENSE)

</div>

---


## 📝 Table of Contents
- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [Authors](#authors)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>
Sysl Proxy is a proxy service to supply .sysl and .pb.json specifications to sysl clients.
It is modelled off go modules, and depends on [gop](https://github.com/joshcarp/gop).
It is designed to be deployed on google cloud functions, using the ServeHTTP function.

## 🏁 Getting Started <a name = "getting_started"></a>
There are a couple of environment variables that need to be set:
- `CacheLocation`: either the bucket or directory that the .sysl contents will be cached to.
- `CacheLocationSyslJson`: either the bucket or directory that the .sysl.pb.json contents will be cached to.
- `FsType`: os/mem/gcs for os filesystem, in memory filesystem, or a gcs bucket respectively
- `PROJECT_NUM`: The unique project number for the google cloud project

Secrets stored in Google secret manager:
- `GH_TOKEN`: A github [personal access token](https://github.com/settings/tokens)
- `HTTP_PROXY`: HTTP proxy url to use

See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites
What things you need to install the software and how to install them.
- Go 1.13: currently google cloud functions only support upto the go 1.13 runtime

## 🔧 Running the tests <a name = "tests"></a>

`go test ./...`

## 🎈 Usage <a name="usage"></a>
`go run ./servehttp/`
- This will run a sysl proxy server on `localhost:8082`


## 🚀 Deployment <a name = "deployment"></a>

- See .github/workflows/cloud-function-deploy.yml

## ⛏️ Built Using <a name = "built_using"></a>
- [Google Cloud Functions](https://cloud.google.com/functions/) - Deployment
- [Google Cloud Storage](https://cloud.google.com/storage/) - Asset caching
- [Google Secret Manager](https://cloud.google.com/secret-manager) - Secret storage
- [Golang](https://golang.org/) - Server 

## ✍️ Authors <a name = "authors"></a>
- [@joshcarp](https://github.com/joshcarp)

## 🎉 Acknowledgements <a name = "acknowledgement"></a>
- Go Modules: Athens Project: https://github.com/gomods/athens 