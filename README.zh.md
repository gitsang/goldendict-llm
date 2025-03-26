# goldendict-llm

使用 Golang 编写的 GoldenDict 外部程序，用于对接 LLM API 进行翻译。

![](./resources/preview.png)

## 使用方式

### 1. 下载可执行文件

从 [Release](https://github.com/gitsang/goldendict-llm/releases/latest) 中下载 `goldendict-llm.exe`

### 2. 编写配置文件

```yaml
adapter: openai

adapters:
  openai:
    url: "https://api.openai.com/v1/chat/completions"
    token: "your-api-key-here"
    model: "gpt-3.5-turbo"
```

### 3. 在词典中配置程序

命令行应该类似如下：`D:\Application\GoldenDict-ng\content\goldendict-llm\goldendict-llm.exe -c "D:\Application\GoldenDict-ng\content\goldendict-llm\configs\config.local.yaml" "%GDWORD%"`

![](./resources/add-external-program.png)
