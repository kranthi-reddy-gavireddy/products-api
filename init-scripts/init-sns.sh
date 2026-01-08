#!/bin/bash
# Wait for LocalStack to be ready
while ! curl -s http://localhost:4566/_localstack/health | grep -q '"sns": "available"'; do
  echo "Waiting for LocalStack SNS to be ready..."
  sleep 2
done

# Create SNS topic
## awslocal sns create-topic --name OrderCreatedTopic

# Create SQS queue
awslocal sqs create-queue --queue-name OrderCreatedTopic

# Subscribe queue to topic (hardcoded ARNs for LocalStack)
## TOPIC_ARN="arn:aws:sns:us-east-1:000000000000:MyTopic"
QUEUE_ARN="arn:aws:sqs:us-east-1:000000000000:OrderCreatedTopic"

awslocal sqs set-queue-attributes --queue-url "$QUEUE_ARN" --attributes "Policy={\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"sns.amazonaws.com\"},\"Action\":\"sqs:SendMessage\",\"Resource\":\"$QUEUE_ARN\"}]}"

#echo "SNS topic and SQS queue created and subscribed"

echo "SQS queue 'OrderCreatedTopic' created."