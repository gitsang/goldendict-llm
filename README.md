# goldendict-llm

GoldenDict external program written in Golang for translating with LLM API.

![preview.png](./resources/preview.png)

## 1. Usage

### 1.1 Download Executable

Download `goldendict-llm.exe` from [Release](https://github.com/gitsang/goldendict-llm/releases/latest)

### 1.2 Write Configuration File

You can define different Adapters to quickly switch by modifying the `adapter` field.

```yaml
adapter: siliconflow

adapters:
  openai:
    url: "https://api.openai.com/v1/chat/completions"
    token: "your-api-key-here"
    model: "gpt-3.5-turbo"
  siliconflow:
    url: "https://api.siliconflow.cn/v1/chat/completions"
    token: "your-api-key-here"
    model: "Qwen/Qwen2.5-7B-Instruct"
```

### 1.3 Configure Program in Dictionary

The command line should look like: `C:\GoldenDict\goldendict-llm.exe -c "C:\GoldenDict\config.yaml" "%GDWORD%"`

![add-external-program.png](./resources/add-external-program.png)
