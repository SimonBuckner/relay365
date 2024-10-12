# Build the service exe and install it as a service.

$rootDir = Get-Location
$svcName = "relay365"
$exeName = "relay365svc.exe"
$exe = Join-Path -Path $rootDir -ChildPath "build\$exeName"

Write-Output "Checking if service is already installed"
$svc = Get-Service -Name $svcName -ErrorAction SilentlyContinue
if ($null -eq $svc) {
    Write-Output "Service note installed. Creating the service"
    New-Service -Name $svcName -BinaryPathName $exe
}

Write-Output "Getting service status"
$running = Get-Service -Name $svcName
if ($running.Status -eq 'Running') {
    Write-Output "Service is running. Stopping the service"
    Stop-Service -Name $svcName -Force
}

Write-Output "Building the new service exe"
go build -ldflags "-s -w" -o "./build/$exeName" "./cmd/relay365svc/relay365svc.go" 

Write-Output "Starting the service"
Start-Service -Name $svcName

