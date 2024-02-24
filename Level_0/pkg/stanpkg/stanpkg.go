package stanpkg

import (
	"log"
	"os"
	"strconv"
	"os/signal"
	"syscall"
	"database/sql"
	"level_0/pkg/database"
	"level_0/pkg/cache"
	"level_0/pkg/validation"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"time"
)

const (
	natsURL    = stan.DefaultNatsURL
	clusterID  = "test-cluster"          
	subject    = "stan-subject" 
	durableID  = "durable-subscriber"     
	queueGroup = "queue-group"
	clientID   = "stan-publisher-1"         
)

func SubscribeAndListen(cache_instance *cache.Cache, db *sql.DB) {
    natsConn, err := nats.Connect(natsURL)
    if err != nil {
        log.Printf("Error connecting to NATS: %v", err)
        return
    }
    defer natsConn.Close()

    stanConn, err := stan.Connect(clusterID, "stan-subscriber-1", stan.NatsConn(natsConn))
    if err != nil {
        log.Printf("Error connecting to STAN: %v", err)
        return
    }
    defer stanConn.Close()

    subscription, err := stanConn.QueueSubscribe(subject, queueGroup, func(msg *stan.Msg) {
        receivedData := string(msg.Data)
        log.Printf("Received message: %s", receivedData)
		if err := validation.Validate(receivedData); err != nil {
			log.Printf("Error data not valid: %v", err)

		} else {
        orderID := cache_instance.Count() + 1
        cache_instance.Set(strconv.Itoa(orderID), receivedData, 20*time.Minute)
        database.InsertToDB(db, receivedData)
		}
        msg.Ack()
    }, stan.DurableName(durableID), stan.SetManualAckMode())
    if err != nil {
        log.Printf("Error subscribing: %v", err)
        return
    }
    defer subscription.Unsubscribe()

    signalCh := make(chan os.Signal, 1)
    signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
    <-signalCh

    log.Println("Shutting down gracefully")
}

func PublishMessage(msg []byte) {
    natsConn, err := nats.Connect(natsURL)
    if err != nil {
        log.Printf("Error connecting to NATS: %v", err)
        return
    }
    defer natsConn.Close()

    stanConn, err := stan.Connect(clusterID, clientID, stan.NatsConn(natsConn))
    if err != nil {
        log.Printf("Error connecting to STAN: %v", err)
        return
    }
    defer stanConn.Close()

    messageData := []byte(msg)
    err = stanConn.Publish(subject, messageData)
    if err != nil {
        log.Printf("Error publishing message: %v", err)
        return
    }

    log.Printf("Message published: %s", messageData)
}
