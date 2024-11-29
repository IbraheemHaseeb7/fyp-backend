@CreatedAt echo off
SETLOCAL

cls

REM Checkinf if an arg was passed
IF "%~1"=="" (
	echo No argument was passed, FAILED !!!
	exit /b 1
)

REM Checking the arg and executing commands

REM starting all services
IF /I "%~1"=="start" (
	echo Starting auth service...
	go build -o auth.exe main.go
	IF ERRORLEVEL 1 (
		echo Auth Build failed...
		exit /b 1
	)
	start /b auth.exe

	echo Starting db service...
	cd micro
	go build -o db.exe main.go
	if ERRORLEVEL 1 (
		echo DB build failed...
		exit /b 1
	)
	start /b db.exe

	echo Successfully ran both services...
)

REM starting the auth service
IF /I "%~1"=="startAuth" (
	echo Starting auth service...
	go build -o auth.exe main.go
	IF ERRORLEVEL 1 (
		echo Auth Build failed...
		exit /b 1
	)
	start /b auth.exe
)

REM starting the db service
IF /I "%~1"=="startDB" (
	echo Starting db service...
	cd micro
	go build -o db.exe main.go
	if ERRORLEVEL 1 (
		echo DB build failed...
		exit /b 1
	)
	start /b db.exe
)

REM killing all the processes
IF /I "%~1"=="stop" (
	echo Stopping all processes...
	taskkill /IM auth.exe /F
	taskkill /IM db.exe /F
)

REM killing the auth process
IF /I "%~1"=="stopAuth" (
	echo Stopping auth process...
	taskkill /IM auth.exe /F
)

REM killing the db process
IF /I "%~1"=="stopDB" (
	echo Stopping db process...
	taskkill /IM db.exe /F
)

REM restarting all the services
IF /I "%~1"=="restart" (

	echo Stopping all processes...
	taskkill /IM auth.exe /F
	taskkill /IM db.exe /F

	echo Starting auth service...
	go build -o auth.exe main.go
	IF ERRORLEVEL 1 (
		echo Auth Build failed...
		exit /b 1
	)
	start /b auth.exe

	echo Starting db service...
	cd micro
	go build -o db.exe main.go
	if ERRORLEVEL 1 (
		echo DB build failed...
		exit /b 1
	)
	start /b db.exe

	echo Successfully ran both services...
)

