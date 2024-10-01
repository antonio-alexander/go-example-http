package internal

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
)

const (
	CERT_FILE  string = "CERT_FILE"
	KEY_FILE   string = "KEY_FILE"
	HTTP_PORT  string = "HTTP_PORT"
	HTTPS_PORT string = "HTTPS_PORT"
)

var (
	httpPort          string = "80"
	httpsPort         string = "443"
	certFile          string = "./certs/ssl.crt"
	keyFile           string = "./certs/ssl.key"
	templateIndexHtml string = "./templates/index.tmpl"
)

type payloadIndexHtml struct {
	Subhead string `json:"subhead"`
}

func Main(pwd string, args []string, envs map[string]string, osSignal chan os.Signal) error {
	var tlsCertificate tls.Certificate
	var wg sync.WaitGroup
	var err error

	//print version
	fmt.Printf("Version: \"%s\"\n", Version)
	fmt.Printf("Git Commit: \"%s\"\n", GitCommit)
	fmt.Printf("Git Branch: \"%s\"\n", GitBranch)

	//load certs
	certFile, keyFile = envs[CERT_FILE], envs[KEY_FILE]
	httpPort, httpsPort = envs[HTTP_PORT], envs[HTTPS_PORT]
	tlsCertificate, err = GetCertificates()
	if err != nil {
		return err
	}

	// generate template
	bytes, err := os.ReadFile(templateIndexHtml)
	if err != nil {
		return err
	}
	tmplIndex, err := template.New("index.html").Parse(string(bytes))
	if err != nil {
		return err
	}

	//generate http handlers, create and start server
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		default:
			http.FileServer(http.Dir("./static")).ServeHTTP(writer, request)
		case "/":
			if err := tmplIndex.Execute(writer, &payloadIndexHtml{
				Subhead: "...for the world's best tasting cat food, sourced with only " +
					"the finest ingredients. Made with love from the heart of Mississippi.",
			}); err != nil {
				fmt.Printf("error while executing template: %s", err)
			}
		}
	})
	httpsServer := &http.Server{
		Addr: ":" + httpsPort,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCertificate},
		},
	}
	httpServer := &http.Server{
		Addr: ":" + httpPort,
	}
	fmt.Printf("...starting http server on :%s\n", httpPort)
	fmt.Printf("...starting https server on :%s\n", httpsPort)
	stopped := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()

		if err := httpsServer.ListenAndServeTLS(certFile, keyFile); err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()

		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()
	select {
	case <-stopped:
	case <-osSignal:
		if err := httpsServer.Close(); err != nil {
			fmt.Println(err)
		}
		if err := httpServer.Close(); err != nil {
			fmt.Println(err)
		}
		return nil
	}
	wg.Wait()
	return nil
}
