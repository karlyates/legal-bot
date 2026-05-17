@echo off
setlocal

REM Open the latest review output

set MATTER=%1
if "%MATTER%"=="" set MATTER=TRO-MTE-GAL-SM

set "SCRIPT_DIR=%~dp0"
set "ROOT_DIR=%SCRIPT_DIR%.."
set "OUTPUT_DIR=%ROOT_DIR%\matters\%MATTER%\output"

if exist "%OUTPUT_DIR%\final_review_packet.md" (
    start "" "%OUTPUT_DIR%\final_review_packet.md"
) else (
    echo No review packet found. Run a review first.
    echo   scripts\run-review.bat %MATTER%
    pause
)
