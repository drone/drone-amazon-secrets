# drone-amazon-extension

A secret extension that provides optional support for sourcing secrets from the AWS Secrets Manager. _Please note this project requires Drone server version 1.3 or higher._

## Installation

Create a shared secret:

```text
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```text
$ docker run -d \
  --publish=3000:3000 \
  --env=DEBUG=true \
  --env=SECRET_KEY=bea26a2221fd8090ea38720fc445eca6 \
  --env=AWS_ACCESS_KEY_ID=... \
  --env=AWS_SECRET_ACCESS_KEY=... \
  --restart=always \
  --name=amazon-secrets drone/amazon-secrets
```

Update your Drone runner configuration to include the plugin address and the shared secret.

```text
DRONE_SECRET_ENDPOINT=http://1.2.3.4:3000
DRONE_SECRET_SECRET=bea26a2221fd8090ea38720fc445eca6
```
