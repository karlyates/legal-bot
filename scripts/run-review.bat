@echo off
setlocal

REM Legal-Bot Review Pipeline
REM Reviews a draft against the knowledge layer

set MATTER=%1
if "%MATTER%"=="" set MATTER=TRO-MTE-GAL-SM

echo.
echo  ============================================
echo   Legal-Bot Draft Review Pipeline
echo   Matter: %MATTER%
echo  ============================================
echo.

set "SCRIPT_DIR=%~dp0"
set "ROOT_DIR=%SCRIPT_DIR%.."

REM Check for knowledge layer
if not exist "%ROOT_DIR%\matters\%MATTER%\knowledge\document_register.csv" (
    echo ERROR: No knowledge layer found for matter %MATTER%.
    echo Run intake first: scripts\run-intake.bat %MATTER%
    exit /b 1
)

REM Check for draft files
dir /b "%ROOT_DIR%\matters\%MATTER%\drafts\*.txt" "%ROOT_DIR%\matters\%MATTER%\drafts\*.md" >nul 2>&1
if errorlevel 1 (
    echo ERROR: No draft files found in matters\%MATTER%\drafts\
    echo Drop your attorney draft there first.
    exit /b 1
)

REM Build if needed
if not exist "%ROOT_DIR%\go\bin\legal-bot.exe" (
    echo Building legal-bot CLI...
    pushd "%ROOT_DIR%\go"
    go build -o bin\legal-bot.exe .\cmd\legal-bot
    popd
)

REM Run review
"%ROOT_DIR%\go\bin\legal-bot.exe" review "%MATTER%"

echo.
echo  Output: matters\%MATTER%\output\final_review_packet.md
echo.
pause
