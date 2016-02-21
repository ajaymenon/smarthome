// Package vera_wrapper contains utility functions for querying VERA upnp api
package vera_wrapper

import (
	"encoding/json"
	// 	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/franela/goreq"
	"log"
	"strconv"
)

type DataRequestParams struct {
	Id           string
	OutputFormat string
}

type Room struct {
	Id   int
	Name string
}

type Device struct {
	Name         string
	Id           int
	Room         int
	RoomName     string
	Status       string
	Watts        string
	Category     int
	CategoryName string
}

type Category struct {
	Name string
	Id   int
}

func TurnOnDevice(device *Device) {
	SetDeviceState(device, "1")
}

func TurnOffDevice(device *Device) {
	SetDeviceState(device, "0")
}

func SetDeviceState(device *Device, newValue string) {
	_, err := goreq.Request{
		Uri: "http://10.0.1.61:3480/data_request?id=action&output_format=json&DeviceNum=" + strconv.Itoa(device.Id) + "&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&newTargetValue=" + newValue,
	}.Do()
	if err != nil {
		log.Fatal(err)
	}
	device.Status = newValue
}

func LoadDevices() (devices map[int]*Device, categories map[int]*Category, rooms map[int]*Room) {
	res, _ := goreq.Request{
		Uri: "http://10.0.1.61:3480/data_request?id=sdata&output_format=json",
	}.Do()
	json_response, _ := jason.NewObjectFromReader(res.Response.Body)
	categories_json, _ := json_response.GetObjectArray("categories")
	devices_json, _ := json_response.GetObjectArray("devices")
	rooms_json, _ := json_response.GetObjectArray("rooms")

	devices = make(map[int]*Device)
	categories = make(map[int]*Category)
	rooms = make(map[int]*Room)
	for _, room_json := range rooms_json {
		room := new(Room)
		room_bytes, _ := room_json.Marshal()
		json.Unmarshal(room_bytes, &room)
		rooms[room.Id] = room
	}
	for _, category_json := range categories_json {
		category := new(Category)
		category_bytes, _ := category_json.Marshal()
		json.Unmarshal(category_bytes, &category)
		categories[category.Id] = category
	}
	for _, device_json := range devices_json {
		device := new(Device)
		device_bytes, _ := device_json.Marshal()
		json.Unmarshal(device_bytes, &device)
		device.RoomName = rooms[device.Room].Name
		device.CategoryName = categories[device.Category].Name
		devices[device.Id] = device
	}

	return devices, categories, rooms
}
