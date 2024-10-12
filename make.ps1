
param (
    [ValidateSet("Svc", "Cli", "Both")]
    [string]$Build,

    [switch]$RunCli = $false,
    [switch]$InstallService = $False,

    [switch]$UnistallService = $False
)

$appName = "relay365"
$svcFname = "$($appName)_service.exe"
$cliFname = "$($appName).exe"
$svcDisplayName = "Local MS 365 SMTP to Graph API Relay service"
$serviceDir = "C:\Program Files\SimonBuckner"
$servicePath = "$serviceDir\$appName"


if ($Build -eq "Svc" -or $Build -eq "Both") {
    Write-Output "Building $svcFname"
    go build -ldflags "-s -w" -o "./build/$svcFname" "./cmd/relay365svc/relay365svc.go" 
}

if ($Build -eq "Cli" -or $Build -eq "Both") {
    Write-Output "Building $cliFname"
    go build -ldflags "-s -w" -o "./build/$cliFname" "./cmd/relay365cli/relay365cli.go" 
}

if ($RunCli) {
    go run "./cmd/relay365cli/relay365cli.go" 
}
if ($UnistallService) {

    $isAdmin = ( [Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if ($isAdmin -eq $false) {
        Write-Error "Please run as admin" -ErrorAction Stop
    }

    $svc = Get-Service -Name $appName -ErrorAction SilentlyContinue

    $ErrorActionPreference = 'Stop'
    if ($null -ne $svc) {
        Write-Host "Stopping service"
        Stop-Service -Name $appName -Force
    }

    if ($null -ne $svc) {
        Write-Output "Removing the service $svcDisplayName"
        Remove-Service -Name $appName
    }

    # TODO // Delete the exe and dir
}

if ($InstallService) {

    $isAdmin = ( [Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if ($isAdmin -eq $false) {
        Write-Error "Please run as admin" -ErrorAction Stop
    }

    $svc = Get-Service -Name $appName -ErrorAction SilentlyContinue

    $ErrorActionPreference = 'Stop'
    if ($null -ne $svc) {
        Write-Host "Stopping service"
        Stop-Service -Name $appName -Force
    }
    
    if (-not (Test-Path -Path $servicePath)) {
        Write-Output "Creating service directory $servicePath"
        New-Item -Path $servicePath -Type Directory -Force | Out-Null
    }
   
    Write-Output "Install $svcFname into $serviceDir"
    Copy-Item -Path "./build/$svcFname" -Destination $servicePath -Force

    if ($null -eq $svc) {
        Write-Output "Creating the service $svcDisplayName"
        New-Service -Name $appName -DisplayName $svcDisplayName -BinaryPathName "$servicePath\$svcFname"
    }

    Write-Output "Starting the service"
    Start-Service -Name $appName
    
}


