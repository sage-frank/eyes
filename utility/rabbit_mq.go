package utility

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

//
//// Exchange 交换机
//type Exchange struct {
//	Name                string // exchange名称
//	Type                string // exchange类型，支持direct、topic、fanout、headers、x-delayed-message
//	RoutingKey          string // 路由key
//	XDelayedMessageType string // 延时消息类型，支持direct、topic、fanout、headers
//}
//
//// Producer 生产者对象
//type Producer struct {
//	queueName string
//	exchange  *Exchange
//	conn      *amqp.Connection
//	ch        *amqp.Channel
//}
//
//type RMQ struct {
//	client *amqp.Connection
//	logger *log.Logger
//	err    error
//}
//
//func NewRMQ(client *amqp.Connection, logger *log.Logger) *RMQ {
//
//	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//	return &RMQ{
//		client: client,
//		logger: logger,
//	}
//}
//
//func (r *RMQ) Error() string {
//	return r.err.Error()
//}
//
//func (r *RMQ) Publish(name string, msg []byte) {
//
//	ch, err := r.client.Channel()
//	if err != nil {
//		r.err = err
//		return
//	}
//
//	// 定义一个队列用来接受数据
//	q, err := ch.QueueDeclare(
//		name,  // name
//		false, // durable
//		false, // delete when unused
//		false, // exclusive
//		false, // no-wait
//		nil,   // arguments
//	)
//	if err != nil {
//		r.err = errors.Join(r.err, err)
//		return
//	}
//
//	// 发送订单消息，并设置expiration为30分钟
//	err = ch.Publish(
//		"",     // exchange
//		q.Name, // routing key
//		false,  // mandatory
//		false,  // immediate
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        msg,
//			Expiration:  "60000", // 30 minutes in milliseconds
//		},
//	)
//
//	if err != nil {
//		r.err = errors.Join(r.err, err)
//		return
//	}
//}
//
//func (r *RMQ) Consume(name string) {
//	ch, err := r.client.Channel()
//	if err != nil {
//		r.err = err
//		return
//	}
//
//	// 声明订单队列
//	_, err = ch.QueueDeclare(
//		name,  // name
//		false, // durable
//		false, // delete when unused
//		false, // exclusive
//		false, // no-wait
//		nil,   // arguments
//	)
//
//	if err != nil {
//		r.err = err
//		return
//	}
//
//	// 声明死信队列
//	_, err = ch.QueueDeclare(
//		"orders_dlq", // name
//		false,        // durable
//		false,        // delete when unused
//		false,        // exclusive
//		false,        // no-wait
//		nil,          // arguments
//	)
//
//	if err != nil {
//		r.err = err
//	}
//
//	// 绑定死信队列到交换机（如果需要的话）
//	// 这里假设你已经设置了DLX和相应的路由键
//	messages, err := ch.Consume(
//		"orders", // queue
//		"",       // consumer
//		true,     // auto-ack
//		false,    // exclusive
//		false,    // no-local
//		false,    // no-wait
//		nil,      // args
//	)
//	if err != nil {
//		r.err = err
//	}
//
//	forever := make(chan bool)
//
//	go func() {
//		for d := range messages {
//			r.logger.Printf("Received a message: %s", d.Body)
//			// 模拟处理订单和付款的流程
//			if processOrder(d.Body) {
//				r.logger.Printf("Order paid: %s", d.Body)
//			} else {
//				// 如果订单未付款，我们什么也不做，消息将因为TTL到期而死信
//				r.logger.Printf("Order not paid within 30 minutes: %s", d.Body)
//			}
//		}
//	}()
//
//	r.logger.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//	<-forever
//}
//
//// 模拟处理订单的函数
//func processOrder(order []byte) bool {
//	// 假设这里有一个付款的逻辑，如果在30分钟内返回true，则订单付款成功
//	// 这里仅用于演示，所以立即返回false
//	time.Sleep(2 * time.Minute) // 注意：这行代码仅用于演示，实际上应该移除或调整
//	return false
//}

var (
	uri           = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange      = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType  = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue         = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey    = flag.String("key", "test-key", "AMQP binding key")
	consumerTag   = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime      = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
	verbose       = flag.Bool("verbose", true, "enable verbose output of message data")
	autoAck       = flag.Bool("auto_ack", false, "enable message auto-ack")
	ErrLog        = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	Log           = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
	deliveryCount = 0
)

func init() {
	// flag.Parse()
}

func test() {
	c, err := NewConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	SetupCloseHandler(c)

	if *lifetime > 0 {
		Log.Printf("running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		Log.Printf("running until Consumer is done")
		<-c.done
	}

	Log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		ErrLog.Fatalf("error during shutdown: %s", err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func SetupCloseHandler(consumer *Consumer) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Log.Printf("Ctrl+C pressed in Terminal")
		if err := consumer.Shutdown(); err != nil {
			ErrLog.Fatalf("error during shutdown: %s", err)
		}
		os.Exit(0)
	}()
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, tag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     tag,
		done:    make(chan error),
	}

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	Log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	go func() {
		Log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	Log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}

	Log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %w", err)
	}

	Log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %w", err)
	}

	Log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %w", err)
	}

	Log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		*autoAck,   // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %w", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	defer Log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()

	for d := range deliveries {
		deliveryCount++
		if *verbose == true {
			Log.Printf(
				"got %dB delivery: [%v] %q",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
			)
		} else {
			if deliveryCount%65536 == 0 {
				Log.Printf("delivery count %d", deliveryCount)
			}
		}
		if *autoAck == false {
			err := d.Ack(false)
			if err != nil {
			}
		}
	}
}
