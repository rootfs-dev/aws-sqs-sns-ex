package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	region      = flag.String("region", "us-east-1", "aws region")
	credFile    = flag.String("credFile", "", " path to aws credentials file")
	credProfile = flag.String("credProfile", "default", "credential profile")
	topic       = flag.String("topic", "foobar", "SNS topic name")
)

func main() {
	flag.Parse()
	sess := session.New(&aws.Config{
		Region:      aws.String(*region),
		Credentials: credentials.NewSharedCredentials(*credFile, *credProfile),
	})

	svc := sns.New(sess)
	createTopicParam := &sns.CreateTopicInput{
		Name: topic,
	}
	createTopicRes, err := svc.CreateTopic(createTopicParam)
	if err != nil {
		glog.Fatalf("failed to create topic %q: %v", *topic, err)
	}

	params := &sns.PublishInput{
		Message:  aws.String("ping"),
		TopicArn: aws.String(*createTopicRes.TopicArn),
	}

	resp, err := svc.Publish(params)

	if err != nil {
		glog.Fatalf("failed to publish message %v", err)
		return
	}

	fmt.Println(resp)
}
