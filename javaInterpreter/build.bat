@echo off
if not exist "out" mkdir out
dir /s /B *.java > sources.txt
javac -d out @sources.txt
del sources.txt
echo Build complete.
pause
