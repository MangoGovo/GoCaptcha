# GoServerBase

一个基于gin的Web服务开发模板

### Lint(代码格式检查)

#### 手动格式化+lint检测

```shell
gofmt -w .
gci write . -s standard -s default
golangci-lint run --config .golangci.yml
```

#### 集成到IDE中

[配置方法](https://golangci-lint.run/welcome/integrations/)