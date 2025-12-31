#!/bin/bash
# Wait for LocalStack to be ready
while ! curl -s http://localhost:4566/_localstack/health | grep -q '"sns": "available"'; do
  echo "Waiting for LocalStack SNS to be ready..."
  sleep 2
done

# Create SNS topic
awslocal sns create-topic --name MyTopic

# Create SQS queue
awslocal sqs create-queue --queue-name MyQueue

# Subscribe queue to topic (hardcoded ARNs for LocalStack)
TOPIC_ARN="arn:aws:sns:us-east-1:000000000000:MyTopic"
QUEUE_ARN="arn:aws:sqs:us-east-1:000000000000:MyQueue"

awslocal sns subscribe --topic-arn "$TOPIC_ARN" --protocol sqs --notification-endpoint "$QUEUE_ARN"

echo "SNS topic and SQS queue created and subscribed"