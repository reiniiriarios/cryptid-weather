package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

const FETCH_INTERVAL = 10

func main() {
	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Get config
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %s", err)
	}

	// Logs
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Create client
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://cryptid:public@172.16.0.131:1883")
	opts.SetClientID("cryptidWeather")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(1 * time.Second)

	// Connect
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Updates
	go func() {
		for range time.Tick(time.Second * FETCH_INTERVAL) {
			weatherUpdate(&mqttClient)
		}
	}()

	// Run server until interrupted
	<-done

	// Cleanup
	println("\nClosing...")
	mqttClient.Disconnect(250)
	println("Cryptid Weather Closed")
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func weatherUpdate(mqttClient *mqtt.Client) {
	weather, err := getCurrentWeather()
	if err != nil {
		println("Error fetching weather.", err.Error())
		return
	}

	temp_c := fmt.Sprintf("%.4f", weather.TempC)
	feelslike := fmt.Sprintf("%.4f", weather.FeelsLikeC)
	humidity := fmt.Sprint(weather.Humidity)

	wait := time.Second * 10
	_ = (*mqttClient).Publish("weather/temperature", 0, false, temp_c).WaitTimeout(wait)
	_ = (*mqttClient).Publish("weather/feelslike", 0, false, feelslike).WaitTimeout(wait)
	_ = (*mqttClient).Publish("weather/humidity", 0, false, humidity).WaitTimeout(wait)
	_ = (*mqttClient).Publish("weather/condition", 0, false, weather.Code).WaitTimeout(wait)
}
