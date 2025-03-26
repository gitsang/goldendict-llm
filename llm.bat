@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

if "%~1"=="" (
    echo 错误：请提供查询内容
    echo 用法：llm.bat "您的问题"
    pause
    exit /b 1
)

set "user_query=%~1"

if not exist llm.conf (
    echo 错误：找不到 llm.conf 文件
    echo 请创建 llm.conf 文件并包含以下内容：
    echo url: https://example.com/llm/endpoint
    echo token: your-token-here
    pause
    exit /b 1
)

set "api_url="
set "api_token="
for /f "tokens=1,2 delims==" %%a in (llm.conf) do (
    if "%%a"=="url" set "api_url=%%b"
    if "%%a"=="token" set "api_token=%%b"
)

if "!api_url!"=="" (
    echo 错误：无法从 llm.conf 读取URL
    echo 当前文件内容:
    type  llm.conf 
    pause
    exit /b 1
)

if "!api_token!"=="" (
    echo 错误：无法从 llm.conf 读取Token
    pause
    exit /b 1
)

set "user_query=!user_query:"=\"!"
set "temp_file=%TEMP%\llm_response.json"

curl -s "!api_url!" ^
  -H "Authorization: Bearer !api_token!" ^
  -H "Content-Type: application/json" ^
  -d "{\"model\": \"deepseek-r1\", \"messages\": [{\"role\": \"user\", \"content\": \"!user_query!\"}], \"stream\": false}" > "%temp_file%"

powershell -Command "(Get-Content '%temp_file%' -Raw | ConvertFrom-Json).choices[0].message.content"

del "%temp_file%" 2>nul

echo.
pause
endlocal
