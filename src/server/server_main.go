package main

import (
	"encoding/json"
	"fmt"
	"github.com/drone/routes"
	// "io/ioutil"
	"log"
	"net/http"
	// "github.com/bitly/go-simplejson"
	// "github.com/davecgh/go-spew/spew"
	"github.com/boltdb/bolt"
	// "math/rand"
	// "strconv"
	"vera_wrapper"
)

type DeviceSwitch struct {
	Id     string
	Name   string
	Room   string
	Status string
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	/*
		params := r.URL.Query()
		light_id := params.Get(":light")
		desired_state := params.Get(":state")
		fmt.Fprint(w, "requesting switch ", light_id, "\n")
		device := switches[light_id]
		json.NewEncoder(w).Encode(device)
		fmt.Println(device)
		url := "http://10.0.1.61/port_3480/data_request?id=lu_action&output_format=json&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&rand=0.123&DeviceNum=3&newTargetValue=" + desired_state
		resp, err := http.Get(url)
		if err != nil {
			return
		}

		defer resp.Body.Close()
	*/

}

func getDevices(w http.ResponseWriter, r *http.Request) {
	// switches := LoadDevices()
	// fmt.Println(switches)

	var switches = make(map[int]DeviceSwitch)
	json.NewEncoder(w).Encode(switches)
}

func loadDevices() map[int]*vera_wrapper.Device {
	db, err := bolt.Open("vera_status.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	devices := make(map[int]*vera_wrapper.Device)
	if err := db.View(func(tx *bolt.Tx) error {
		device_bucket := tx.Bucket([]byte("devices"))
		if err2 := device_bucket.ForEach(func(k, v []byte) error {
			device := new(vera_wrapper.Device)
			json.Unmarshal(v, &device)
			devices[device.Id] = device
			return nil
		}); err2 != nil {
			return err2
		}
		return nil
	}); err != nil {
		log.Fatal(err)
		return devices
	}

	return devices
}

func initDB() {
	db, err := bolt.Open("vera_status.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	devices, categories, rooms := vera_wrapper.LoadDevices()
	for _, device := range devices {
		fmt.Printf("%+v\n", *device)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("devices"))
			if err != nil {
				return err
			}
			encoded, err := json.Marshal(device)
			if err != nil {
				return err
			}
			return b.Put([]byte(device.Name), []byte(encoded))
		})
	}
	for _, category := range categories {
		fmt.Printf("%+v\n", *category)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("categories"))
			if err != nil {
				return err
			}
			encoded, err := json.Marshal(category)
			if err != nil {
				return err
			}
			return b.Put([]byte(category.Name), []byte(encoded))
		})

	}
	for _, room := range rooms {
		fmt.Printf("%+v\n", *room)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("rooms"))
			if err != nil {
				return err
			}
			encoded, err := json.Marshal(room)
			if err != nil {
				return err
			}
			return b.Put([]byte(room.Name), []byte(encoded))
		})

	}
	defer db.Close()
}

func main() {
	initDB()
	mux := routes.New()
	mux.Get("/api/devices", getDevices)
	http.Handle("/", mux)
	http.ListenAndServe(":8088", nil)
}
