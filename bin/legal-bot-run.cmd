@echo off
REM legal-bot-run: Spawn different LLM subprocesses for Legal-Bot (Windows)
REM
REM Delegates to the Go legal-bot CLI binary. Auto-downloads if not present.
REM
REM Usage: legal-bot-run <llm> "<prompt>" [output_file] [persona]

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "VERN_CLI=%SCRIPT_DIR%..\go\bin\legal-bot.exe"

REM Auto-build if Go binary doesn't exist
if not exist "%VERN_CLI%" (
    where go >nul 2>nul
    if not errorlevel 1 (
        if exist "%SCRIPT_DIR%..\go\cmd\legal-bot\main.go" (
            echo [legal-bot-run] Building legal-bot CLI from source... 1>&2
            pushd "%SCRIPT_DIR%..\go"
            go build -o bin\legal-bot.exe .\cmd\legal-bot 2>nul
            popd
        )
    )
)

REM Auto-download if still not available
if not exist "%VERN_CLI%" (
    if exist "%SCRIPT_DIR%install-legal-bot-cli.cmd" (
        call "%SCRIPT_DIR%install-legal-bot-cli.cmd"
    )
)

if not exist "%VERN_CLI%" (
    echo [legal-bot-run] Error: legal-bot CLI not available. Run: cd go ^&^& go build -o bin\legal-bot.exe .\cmd\legal-bot 1>&2
    exit /b 1
)

REM Map positional args to Go CLI flags
set "LLM=%~1"
set "PROMPT=%~2"
set "OUTPUT_FILE=%~3"
set "PERSONA=%~4"

if "%LLM%"=="" (
    echo Usage: legal-bot-run ^<llm^> "^<prompt^>" [output_file] [persona]
    echo   llm: claude ^| codex ^| gemini
    exit /b 1
)
if "%PROMPT%"=="" (
    echo Usage: legal-bot-run ^<llm^> "^<prompt^>" [output_file] [persona]
    exit /b 1
)

set "ARGS=run %LLM% "%PROMPT%""
if not "%OUTPUT_FILE%"=="" set "ARGS=%ARGS% --output "%OUTPUT_FILE%""
if not "%PERSONA%"=="" set "ARGS=%ARGS% --persona %PERSONA%"

"%VERN_CLI%" %ARGS%
