@echo off
setlocal

set MATTER=%1

set "SCRIPT_DIR=%~dp0"
set "ROOT_DIR=%SCRIPT_DIR%.."

if not exist "%ROOT_DIR%\go\bin\legal-bot.exe" (
    echo Building legal-bot CLI...
    pushd "%ROOT_DIR%\go"
    go build -o bin\legal-bot.exe .\cmd\legal-bot
    popd
)

if "%MATTER%"=="" (
    "%ROOT_DIR%\go\bin\legal-bot.exe" guide
) else (
    "%ROOT_DIR%\go\bin\legal-bot.exe" guide "%MATTER%"
)

pause
