package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	region      = flag.String("region", "us-east-1", "aws region")
	credFile    = flag.String("credFile", "", " path to aws credentials file")
	credProfile = flag.String("credProfile", "default", "credential profile")
	queue       = flag.String("queue", "channel-awssqs-default-aloha", "sqs queue name")
	msg         = flag.String("message", "hello world", "message body")
)

func main() {
	flag.Parse()
	sess := session.New(&aws.Config{
		Region:      aws.String(*region),
		Credentials: credentials.NewSharedCredentials(*credFile, *credProfile),
	})

	svc := sqs.New(sess)
	// create queue
	createParam := &sqs.CreateQueueInput{
		QueueName: queue,
	}

	createQueueResp, err := svc.CreateQueue(createParam)
	if err != nil {
		fmt.Println("failed to create queue %v", err)
		return
	}
	fmt.Println("queue url", *createQueueResp.QueueUrl)
	// list queue
	param := &sqs.ListQueuesInput{
		QueueNamePrefix: queue,
	}
	resp, err := svc.ListQueues(param)
	if err != nil {
		fmt.Println("failed to list queues", err)
		return
	}
	for i := range resp.QueueUrls {
		fmt.Println("queue: ", *resp.QueueUrls[i])
	}

	// send message
	sendParam := &sqs.SendMessageInput{
		MessageBody: aws.String(*msg),
		QueueUrl:    aws.String(*createQueueResp.QueueUrl),
	}
	sendResp, err := svc.SendMessage(sendParam)
	if err != nil {
		fmt.Println("failed to send", err)
		return
	}
	fmt.Println("send resp", sendResp)

	// receive message
	receiveParam := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(*createQueueResp.QueueUrl),
	}
	receiveResp, err := svc.ReceiveMessage(receiveParam)
	if err != nil {
		fmt.Println("failed to receive", err)
		return
	}
	fmt.Println("recv", receiveResp)

	// Delete message
	for _, message := range receiveResp.Messages {
		deleteParams := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(*createQueueResp.QueueUrl),
			ReceiptHandle: message.ReceiptHandle,
		}
		_, err := svc.DeleteMessage(deleteParams) // No response returned when successed.
		if err != nil {
			fmt.Println("failed to delete msg", err)
		}
		fmt.Println("Message has beed deleted:", *message.MessageId)
	}
}
