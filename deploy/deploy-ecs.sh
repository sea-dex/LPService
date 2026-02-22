#!/usr/bin/env bash

SRV_NAME=$1
ECS_SERVICE_NAME=$2
ECS_TASK_DEFINITION_NAME=$3

echo "Image ${DOCKER_REGISTRY}/${APP_FULL_NAME}:${DOCKER_TAG_COMMIT}"

echo "Deploy ${SRV_NAME} service"

TASK_DEFINITION=$(aws ecs describe-task-definition \
  --task-definition ${ECS_TASK_DEFINITION_NAME} --region ${AWS_REGION})

TASK_ROLE=$(echo ${TASK_DEFINITION} | jq -r '.taskDefinition.taskRoleArn')
EXECUTE_ROLE=$(echo ${TASK_DEFINITION} | jq -r '.taskDefinition.executionRoleArn')
CPU=$(echo ${TASK_DEFINITION} | jq -r '.taskDefinition.cpu')
MEMORY=$(echo ${TASK_DEFINITION} | jq -r '.taskDefinition.memory')

NEW_CONTAINER_DEFINITION=$(echo ${TASK_DEFINITION} | \
  jq --arg IMAGE ${DOCKER_REGISTRY}/${APP_FULL_NAME}:${DOCKER_TAG_COMMIT} \
  '.taskDefinition.containerDefinitions[0].image = $IMAGE | .taskDefinition.containerDefinitions[0]') \

echo "Registering new ${SRV_NAME} container definition..."

aws ecs register-task-definition --region ${AWS_REGION} --family ${ECS_TASK_DEFINITION_NAME} \
  --network-mode awsvpc --requires-compatibilities FARGATE \
  --task-role-arn ${TASK_ROLE} --execution-role-arn ${EXECUTE_ROLE} \
  --cpu ${CPU} --memory ${MEMORY} \
  --container-definitions "${NEW_CONTAINER_DEFINITION}"

NEW_TASK_DEFINITION=$(aws ecs describe-task-definition \
  --task-definition ${ECS_TASK_DEFINITION_NAME} --region ${AWS_REGION})

ARN=$(echo ${NEW_TASK_DEFINITION} | jq -r '.taskDefinition.taskDefinitionArn')

echo "Tagging new ${SRV_NAME} task definition $ARN ..."

aws ecs tag-resource --region ${AWS_REGION} --resource-arn "${ARN}" \
  --tags key=Application,value=${APP_NAME} key=Service,value=${SRV_NAME} key=Stage,value=${ENV_NAME} \
  key=Repository,value=${CI_PROJECT_NAME} key=Branch,value=${CI_COMMIT_BRANCH} \
  key=GitlabCI,value=true

echo "Updating the ${SRV_NAME} service..."

aws ecs update-service --region ${AWS_REGION} \
  --cluster ${ECS_CLUSTER_NAME} --service ${ECS_SERVICE_NAME} \
  --task-definition ${ECS_TASK_DEFINITION_NAME}
