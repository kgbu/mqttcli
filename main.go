package main

import (
	"bufio"
	"os"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var usage = `
Usage here
`

func initFunc() {
	log.SetLevel(log.WarnLevel)
}

// connects MQTT broker
func connect(c *cli.Context, opts *MQTT.ClientOptions) (*MQTTClient, error) {
	log.Info("Connecting...")

	client := &MQTTClient{Opts: opts}
	_, err := client.Connect()
	if err != nil {
		return nil, err
	}
	log.Info("Connected")

	return client, nil
}

func publish(c *cli.Context) {
	if c.Bool("d") {
		log.SetLevel(log.DebugLevel)
	}

	opts := NewOption(c)

	willPayload := c.String("will-payload")
	willQoS := c.Int("will-qos")
	willRetain := c.Bool("will-retain")
	willTopic := c.String("will-topic")
	if willPayload != "" && willTopic != "" {
		opts.SetWill(willTopic, willPayload, MQTT.QoS(willQoS), willRetain)
	}

	client, err := connect(c, opts)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	qos := c.Int("q")
	topic := c.String("t")
	if topic == "" {
		log.Errorf("Please specify topic")
		os.Exit(1)
	}
	log.Infof("Topic: %s", topic)

	retain := c.Bool("r")

	if c.Bool("s") {
		// Read from Stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			err = client.Publish(topic, []byte(scanner.Text()), qos, retain)
			if err != nil {
				log.Error(err)
			}

		}
	} else {
		payload := c.String("m")
		err = client.Publish(topic, []byte(payload), qos, retain)
		if err != nil {
			log.Error(err)
		}

	}
	log.Info("Published")
}
func subscribe(c *cli.Context) {
	if c.Bool("d") {
		log.SetLevel(log.DebugLevel)
	}
	opts := NewOption(c)
	client, err := connect(c, opts)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	qos := c.Int("q")
	topic := c.String("t")
	if topic == "" {
		log.Errorf("Please specify topic")
		os.Exit(1)
	}
	log.Infof("Topic: %s", topic)

	err = client.Subscribe(topic, qos)
	if err != nil {
		log.Error(err)
	}
}

func main() {
	initFunc()

	app := cli.NewApp()
	app.Name = "mqttcli"
	app.Usage = usage

	commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			Usage:  "mqtt host to connect to. Defaults to localhost",
			EnvVar: "MQTT_HOST"},
		cli.IntFlag{
			Name:   "p, port",
			Value:  1883,
			Usage:  "network port to connect to. Defaults to 1883",
			EnvVar: "MQTT_PORT"},
		cli.StringFlag{
			Name:   "u,user",
			Value:  "",
			Usage:  "provide a username",
			EnvVar: "MQTT_USERNAME"},
		cli.StringFlag{
			Name:   "P,password",
			Value:  "",
			Usage:  "provide a password",
			EnvVar: "MQTT_PASSWORD"},
		cli.StringFlag{"t", "", "mqtt topic to publish to.", ""},
		cli.IntFlag{"q", 0, "QoS", ""},
		cli.StringFlag{"cafile", "", "CA file", ""},
		cli.StringFlag{"i", "", "ClientiId. Defaults random.", ""},
		cli.StringFlag{"m", "test message", "Message body", ""},
		cli.BoolFlag{"r", "message should be retained.", ""},
		cli.BoolFlag{"d", "enable debug messages", ""},
		cli.BoolFlag{"insecure", "do not check that the server certificate", ""},
		//		cli.StringFlag{"conf", "~/.mqtt.cfg", "config file", ""},
	}
	pubFlags := append(commonFlags,
		cli.BoolFlag{"s", "read message from stdin, sending line by line as a message", ""},
		cli.StringFlag{
			Name:  "will-payload",
			Value: "",
			Usage: "payload for the client Will",
		},
		cli.IntFlag{
			Name:  "will-qos",
			Value: 0,
			Usage: "QoS level for the client Will",
		},
		cli.BoolFlag{
			Name:  "will-retain",
			Usage: "if given, make the client Will retained",
		},
		cli.StringFlag{
			Name:  "will-topic",
			Value: "",
			Usage: "the topic on which to publish the client Will",
		},
	)
	subFlags := commonFlags

	app.Commands = []cli.Command{
		{
			Name:   "pub",
			Usage:  "publish",
			Flags:  pubFlags,
			Action: publish,
		},
		{
			Name:   "sub",
			Usage:  "subscribe",
			Flags:  subFlags,
			Action: subscribe,
		},
	}
	app.Run(os.Args)
}
