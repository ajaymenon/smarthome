package main

import (
	"encoding/gob"
	// "encoding/json"
	"fmt"
	"github.com/drone/routes"
	// "io/ioutil"
	"net/http"
	"os"
	//	"log"
	"strings"
)

type DeviceSwitch struct {
	Id     string
	Name   string
	Room   string
	Status int
}

func getLight(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	light_id := params.Get(":light")
	fmt.Fprint(w, "requesting switch ", light_id, "\n")
	var switches []DeviceSwitch
	Load("test_switches", &switches)
	for _, device := range switches {
		if device.Id == light_id {
			fmt.Fprint(w, device, "\n")
		}
	}
}

func printLights(w http.ResponseWriter, r *http.Request) {
	var switches []DeviceSwitch
	fmt.Fprintf(w, "Loading switches\n")
	Load("test_switches", &switches)

	for _, device := range switches {
		// b, _ := json.Marshal(device)
		fmt.Fprint(w, device, "\n")
	}
}

func printTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "HEllo world!")
}

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	} else {
		fmt.Println(err)
	}

	file.Close()
	return err
}

func main() {
	var switches []DeviceSwitch
	switches = append(switches, DeviceSwitch{"1", "test1", "room", 0})
	switches = append(switches, DeviceSwitch{"2", "test2", "room2", 1})
	Save("test_switches", switches)
	mux := routes.New()
	mux.Get("/lights", printLights)
	mux.Get("/lights/:light", getLight)
	mux.Get("/", printTest)
	http.Handle("/", mux)
	http.ListenAndServe(":8088", nil)
}
