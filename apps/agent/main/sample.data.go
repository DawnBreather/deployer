package main

var SampleData = `{
	"environments": {
		"ZGVtby5hZGhlcmVpdC5jb2hlcm8taGVhbHRoLmNvbQ": {
			"name": "demo.adhereit.cohero-health.com",
			"tier": {
				"name": "backend",
				"credentials": {
					"aws": {
						"token": "AKIATAMBS6UH654REHFT",
						"secret": "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",
						"region": "us-east-1"
					}
				},
				"artifacts": [{
					"source": [{
						"from": {
							"s3": {
								"bucket": "builds-206967076111",
								"region": "us-east-1",
								"object": "build-307.zip"
							}
						},
						"to": "webapi/"
					}],
					"middleware": {
						"envsubst": {
							"target": [
								"webapi/"
							],
							"variables": {
								"VARIABLE1": "value1",
								"VARIABLE2": "value2",
								"VARIABLE3": "value3"
							}
						}

					},
					"deploy": [{
						"from": "/tmp/webapi/build-307/Chat",
						"to": "/tmp/app/Chat",
            "mode": "delete"
					}],
					"post_deploy_script": {
						"shell": "powershell",
						"environment_variables": {
							"VARIABLE1": "value1",
							"VARIABLE2": "value2"
						},
						"command": [
							"inetmgr restart"
						]
					}
				}],
				"health-checks": {
					"tcp": {
						"targets": [{
							"host": "127.0.0.1",
							"port": 3000
						}],
						"options": {
							"timeout": 30,
							"interval": 30,
							"unhealthy_threshold": 2,
							"healthy_threshold": 2
						}
					},
					"mssql": {
						"targets": [{
							"host": "127.0.0.1",
							"username": "",
							"password": "",
							"database": "",
							"port": 1433
						}],
						"options": {
							"timeout": 30,
							"interval": 30,
							"unhealthy_threshold": 2,
							"healthy_threshold": 2
						}
					},
					"redis": {
						"targets": [{
							"host": "127.0.0.1",
							"username": "",
							"password": "",
							"port": 6379,
							"ssl": false
						}],
						"options": {
							"timeout": 30,
							"interval": 30,
							"unhealthy_threshold": 2,
							"healthy_threshold": 2
						}
					},
					"http": {
						"targets": [{
							"url": "http://127.0.0.1:80",
							"status": "200-399"
						}],
						"options": {
							"timeout": 30,
							"interval": 30,
							"unhealthy_threshold": 2,
							"healthy_threshold": 2
						}
					}
				}
			},
      "status": {
        "inventory": {
          "i-02fb3a3c2803989c1": {
            "agent_heartbeat_at": "",
            "installation": {
              "state": "accomplished",
              "message": "base64_encoded_string_here"
            },
            "health_checks": {
              "tcp": [
                {
                  "updated_at": "",
                  "target": "",
                  "status": "healthy",
                  "message": ""
                }
              ],
              "mssql": [
                {
                  "updated_at": "",
                  "target": "",
                  "status": "healthy",
                  "message": ""
                }
              ],
              "redis": [
                {
                  "updated_at": "",
                  "target": "",
                  "status": "healthy",
                  "message": ""
                }
              ],
              "http": [
                {
                  "updated_at": "",
                  "target": "",
                  "status": "healthy",
                  "message": ""
                }
              ]
            }
          }
        }
      }
		}
	}
}`
