package main

import (
	"embed"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/google/gousb"
	"github.com/s-urbaniak/uevent"
)

const (
	tvid = "a4e"
	tpid = "1000"
)

type Firmware []struct {
	RType uint8  `json:"bmRequestType"`
	Req   uint8  `json:"bRequest"`
	Val   uint16 `json:"wValue"`
	Idx   uint16 `json:"wIndex"`
	Data  []byte `json:"data_fragment"`
}

type Event struct {
	*uevent.Uevent
}

func (e *Event) Product() (vid string, pid string) {
	if val, ok := e.Vars["PRODUCT"]; ok {
		parts := strings.Split(val, "/")
		vid = parts[0]
		pid = parts[1]
	}
	return
}

func (e *Event) ToUint16(str string) uint16 {
	val, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		log.Fatal("here", err)
	}
	return uint16(val)
}

func (e *Event) HandleEvent(fw Firmware) {
	vid, pid := e.Product()
	if e.Action == "add" && e.Vars["DEVTYPE"] == "usb_device" && vid == tvid && pid == tpid {

		ctx := gousb.NewContext()
		defer ctx.Close()

		log.Println("Connecting to device", tvid, "Model", pid)
		dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(e.ToUint16(vid)), gousb.ID(e.ToUint16(pid)))
		if err != nil {
			log.Fatalf("Could not open a device: %v", err)
		}
		defer dev.Close()

		log.Println("Attempting firmware update")
		for i := 0; i < len(fw); i++ {
			p := fw[i]
			_, err := dev.Control(p.RType, p.Req, p.Val, p.Idx, p.Data)
			if err != nil {
				log.Fatalf("Could not open a device: %v", err)
			}
		}

		log.Println("Firmware update completed")

	}
}

//go:embed fw.json
var fs embed.FS

func main() {
	//load firmware
	file, err := fs.ReadFile("fw.json")
	if err != nil {
		log.Fatal(err)
	}

	fw := Firmware{}
	err = json.Unmarshal([]byte(file), &fw)
	if err != nil {
		log.Fatal(err)
	}

	r, err := uevent.NewReader()
	if err != nil {
		log.Fatal(err)
	}

	dec := uevent.NewDecoder(r)
	for {
		e, err := dec.Decode()
		if err != nil {
			log.Fatal(err)
		}

		evt := &Event{e}
		evt.HandleEvent(fw)
	}
}
