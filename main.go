package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	responseString string
	title          string
)

func main() {
	var port int
	var rootContext string
	flag.IntVar(&port, "port", 8080, "Port the webserver should launch with")
	flag.StringVar(&rootContext, "rootContext", "", "Root context for the webserver")
	flag.StringVar(&responseString, "response", "Default hello from go code", "Content of respone element")
	flag.StringVar(&title, "title", "Default title", "Content of title element")
	flag.Parse()
	
	log.Printf("Staring application on port %d...\n", port)
	if !strings.HasSuffix(rootContext, "/") {
		rootContext = rootContext + "/"
	}
	if !strings.HasPrefix(rootContext, "/") {
		rootContext = "/" + rootContext
	}

	http.HandleFunc(fmt.Sprintf("%s", rootContext), helloHandler)
	http.HandleFunc(fmt.Sprintf("%serrors", rootContext), errorHandler)

	server := http.Server{Addr: fmt.Sprintf(":%d", port)}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown received, exiting...")

	server.Shutdown(context.Background())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := WebResponse{
		Message:  responseString,
		Title:    title,
		Hostname: os.Getenv("HOSTNAME"),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write(generateHTML(response))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	serverErrorRate := readIntParam("500", r)
	rNumber := rand.Intn(101)
	if rNumber < serverErrorRate {
		http.Error(w, fmt.Sprintf("500 error returned as %d < %d", rNumber, serverErrorRate), http.StatusInternalServerError)
	} else {
		w.Write([]byte(fmt.Sprintf("OK as %d >= %d", rNumber, serverErrorRate)))
	}
}

func readIntParam(name string, r *http.Request) (value int) {
	serverErrorRate := r.URL.Query().Get(name)
	if i, err := strconv.Atoi(serverErrorRate); err == nil {
		value = i
	}
	return
}

func generateHTML(input WebResponse) []byte {
	templateDir := os.Getenv("KO_DATA_PATH")
	if templateDir == "" {
		templateDir = "kodata"
	}
	tmpl, err := template.New("index.html").ParseFiles(path.Join(templateDir, "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, input)
	return []byte(buf.Bytes())
}

//WebResponse struct returned from webendpoint
type WebResponse struct {
	Title    string `json:"info"`
	Message  string `json:"response"`
	Hostname string `json:"hostname"`
}
