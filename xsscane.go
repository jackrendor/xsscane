package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var banner string = `
___    ___ ________   ________  ________  ________  ________   _______      
|\  \  /  /|\   ____\ |\   ____\|\   ____\|\   __  \|\   ___  \|\  ___ \     
\ \  \/  / | \  \___|_\ \  \___|\ \  \___|\ \  \|\  \ \  \\ \  \ \   __/|    
 \ \    / / \ \_____  \\ \_____  \ \  \    \ \   __  \ \  \\ \  \ \  \_|/__  
  /     \/   \|____|\  \\|____|\  \ \  \____\ \  \ \  \ \  \\ \  \ \  \_|\ \ 
 /  /\   \     ____\_\  \ ____\_\  \ \_______\ \__\ \__\ \__\\ \__\ \_______\
/__/ /\ __\   |\_________\\_________\|_______|\|__|\|__|\|__| \|__|\|_______|
|__|/ \|__|   \|_________\|_________|                          by @jackrendor`

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//http.CanonicalHeaderKey()
	fmt.Printf("  \u001b[30;102m[%s] Connected\u001b[0m\n", r.RemoteAddr)
	fmt.Printf("%s %s%s\n", r.Method, r.Host, r.RequestURI)
	for key, value := range r.Header {
		fmt.Printf("\u001b[93m%s\u001b[0m: \u001b[96m%s\u001b[0m\n", key, strings.Join(value, ", "))
	}
	body, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		fmt.Println("Error reading body:", bodyErr.Error())
		return
	}
	fmt.Println("\u001b[96m" + string(body) + "\u001b[0m\n\n-----------------")
}

func callBackFromScript()

func main() {
	var listeningVar string
	flag.StringVar(&listeningVar, "listen", "0.0.0.0:8000", "Listen address, default: 0.0.0.0:8000")
	flag.Parse()

	http.HandleFunc("/", rootHandler)
	fmt.Printf("\u001b[36m%s\u001b[0m\n", banner)
	fmt.Println("\nListening on:", listeningVar)
	fmt.Println("-----------------")
	http.ListenAndServe(":8000", nil)
}
