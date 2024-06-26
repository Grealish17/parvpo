package server

import (
	"github.com/Grealish17/parvpo/infrastructure/logger"
	"github.com/Grealish17/parvpo/internal/model"
)

func Distribute(msgChan <-chan model.Message) {
	for msg := range msgChan {
		respChans[msg.ID] <- msg
		close(respChans[msg.ID])
		delete(respChans, msg.ID)
	}
	//log.Println("Api channel closed, exiting distribute goroutine")
	logger.Info("Api channel closed, exiting distribute goroutine")
}
