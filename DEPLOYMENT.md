# Deployment Guide

This guide covers deploying the School Management Microservices to different environments.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Compose Production](#docker-compose-production)
3. [Kubernetes Deployment](#kubernetes-deployment)
4. [Environment Variables](#environment-variables)
5. [Monitoring & Logging](#monitoring--logging)
6. [Security Considerations](#security-considerations)

## Local Development

### Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd skool-management

# Set up development environment
make setup

# Verify all services are running
make health
```

### Individual Service Development

```bash
# Run services individually for development
make dev-auth     # Auth service on :8081
make dev-school   # School service on :8082
make dev-student  # Student service on :8083
make dev-gateway  # API Gateway on :8080
```

## Docker Compose Production

### Prerequisites

- Docker 20.0+
- Docker Compose 2.0+
- At least 4GB RAM
- 10GB free disk space

### Production Setup

1. **Create environment file:**

```bash
cp .env.example .env
# Edit .env with production values
```

2. **Update production secrets:**

```bash
# Generate secure JWT secrets
openssl rand -base64 64 > jwt_secret.txt
openssl rand -base64 64 > jwt_refresh_secret.txt

# Update .env file with these secrets
```

3. **Create production docker-compose:**

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  # Extend base docker-compose.yml
  auth-service:
    restart: unless-stopped
    environment:
      - JWT_SECRET=${JWT_SECRET}
      - JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Add similar configuration for other services
```

4. **Deploy:**

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### SSL/TLS Setup

For production, use a reverse proxy like Nginx:

```nginx
# nginx.conf
upstream api_gateway {
    server localhost:8080;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://api_gateway;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster 1.20+
- kubectl configured
- Helm 3.0+ (optional)

### Namespace Setup

```bash
kubectl create namespace school-management
kubectl config set-context --current --namespace=school-management
```

### Secrets Management

```bash
# Create JWT secrets
kubectl create secret generic jwt-secrets \
  --from-literal=jwt-secret="$(openssl rand -base64 64)" \
  --from-literal=jwt-refresh-secret="$(openssl rand -base64 64)"

# Create database secrets
kubectl create secret generic db-secrets \
  --from-literal=mongo-password="$(openssl rand -base64 32)" \
  --from-literal=school-db-password="$(openssl rand -base64 32)" \
  --from-literal=student-db-password="$(openssl rand -base64 32)"
```

### Database Deployments

**MongoDB (Auth Service):**

```yaml
# mongodb-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
        - name: mongodb
          image: mongo:7
          env:
            - name: MONGO_INITDB_ROOT_USERNAME
              value: admin
            - name: MONGO_INITDB_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db-secrets
                  key: mongo-password
          ports:
            - containerPort: 27017
          volumeMounts:
            - name: mongodb-storage
              mountPath: /data/db
      volumes:
        - name: mongodb-storage
          persistentVolumeClaim:
            claimName: mongodb-pvc
```

### Service Deployments

**Auth Service:**

```yaml
# auth-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
        - name: auth-service
          image: school-management/auth-service:latest
          env:
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: jwt-secrets
                  key: jwt-secret
            - name: MONGODB_URI
              value: "mongodb://admin:$(MONGO_PASSWORD)@mongodb:27017/authdb?authSource=admin"
          ports:
            - containerPort: 8081
          livenessProbe:
            httpGet:
              path: /health
              port: 8081
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8081
            initialDelaySeconds: 5
            periodSeconds: 5
```

### Ingress Configuration

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: school-management-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
    - hosts:
        - your-domain.com
      secretName: school-management-tls
  rules:
    - host: your-domain.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-gateway-service
                port:
                  number: 8080
```

### Deploy to Kubernetes

```bash
# Apply all configurations
kubectl apply -f k8s/

# Check deployment status
kubectl get pods
kubectl get services
kubectl get ingress
```

## Environment Variables

### Required Variables

| Variable             | Description               | Default   | Required |
| -------------------- | ------------------------- | --------- | -------- |
| `JWT_SECRET`         | JWT signing secret        | -         | Yes      |
| `JWT_REFRESH_SECRET` | JWT refresh token secret  | -         | Yes      |
| `MONGODB_URI`        | MongoDB connection string | -         | Yes      |
| `DB_HOST`            | PostgreSQL host           | localhost | Yes      |
| `DB_USER`            | PostgreSQL user           | -         | Yes      |
| `DB_PASSWORD`        | PostgreSQL password       | -         | Yes      |
| `DB_NAME`            | PostgreSQL database name  | -         | Yes      |

### Optional Variables

| Variable    | Description      | Default |
| ----------- | ---------------- | ------- |
| `PORT`      | Service port     | 8080    |
| `LOG_LEVEL` | Logging level    | info    |
| `DEV_MODE`  | Development mode | false   |

## Monitoring & Logging

### Prometheus Metrics

Add Prometheus metrics to each service:

```go
// Add to main.go
import "github.com/prometheus/client_golang/prometheus/promhttp"

http.Handle("/metrics", promhttp.Handler())
```

### Grafana Dashboard

Create dashboards for:

- Request rate and latency
- Database connections
- Memory and CPU usage
- Error rates

### Centralized Logging

Use ELK stack or similar:

```yaml
# filebeat.yml
filebeat.inputs:
  - type: docker
    containers.ids:
      - "*"
    processors:
      - add_docker_metadata: ~

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

## Security Considerations

### Production Checklist

- [ ] Use strong, unique JWT secrets
- [ ] Enable SSL/TLS
- [ ] Use secure database passwords
- [ ] Implement rate limiting
- [ ] Enable CORS properly
- [ ] Use network policies in Kubernetes
- [ ] Regular security updates
- [ ] Database encryption at rest
- [ ] Secret management (Vault, etc.)
- [ ] Network segmentation

### Database Security

```bash
# PostgreSQL security
ALTER USER postgres PASSWORD 'strong_password';
CREATE USER app_user WITH PASSWORD 'strong_password';
GRANT CONNECT ON DATABASE school_db TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
```

### Network Security

```yaml
# Kubernetes Network Policy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: school-management-netpol
spec:
  podSelector:
    matchLabels:
      app: school-management
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: api-gateway
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: database
```

## Backup & Recovery

### Database Backups

```bash
# MongoDB backup
mongodump --uri="mongodb://admin:password@localhost:27017/authdb" --out=/backup/mongo

# PostgreSQL backup
pg_dump -h localhost -U user -d schooldb > backup_school.sql
pg_dump -h localhost -U user -d studentdb > backup_student.sql
```

### Automated Backups

```yaml
# Kubernetes CronJob for backups
apiVersion: batch/v1
kind: CronJob
metadata:
  name: database-backup
spec:
  schedule: "0 2 * * *" # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: backup
              image: backup-image:latest
              command: ["/backup-script.sh"]
              volumeMounts:
                - name: backup-storage
                  mountPath: /backups
          restartPolicy: OnFailure
```

## Scaling

### Horizontal Scaling

```bash
# Scale services
kubectl scale deployment auth-service --replicas=3
kubectl scale deployment school-service --replicas=3
kubectl scale deployment student-service --replicas=3
```

### Database Scaling

- MongoDB: Replica sets and sharding
- PostgreSQL: Read replicas and connection pooling

### Auto-scaling

```yaml
# HorizontalPodAutoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```
