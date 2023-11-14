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

const FETCH_INTERVAL = 60

func main() {
	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	plog("Starting...")

	// Get config
	plog("Reading config...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %s", err)
	}

	// Logs
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Create client
	plog("Creating client...")
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
			plog("Updating weather...")
			weatherUpdate(&mqttClient)
		}
	}()

	// Run server until interrupted
	<-done

	// Cleanup
	println()
	plog("Closing...")
	mqttClient.Disconnect(250)
	plog("Cryptid Weather Closed")
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func weatherUpdate(mqttClient *mqtt.Client) {
	weather, err := getCurrentWeather()
	if err != nil {
		plog("Error fetching weather.")
		plog(err.Error())
		return
	}
	plog(weather)

	// Stringify all data!
	temp_c := fmt.Sprintf("%.4f", weather.TempC)
	feelslike := fmt.Sprintf("%.4f", weather.FeelsLikeC)
	humidity := fmt.Sprint(weather.Humidity)
	code := fmt.Sprint(weather.Code)
	is_day := "1"
	if !weather.IsDay {
		is_day = "0"
	}

	// Publish individually.
	publish(mqttClient, "weather/temperature", temp_c)
	publish(mqttClient, "weather/feelslike", feelslike)
	publish(mqttClient, "weather/humidity", humidity)
	publish(mqttClient, "weather/condition", weather.Condition)
	publish(mqttClient, "weather/code", code)
	publish(mqttClient, "weather/isday", is_day)
}

func publish(c *mqtt.Client, topic string, payload string) error {
	t := (*c).Publish(topic, 0, false, payload)
	if !t.WaitTimeout(time.Second * 10) {
		plog(t.Error().Error())
		return t.Error()
	}
	return nil
}
