package broker

// func (ox *RabbitMQSvc) handlerTransactionNotification(msg *amqp.Delivery) (errx serror.SError) {
// 	defer func() {
// 		if errx != nil {
// 			logger.Err(errx)
// 		}
// 	}()

// 	var dataMsg models.CreateTransactionRequest

// 	err := json.Unmarshal(msg.Body, &dataMsg)
// 	if err != nil {
// 		errx = serror.NewFromErrorc(err, "[handler][handlerTransactionNotification] failed while Unmarshal rabbitMQ dataMsg")
// 		return
// 	}

// 	_, errx = ox.transactionUsecase.TransactionNotificationProcess(context.Background(), dataMsg)
// 	if errx != nil {
// 		errx.AddCommentf("[handler][handlerTransactionNotification][RequestID:%d] while transactionUsecase.TransactionNotification", dataMsg.RequestID)
// 		logger.Err(errx)
// 		return
// 	}

// 	return
// }
