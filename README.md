# goldendict-llm

A GoldenDict external program written in Golang, using LLM API for translation.

![](./resources/preview.png)

## Usage

### 1. Download Executable

Download `goldendict-llm.exe` from [Release](https://github.com/gitsang/goldendict-llm/releases/latest)

### 2. Write Configuration File

```yaml
adapter: openai

adapters:
  openai:
    url: "https://api.openai.com/v1/chat/completions"
    token: "your-api-key-here"
    model: "gpt-3.5-turbo"
```

### 3. Configure Program in Dictionary

The command line should look similar to: `D:\Application\GoldenDict-ng\content\goldendict-llm\goldendict-llm.exe -c "D:\Application\GoldenDict-ng\content\goldendict-llm\configs\config.local.yaml" "%GDWORD%"`

![](./resources/add-external-program.png)
