# Kratos Project Template

mono repo kratos template

## Init

```bash
make mono
git init
```

## Create a service

```bash
make service name=<SERVICE_NAME>
```

## Build service (Docker)

```bash
docker build -f ./build/Dockerfile --build-arg SERVICE_NAME=<SERVICE_NAME> --build-arg CMD_NAME=<CMD_NAME> .
```

## Next step

## 约定

### Module Format

```go
moduleName := strings.ReplaceAll(serviceName, "/", "-")
```

### Git Tag

```
{type}/{module}/{version}
```
