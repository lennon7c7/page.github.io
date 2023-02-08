@echo off

rem -------------------------------------------------------------
rem  Docker init script for Windows
rem -------------------------------------------------------------


@REM step1: init env
set "_path=%~dp0"
for %%a in ("%_path%") do (set "p_dir=%%~dpa")
for %%a in (%p_dir:~0,-1%) do (set p2_dir=%%~dpa&&set p2_folder=%%~nxa)
for %%a in (%p2_dir:~0,-1%) do (set p3_dir=%%~dpa&&set p3_folder=%%~nxa)
for %%a in (%p3_dir:~0,-1%) do (set p4_dir=%%~dpa&&set p4_folder=%%~nxa)
for %%a in (%p4_dir:~0,-1%) do (set p5_dir=%%~dpa&&set p5_folder=%%~nxa)
for %%a in (%p5_dir:~0,-1%) do (set p6_dir=%%~dpa&&set p6_folder=%%~nxa)

@REM step2: docker delete
docker rm -f "%p5_folder%"
docker image rm "%p5_folder%"


@REM step3: docker start
@cd /d %p_dir%
docker build -t "%p5_folder%" -f Dockerfile .
docker run -d --restart=always --name "%p5_folder%" -v %p4_dir%:/var/www/html -p 9980:80 "%p5_folder%"
