package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/adapter"
	"github.com/muka/go-bluetooth/bluez/device"
)

type BTScanner struct {
	adapter *adapter.Adapter1
	devices map[string]*device.Device1
}

func NewScanner() *BTScanner {
	conn, err := api.NewConn()
	if err != nil {
		log.Fatal("Erro BlueZ:", err)
	}

	adapterID := "hci0"
	adp, err := adapter.NewAdapter1(conn, adapterID)
	if err != nil {
		log.Fatal("Erro adapter:", err)
	}

	return &BTScanner{
		adapter: adp,
		devices: make(map[string]*device.Device1),
	}
}

func (s *BTScanner) ScanClassic(ctx context.Context, duration int) {
	fmt.Printf("\n🔍 Classic Scan (%ds)...\n", duration)
	
	properties := map[string]interface{}{
		"DiscoveryFilter": map[string]interface{}{
			"Transport": "bredr",
		},
	}
	
	err := s.adapter.StartDiscovery(properties)
	if err != nil {
		fmt.Println("Erro start discovery:", err)
		return
	}
	defer s.adapter.StopDiscovery()

	time.Sleep(time.Duration(duration) * time.Second)
}

func (s *BTScanner) ListDevices() {
	fmt.Println("\n📱 Dispositivos Encontrados:")
	for addr, dev := range s.devices {
		name, _ := dev.GetName()
		fmt.Printf("  %s - %s\n", addr, name)
	}
}

func main() {
	scanner := NewScanner()
	
	fmt.Println("🔥 Go Bluetooth Pentest Scanner")
	fmt.Println("================================")
	
	// Ctrl+C handler
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n[+] Parando...")
		cancel()
		os.Exit(0)
	}()

	scanner.ScanClassic(ctx, 15)
	scanner.ListDevices()
}
