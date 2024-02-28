package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type Scheduler interface {
	Schedule(date time.Time, data any) error
}

type EventBridgeScheduler struct {
	client *eventbridge.EventBridge
}

func NewEventBridgeScheduler(sess *session.Session) *EventBridgeScheduler {
	return &EventBridgeScheduler{
		client: eventbridge.New(sess),
	}
}

func (eventBridgeScheduler *EventBridgeScheduler) Schedule(date time.Time, data any) error {
	name := aws.String(fmt.Sprintf("lemmeknow-event-%d", GenerateID()))
	schedule := fmt.Sprintf("at(%s)", date.Format("2006-01-02T15:04:05"))
	request, response := eventBridgeScheduler.client.PutRuleRequest(&eventbridge.PutRuleInput{
		Description:        nil,
		EventPattern:       nil,
		Name:               name,
		ScheduleExpression: aws.String(schedule),
		State:              aws.String(eventbridge.RuleStateEnabled),
	})
	if err := request.Send(); err != nil {
		return err
	}
	if response.RuleArn == nil {
		return fmt.Errorf("no rule ARN")
	}
	request2, _ := eventBridgeScheduler.client.PutTargetsRequest(&eventbridge.PutTargetsInput{
		Rule: name,
		Targets: []*eventbridge.Target{
			{
				Id:  aws.String("backend"),
				Arn: aws.String(os.Getenv("AWS_LAMBDA_ARN")),
				//RoleArn: aws.String(os.Getenv("LAMBDA_EVENTBRIDGE_ROLE_ARN")),
				Input: aws.String(string(mustMarshal(data))),
			},
		},
	})
	if err := request2.Send(); err != nil {
		return err
	}
	return nil
}
