@echo off
SETLOCAL EnableDelayedExpansion

SET "PROJECT_PATH=%~dp0"

SET "GOOS_LIST=windows linux"
SET "GOARCH_LIST=amd64"

FOR %%A IN (%GOOS_LIST%) DO (
    FOR %%B IN (%GOARCH_LIST%) DO (
        SET "SUFFIX="
        IF "%%A"=="windows" (
            SET "SUFFIX=.exe"
        )
        ECHO Building for %%A/%%B...
        SET "GOOS=%%A"
        SET "GOARCH=%%B"
        go build -o "output\webbenchmark_%%A_%%B!SUFFIX!" "%PROJECT_PATH%cmd\main.go"
    )
)

ENDLOCAL
@echo on
