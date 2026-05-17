@echo off
REM Legacy VernHole wrapper.
REM
REM This wrapper is not part of the current Legal-Bot legal workflow.
REM The current operator-facing legal surface is guide, intake, review,
REM and workflow. See docs/legacy_surfaces.md.

echo This is a legacy VernHole wrapper and is not part of the current Legal-Bot legal workflow. 1>&2
echo Use legal-bot tui only if you explicitly want the legacy utility UI. 1>&2
echo See docs/legacy_surfaces.md. 1>&2
exit /b 1
