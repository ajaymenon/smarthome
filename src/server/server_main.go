package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/drone/routes"
	"io/ioutil"
	"net/http"
	"os"
	//	"log"
	"github.com/bitly/go-simplejson"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"strconv"
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
	desired_state := params.Get(":state")
	fmt.Fprint(w, "requesting switch ", light_id, "\n")
	var switches map[string]DeviceSwitch
	switches = make(map[string]DeviceSwitch)
	Load("test_switches", &switches)
	device := switches[light_id]
	json.NewEncoder(w).Encode(device)
	fmt.Println(device)
	url := "http://10.0.1.61/port_3480/data_request?id=lu_action&output_format=json&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&rand=0.123&DeviceNum=3&newTargetValue=" + desired_state
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

}

func printLights(w http.ResponseWriter, r *http.Request) {
	var switches map[string]DeviceSwitch
	switches = make(map[string]DeviceSwitch)
	Load("test_switches", &switches)
	fmt.Println(switches)

	json.NewEncoder(w).Encode(switches)
}

func printTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parse arguments, you have to call this by yourself
	random := strconv.FormatFloat(rand.Float64(), 'f', 2, 32)
	url := "http://10.0.1.61/port_3480/data_request?id=user_data&rand=" + string(random)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	raw_body, err := ioutil.ReadAll(resp.Body)
	body := []byte(string(raw_body))
	var parsed_body interface{}
	err = json.Unmarshal(body, &parsed_body)
	// pretty_body, err := json.MarshalIndent(parsed_body, "", "  ")
	// fmt.Println(parsed_body)
	// spew.Dump(parsed_body)
	w.Write(body)
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

func LoadDevices() {
	random := strconv.FormatFloat(rand.Float64(), 'f', 2, 32)
	url := "http://10.0.1.61/port_3480/data_request?id=user_data&rand=" + string(random)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	raw_body, err := ioutil.ReadAll(resp.Body)
	body := []byte(string(raw_body))

	devices := make(map[string]DeviceSwitch)
	fmt.Println(devices)
	parsed_body, err := simplejson.NewJson(body)
	// fmt.Println(parsed_body.Get("devices"))
	parsed_keys := parsed_body.Get("devices").MustArray()
	for k, v := range parsed_keys {
		var deviceSwitch DeviceSwitch
		deviceSwitch.Id = strconv.Itoa(k)
		for k1, _ := range v.(map[string]interface{}) {
			fmt.Println("\t", k1)
		}
		devices[strconv.Itoa(k)] = deviceSwitch
	}
	fmt.Println(devices)
	spew.Sdump(parsed_body.Get("devices"))
}

func main() {
	var switches map[string]DeviceSwitch
	switches = make(map[string]DeviceSwitch)
	switches["1"] = DeviceSwitch{"1", "test1", "room", 0}
	switches["2"] = DeviceSwitch{"2", "test2", "room2", 1}
	LoadDevices()
	Save("test_switches", switches)
	mux := routes.New()
	mux.Get("/lights", printLights)
	mux.Get("/lights/:light/:state", getLight)
	mux.Get("/", printTest)
	http.Handle("/", mux)
	http.ListenAndServe(":8088", nil)
}
