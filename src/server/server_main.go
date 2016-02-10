package main

import (
	"encoding/json"
	"fmt"
	"github.com/drone/routes"
	"io/ioutil"
	"net/http"
	//	"log"
	"github.com/bitly/go-simplejson"
	// "github.com/davecgh/go-spew/spew"
	"math/rand"
	"strconv"
)

type DeviceSwitch struct {
	Id     string
	Name   string
	Room   string
	Status string
}

func getLight(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	light_id := params.Get(":light")
	desired_state := params.Get(":state")
	fmt.Fprint(w, "requesting switch ", light_id, "\n")
	switches := LoadDevices()
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
	switches := LoadDevices()
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

func LoadDevices() map[string]DeviceSwitch {
	random := strconv.FormatFloat(rand.Float64(), 'f', 2, 32)
	url := "http://10.0.1.61/port_3480/data_request?id=user_data&rand=" + string(random)
	resp, err := http.Get(url)
	devices := make(map[string]DeviceSwitch)
	if err != nil {
		return devices
	}

	defer resp.Body.Close()
	raw_body, err := ioutil.ReadAll(resp.Body)
	body := []byte(string(raw_body))

	parsed_body, err := simplejson.NewJson(body)
	// fmt.Println(parsed_body.Get("devices"))
	parsed_keys := parsed_body.Get("devices").MustArray()
	for _, v := range parsed_keys {
		var deviceSwitch DeviceSwitch
		for property, value := range v.(map[string]interface{}) {
			// fmt.Println("\t", property)
			switch value_key := value.(type) {
			case string:
				if property == "name" {
					deviceSwitch.Name = value_key
				} else if property == "room" {
					deviceSwitch.Room = value_key
				}
				// 	fmt.Println("\t\t", value_key)
				break
			case json.Number:
				if property == "id" {
					deviceSwitch.Id = string(value_key)
				}
				break
			case map[string]interface{}:
				/*
					for sub_property, sub_value := range value_key {
						fmt.Println("\t\t\t", sub_property)
						for sub_valuekey, sub_valuevalue := range sub_value.(map[string]interface{}) {
							fmt.Println("\t\t\t\t", sub_valuekey, "\t", sub_valuevalue)
						}
					}
				*/
				break
			case []interface{}:
				/*
					for _, sub_value := range value_key {
						switch sub_value_type := sub_value.(type) {
						default:
							fmt.Printf("%T", sub_value_type)
						}
					}
				*/
			default:
				/*
					fmt.Printf("%T", value)
					fmt.Println("\t\t", value, " - ")
				*/
			}
		}
		if deviceSwitch.Id == "3" {
			deviceSwitch.Status = parsed_body.Get("devices").GetIndex(2).Get("states").GetIndex(19).Get("value").MustString("-1")

		} else {
			deviceSwitch.Status = "-1"
		}
		devices[deviceSwitch.Id] = deviceSwitch
	}
	return devices
}

func main() {
	mux := routes.New()
	mux.Get("/lights", printLights)
	mux.Get("/lights/:light/:state", getLight)
	mux.Get("/", printTest)
	http.Handle("/", mux)
	http.ListenAndServe(":8088", nil)
}
