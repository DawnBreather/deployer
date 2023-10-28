FROM alpine

#ENV GITLAB_CI_COMMIT_SHORT_SHA "17b51953"
#ENV GITLAB_CI_PROJECT_NAME "breathesmart"
#ENV GITLAB_CI_COMMIT_BRANCH "release-3.6-stallergenes"
#ENV ARTIFACT_OBJECT_REFERENCE_PREFIX "aptar-digital-health/web"
#ENV AWS_GIT_URL_WITH_CREDENTIALS_SECRETS_MANAGER_ARN "arn:aws:secretsmanager:us-east-1:010987917155:secret:git-8X4tHy"

COPY deployer.deployments_status_tracker /usr/bin/deployer.deployments_status_tracker
RUN chmod +x /usr/bin/deployer.deployments_status_tracker

CMD /usr/bin/deployer.deployments_status_tracker