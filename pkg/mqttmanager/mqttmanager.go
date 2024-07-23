package mqttmanager

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"infinilapse-unified/pkg/cloud"
	"infinilapse-unified/pkg/dslrMgmt"
	"infinilapse-unified/pkg/webcamMgmt"
	"os"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func StartMQTTClient(broker string, port int, topic string, captureFunc func()) mqtt.Client {
	opts := mqtt.NewClientOptions()
	brokerConnectString := fmt.Sprintf("tcp://%s:%d", broker, port)
	println(brokerConnectString)
	opts.AddBroker(brokerConnectString)
	opts.SetClientID("infinilapse_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message on topic: %s\n", msg.Topic())
		captureFunc()
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return nil
	}

	return client
}

func CaptureAllCameras() {
	fmt.Printf("Begin cap loop triggered by MQTT.  Setting the stage\n")

	var capturedFiles []string
	if os.Getenv("DSLR_CAPTURE") != "false" {
		capturedFiles = append(capturedFiles, dslrMgmt.CaptureAllDslr()...)
		fmt.Printf("dslrMgmt.CaptureAllDslr()...\n%v\n", capturedFiles)
	}
	if os.Getenv("WEBCAM_CAPTURE") != "false" {
		capturedFiles = append(capturedFiles, webcamMgmt.CaptureWebCams()...)
		fmt.Printf("webcamMgmt.CaptureWebCams()...\n%v\n", capturedFiles)
	}
	fmt.Printf("Finished cap loop triggered by MQTT.  Unsetting the stage\n")

	err := cloud.IndexGoogleCloudStorageAndGraphQL(capturedFiles)
	if err != nil {
		_ = fmt.Errorf("cloud.IndexGoogleCloudStorageAndGraphQL(filePaths) %s\n", err)
	}
}
