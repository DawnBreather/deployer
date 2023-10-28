FROM amazon/aws-cli

WORKDIR /deployer

#COPY --from=dawnbreather/buildtools:latest /usr/bin/deployer_agent /usr/bin/deployer_agent
COPY deployer_agent /usr/bin/deployer_agent

ENTRYPOINT ["/usr/bin/deployer_agent"]
