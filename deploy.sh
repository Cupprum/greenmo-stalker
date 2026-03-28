#!/usr/bin/env bash

set -e
set -u
set -o pipefail

echo "Building..."
sam build --use-container

echo "Deploying..."
sam deploy --no-confirm-changeset --no-fail-on-empty-changeset \
    --stack-name greenmo-stalker-stack \
    --region ${GREENMO_AWS_REGION} \
    --capabilities CAPABILITY_IAM \
    --resolve-s3 \
    --parameter-overrides \
        GreenmoOpenMapsApiToken=${GREENMO_OPEN_MAPS_API_TOKEN} \
        GreenmoDomainName=${GREENMO_DOMAIN_NAME} \
        GreenmoCertificateArn=${GREENMO_CERTIFICATE_ARN}

echo "Finished"