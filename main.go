package main

import (
	"golang.org/x/net/websocket"
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var homeHtml = string(home_html())
var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.New("homeHtml").Parse(homeHtml))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func MonitorMultiFiles(names []string, out chan []string,
	h hub) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range names {
		result, _ := ReadLastNLines(name, 10)
		PrintMultiLines(result)
		MonitorFile(name, out, watcher)
	}
	for {
		select {
		case lines := <-out:
			content := strings.Join(lines, "\n")
			fmt.Print(content)
			h.broadcast <- content
		}
	}
	watcher.Close()
}

var usage = func() {
	fmt.Fprintf(os.Stderr,
		"Usage: gowebtail [FILE]...\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) < 1 {
		usage()
	}
	out := make(chan []string)
	go h.run()
	go MonitorMultiFiles(flag.Args(), out, h)
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
