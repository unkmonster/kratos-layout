# Kratos Project Template

mono repo kratos template

## Init

```bash
make mono
git init
```

## Create a service

```bash
make service NAME=<SERVICE_NAME>
```

## Build service (Docker)

```bash
docker build -f ./build/Dockerfile --build-arg SERVICE_NAME=<SERVICE_NAME> --build-arg CMD_NAME=<CMD_NAME> .
```

## Next step

1. 设置 Makefile

## 约定

### Service Name

service_name 是服务目录相对于 app 的目录名

### Module Format

```go
moduleName := strings.ReplaceAll(serviceName, "/", "-")
```

### Git Tag

```
{type}/{module}/{version}
```
