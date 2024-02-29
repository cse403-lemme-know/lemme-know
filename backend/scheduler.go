package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/scheduler"
)

type Scheduler interface {
	Schedule(date time.Time, data any) error
}

type EventBridgeScheduler struct {
	client *scheduler.Scheduler
}

func NewEventBridgeScheduler(sess *session.Session) *EventBridgeScheduler {
	return &EventBridgeScheduler{
		client: scheduler.New(sess),
	}
}

func (eventBridgeScheduler *EventBridgeScheduler) Schedule(date time.Time, data any) error {
	name := aws.String(fmt.Sprintf("lemmeknow-event-%d", GenerateID()))
	schedule := fmt.Sprintf("at(%s)", date.Format("2006-01-02T15:04:05"))
	_, err := eventBridgeScheduler.client.CreateSchedule(&scheduler.CreateScheduleInput{
		Description:           nil,
		ActionAfterCompletion: aws.String(scheduler.ActionAfterCompletionDelete),
		FlexibleTimeWindow: &scheduler.FlexibleTimeWindow{
			Mode: aws.String(scheduler.FlexibleTimeWindowModeOff),
		},
		Name:               name,
		ScheduleExpression: aws.String(schedule),
		State:              aws.String(scheduler.ScheduleStateEnabled),
		GroupName:          aws.String("lemmeknow-backend"),
		Target: &scheduler.Target{
			Arn:     aws.String(os.Getenv("AWS_LAMBDA_ARN")),
			RoleArn: aws.String(os.Getenv("AWS_SCHEDULER_ROLE_ARN")),
			Input:   aws.String(string(mustMarshal(data))),
		},
	})
	return err
}

type LocalScheduler struct {
}

func NewLocalScheduler() *LocalScheduler {
	return &LocalScheduler{}
}

func (localScheduler *LocalScheduler) Schedule(date time.Time, data any) error {
	dataJson := mustMarshal(data)
	go func() {
		time.Sleep(time.Until(date))
		Cron(dataJson)
	}()
	return nil
}
