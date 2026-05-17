@echo off
setlocal

REM Legal-Bot Intake Pipeline
REM Scans source documents and builds the knowledge layer

set MATTER=%1
if "%MATTER%"=="" set MATTER=TRO-MTE-GAL-SM

echo.
echo  ============================================
echo   Legal-Bot Intake Pipeline
echo   Matter: %MATTER%
echo  ============================================
echo.

REM Find the repo root (parent of scripts/)
set "SCRIPT_DIR=%~dp0"
set "ROOT_DIR=%SCRIPT_DIR%.."

REM Check if matter directory exists
if not exist "%ROOT_DIR%\matters\%MATTER%\input" (
    echo ERROR: Matter directory not found: matters\%MATTER%\input
    echo Create it first or check the matter name.
    exit /b 1
)

REM Check for input files
dir /b "%ROOT_DIR%\matters\%MATTER%\input\*.txt" "%ROOT_DIR%\matters\%MATTER%\input\*.md" >nul 2>&1
if errorlevel 1 (
    echo ERROR: No .txt or .md files found in matters\%MATTER%\input\
    echo Drop your source documents there first.
    exit /b 1
)

REM Build the Go binary if needed
if not exist "%ROOT_DIR%\go\bin\legal-bot.exe" (
    echo Building legal-bot CLI...
    pushd "%ROOT_DIR%\go"
    go build -o bin\legal-bot.exe .\cmd\legal-bot
    popd
)

REM Run intake
"%ROOT_DIR%\go\bin\legal-bot.exe" intake "%MATTER%"

echo.
echo  Output: matters\%MATTER%\knowledge\
echo.
pause
