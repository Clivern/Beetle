## Running Beetle inside Kubernetes Cluster

```bash
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.41.2/deploy/static/provider/cloud/deploy.yaml
$ kubectl get pods -n ingress-nginx \
  -l app.kubernetes.io/name=ingress-nginx --watch
$ kubectl get svc --namespace=ingress-nginx
```

Update `beetle.yaml` with the database credentials. The following part inside the file

```bash
....

# Application Database
database:
    # Database driver (sqlite3, mysql)
    driver: ${BEETLE_DATABASE_DRIVER:-mysql}
    # Hostname
    host: ${BEETLE_DATABASE_MYSQL_HOST:-REPLACE_WITH_MYSQL_HOSTNAME}
    # Port
    port: ${BEETLE_DATABASE_MYSQL_PORT:-3306}
    # Database
    name: ${BEETLE_DATABASE_MYSQL_DATABASE:-REPLACE_WITH_MYSQL_DATABASE}
    # Username
    username: ${BEETLE_DATABASE_MYSQL_USERNAME:-REPLACE_WITH_MYSQL_USERNAME}
    # Password
    password: ${BEETLE_DATABASE_MYSQL_PASSWORD:-REPLACE_WITH_MYSQL_PASSWORD}
....
```

Deploy a sample application and Beetle API server

```bash
$ kubectl apply -f sample_app.yaml --record
$ kubectl apply -f beetle.yaml --record

$ kubectl get ingress

# Update /etc/hosts with the ingress IP
# 167.x.x.x    example.com

$ kubectl describe ingress toad-ing
$ kubectl describe ingress beetle-ing

$ curl http://example.com/toad/_ready
$ curl http://example.com/beetle/_ready
```

Interact with Beetle API server

```bash
# Get clusters
$ curl http://example.com/beetle/api/v1/cluster -H "X-API-KEY: 1234" -s | jq .

# Get cluster
$ curl http://example.com/beetle/api/v1/cluster/production -H "X-API-KEY: 1234" -s | jq .

# Get cluster namespaces
$ curl http://example.com/beetle/api/v1/cluster/production/namespace -H "X-API-KEY: 1234" -s | jq .

# Get namespace
$ curl http://example.com/beetle/api/v1/cluster/production/namespace/default -H "X-API-KEY: 1234" -s | jq .

# Get namespace applications
$ curl http://example.com/beetle/api/v1/cluster/production/namespace/default/app -H "X-API-KEY: 1234" -s | jq .

# Get application `toad`
$ curl http://example.com/beetle/api/v1/cluster/production/namespace/default/app/toad -H "X-API-KEY: 1234" -s | jq .

# Get Async Jobs
$ curl -X GET http://example.com/beetle/api/v1/job -H "X-API-KEY: 1234" -s | jq .

# Deploy a new version with recreate strategy
$ curl -X POST \
     -H "X-API-KEY: 1234" \
     -d '{"version":"0.2.4","strategy":"recreate"}' \
     http://example.com/beetle/api/v1/cluster/production/namespace/default/app/toad/deployment

# Get application `toad` version
$ curl http://example.com/beetle/api/v1/cluster/production/namespace/default/app/toad -H "X-API-KEY: 1234" -s | jq .

# Another deployment with ramped strategy
$ curl -X POST \
     -H "X-API-KEY: 1234" \
     -d '{"version":"0.2.3","strategy":"ramped", "maxSurge": "1", "maxUnavailable": "0"}' \
     http://example.com/beetle/api/v1/cluster/production/namespace/default/app/toad/deployment
```
