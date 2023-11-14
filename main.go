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

const FETCH_INTERVAL = 20

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
			println("Updating weather...")
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
	fmt.Printf("%v\n", weather)

	// Stringify all data!
	temp_c := fmt.Sprintf("%.4f", weather.TempC)
	feelslike := fmt.Sprintf("%.4f", weather.FeelsLikeC)
	humidity := fmt.Sprint(weather.Humidity)
	code := fmt.Sprint(weather.Code)

	wait := time.Second * 10
	var t mqtt.Token
	t = (*mqttClient).Publish("weather/temperature", 0, false, temp_c)
	if !t.WaitTimeout(wait) {
		println(t.Error().Error())
	}
	t = (*mqttClient).Publish("weather/feelslike", 0, false, feelslike)
	if !t.WaitTimeout(wait) {
		println(t.Error().Error())
	}
	t = (*mqttClient).Publish("weather/humidity", 0, false, humidity)
	if !t.WaitTimeout(wait) {
		println(t.Error().Error())
	}
	t = (*mqttClient).Publish("weather/condition", 0, false, weather.Condition)
	if !t.WaitTimeout(wait) {
		println(t.Error().Error())
	}
	t = (*mqttClient).Publish("weather/code", 0, false, code)
	if !t.WaitTimeout(wait) {
		println(t.Error().Error())
	}
}
