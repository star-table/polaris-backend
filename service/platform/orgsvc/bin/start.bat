cd %~dp0
cd ..

set env=%1
if "%env%" == "" (
   set env=local
)
echo %env%
go mod tidy
go run main.go --env=%env%