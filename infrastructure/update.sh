#!/bin/bash

# Set your variables
REGION="<REGION>"
CLUSTER="<CLUSTER>"
SERVICE="<SERVICE>"
TASK_DEFINITION="<TASK_DEFINITION>"
ECS_TASK_EXECUTION_ROLE_ARN="<ECS_TASK_EXECUTION_ROLE_ARN>"

# Describe the task definition and save to a JSON file
TASK_DEF_ARN=$(
aws ecs describe-task-definition \
  --region "$REGION" \
  --task-definition "$TASK_DEFINITION" \
  --query 'taskDefinition.taskDefinitionArn' \
  --output text
)

# Update the service by the same task definition
aws ecs update-service \
  --region "$REGION" \
  --cluster "$CLUSTER" \
  --service "$SERVICE" \
  --task-definition "$TASK_DEF_ARN" \
  --force-new-deployment