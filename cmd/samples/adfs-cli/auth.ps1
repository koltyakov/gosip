$ConfigPath = "./config/private.onprem-wap-adfs.json";

$ConfigJson = Get-Content -Raw -Path $ConfigPath | ConvertFrom-Json;
$SiteUrl = $ConfigJson.siteUrl;
$Domain = ([System.Uri]$SiteUrl).Host -replace '^www\.';

$SpAuthRead = "go run ./cmd/adfs-cli/main.go -configPath $ConfigPath";
$Cookies = Invoke-Expression $SpAuthRead | ConvertFrom-Json;

$Session = New-Object Microsoft.PowerShell.Commands.WebRequestSession;

ForEach($Prop in $Cookies.PSObject.Properties)
{
  $Cookie = New-Object System.Net.Cookie;
  $CookieName = $Prop.Name;
  $Cookie.Name = $CookieName;
  $Cookie.Value = $Cookies.$CookieName;
  $Cookie.Domain = $Domain;
  $Session.Cookies.Add($Cookie);
}

$Response = Invoke-WebRequest "$SiteUrl/_api/web?$select=Title" `
  -WebSession $Session `
  -Method "GET" `
  -Headers @{"accept"="application/json;odata=verbose"};

$Data = $Response | ConvertFrom-Json;

Write-Host $Data.d.Title;