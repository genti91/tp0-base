package common

import (
	"bufio"
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	
	go func() {
		<-signalChan
		close(signalChan)
		log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
		if c.conn != nil {
			c.conn.Close()
		}
		os.Exit(0)
	}()


	c.createClientSocket()
	
	betProt := NewBetProt(c.config.ID)
	writer := bufio.NewWriter(c.conn)

	if err := betProt.SendBet(writer); err != nil {
		log.Errorf("action: apuesta_enviada | result: fail | dni: %v | error: %v",
			betProt.document,
			err,
		)
		return
	}

	c.conn.Close()

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		betProt.document,
		betProt.number,
	)

}
