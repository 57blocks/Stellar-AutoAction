#!/bin/bash

# Set your variables
REGION="<REGION>"
CLUSTER="<CLUSTER>"
SERVICE="<SERVICE>"
TASK_DEFINITION="<TASK_DEFINITION>"
ECS_TASK_EXECUTION_ROLE_ARN="<ECS_TASK_EXECUTION_ROLE_ARN>"


# Describe the task definition and save to a JSON file
aws ecs describe-task-definition \
  --region "$REGION" \
  --task-definition "$TASK_DEFINITION" \
  --query 'taskDefinition.{containerDefinitions:containerDefinitions}' \
  --output json > task-definition.json

TASK_DEF_ARN=$(
aws ecs register-task-definition \
  --region "$REGION" \
  --family "$TASK_DEFINITION" \
  --network-mode awsvpc \
  --cpu 1024 \
  --memory 2048 \
  --requires-compatibilities FARGATE \
  --execution-role-arn "$ECS_TASK_EXECUTION_ROLE_ARN" \
  --cli-input-json file://task-definition.json \
  --query 'taskDefinition.taskDefinitionArn' \
  --output text
)

aws ecs update-service \
  --region "$REGION" \
  --cluster "$CLUSTER" \
  --service "$SERVICE" \
  --task-definition "$TASK_DEF_ARN" \
  --force-new-deployment

# Clean up the JSON file if needed
rm task-definitions.json