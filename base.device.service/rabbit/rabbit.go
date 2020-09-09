package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//RabbitConnection connection to rabbitmq server
var RabbitConnection *Rabbit

//Rabbit connection
type Rabbit struct {
	Connection  *amqp.Connection
	Channel     *amqp.Channel
	Credentials *Credentials
}

//Publish message to given exchange with routing key
func (r *Rabbit) Publish(message string, exchange string, rk string, headers amqp.Table, retry int) {
	if retry == 0 {
		return
	}
	err := r.Channel.Publish(
		exchange, // exchange
		rk,       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
			DeliveryMode:    2,
			Headers:         headers, //make(amqp.Table, 0),
			Body:            []byte(message),
		})
	log.Println("[Rabbit]Publish message to ", exchange, ":", rk, "; message:", message)
	if err != nil {
		log.Println("[Rabbit]Error in publishing message to ", exchange, ":", rk, "; message:", message, "; Error:", err.Error())
		r.ReInitialize()
		retry = retry - 1
		r.Publish(message, exchange, rk, headers, retry)
	}
}

//ReInitialize connection
func (r *Rabbit) ReInitialize() error {
	r.Channel.Close()
	r.Connection.Close()
	newConnection, err := createRabbitConnection(r.Credentials)
	if err == nil {
		r.Connection = newConnection.Connection
		r.Channel = newConnection.Channel
	}
	return err
}

//Credentials represents credentials for rabbit connection
type Credentials struct {
	Host     string
	Port     int
	Username string
	Password string
}

//InitializeRabbitConnection establish new rabbitmq connection
func InitializeRabbitConnection(credentials *Credentials) {
	RabbitConnection, _ = createRabbitConnection(credentials)
	/*
		config.Config.Rabbit.Host,
			config.Config.Rabbit.Port,
			config.Config.Rabbit.Username,
			config.Config.Rabbit.Password,
		)*/
}

//CreateRabbitConnection establish new rabbitmq connection
func createRabbitConnection(cr *Credentials) (*Rabbit, error) {
	connectionStr := fmt.Sprintf("amqp://%v:%v@%v:%v/", cr.Username, cr.Password, cr.Host, cr.Port)
	conn, err := amqp.Dial(connectionStr)
	failOnError(err, "Failed to connect to RabbitMQ; Connection string:"+connectionStr)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	if err != nil {
		return nil, err
	}

	return &Rabbit{
		Connection:  conn,
		Channel:     ch,
		Credentials: cr,
	}, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
