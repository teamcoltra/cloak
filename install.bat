@echo off
SETLOCAL EnableExtensions

:: Check if Go is installed by looking for the GOBIN environment variable
if "%GOBIN%"=="" (
    echo Go is not detected on this system. Go is required to build Cloak.
    exit /b 1
)

:: Build the project
echo Building Cloak...
go build -o build/cloak.exe .

:: Check if the build was successful
if not exist "build/cloak.exe" (
    echo Failed to build Cloak.
    exit /b 1
)

:: Create target directories
echo Creating target directories...
if not exist "C:\Program Files\Cloak" (
    mkdir "C:\Program Files\Cloak"
)
if not exist "C:\Program Files\Cloak\www" (
    mkdir "C:\Program Files\Cloak\www"
)

:: Move the binary and resources
echo Moving binary and resources...
move build/cloak.exe "C:\Program Files\Cloak"

:: Assuming stuff/www contains static web files
xcopy stuff\www "C:\Program Files\Cloak\www" /E /H /C /I

:: Move other necessary files (assuming they are in the stuff directory)
move stuff\domains.txt "C:\Program Files\Cloak"
move stuff\dictionary.txt "C:\Program Files\Cloak"
move stuff\cloak.cfg "C:\Program Files\Cloak"

echo Installation complete.
ENDLOCAL
