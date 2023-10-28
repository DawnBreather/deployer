FROM alpine

RUN apk add git

#ENV GITLAB_CI_COMMIT_SHORT_SHA "17b51953"
#ENV GITLAB_CI_PROJECT_NAME "breathesmart"
#ENV GITLAB_CI_COMMIT_BRANCH "release-3.6-stallergenes"
#ENV ARTIFACT_OBJECT_REFERENCE_PREFIX "aptar-digital-health/web"
#ENV AWS_GIT_URL_WITH_CREDENTIALS_SECRETS_MANAGER_ARN "arn:aws:secretsmanager:us-east-1:010987917155:secret:git-8X4tHy"

COPY deployer.auto_deployments_invoker /usr/bin/deployer.auto_deployments_invoker
RUN chmod +x /usr/bin/deployer.auto_deployments_invoker

RUN git config --global user.email "devops@abcloudz.com"
RUN git config --global user.name "ABCloudz DevOps"

CMD /usr/bin/deployer.auto_deployments_invoker