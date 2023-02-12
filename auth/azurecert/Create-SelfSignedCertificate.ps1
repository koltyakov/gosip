#Requires -RunAsAdministrator
<#
.SYNOPSIS
Creates a Self Signed Certificate for use in server to server authentication
.DESCRIPTION
.EXAMPLE
PS C:\> .\Create-SelfSignedCertificate.ps1 -CommonName "MyCert" -StartDate 2015-11-21 -EndDate 2017-11-21
This will create a new self signed certificate with the common name "CN=MyCert". During creation you will be asked to provide a password to protect the private key.
.EXAMPLE
PS C:\> .\Create-SelfSignedCertificate.ps1 -CommonName "MyCert" -StartDate 2015-11-21 -EndDate 2017-11-21 -Password (ConvertTo-SecureString -String "MyPassword" -AsPlainText -Force)
This will create a new self signed certificate with the common name "CN=MyCert". The password as specified in the Password parameter will be used to protect the private key
.EXAMPLE
PS C:\> .\Create-SelfSignedCertificate.ps1 -CommonName "MyCert" -StartDate 2015-11-21 -EndDate 2017-11-21 -Force
This will create a new self signed certificate with the common name "CN=MyCert". During creation you will be asked to provide a password to protect the private key. If there is already a certificate with the common name you specified, it will be removed first.
#>
Param (
  [Parameter(Mandatory=$true)]
  [string]$CommonName,

  [Parameter(Mandatory=$true)]
  [DateTime]$StartDate,

  [Parameter(Mandatory=$true)]
  [DateTime]$EndDate,

  [Parameter(Mandatory=$false, HelpMessage="Will overwrite existing certificates")]
  [Switch]$Force,

  [Parameter(Mandatory=$false)]
  [SecureString]$Password
)

function CreateSelfSignedCertificate() {
  #Remove and existing certificates with the same common name from personal and root stores
  #Need to be very wary of this as could break something
  if ($CommonName.ToLower().StartsWith("cn=")) {
    # Remove CN from common name
    $CommonName = $CommonName.Substring(3)
  }
  $certs = Get-ChildItem -Path Cert:\LocalMachine\my | Where-Object { $_.Subject -eq "CN=$CommonName" }
  if ($certs -ne $null -and $certs.Length -gt 0) {
    if ($Force) {
      foreach ($c in $certs) {
        Remove-Item $c.PSPath
      }
    } else {
      Write-Host -ForegroundColor Red "One or more certificates with the same common name (CN=$CommonName) are already located in the local certificate store. Use -Force to remove them";
      return $false
    }
  }

  $name = New-Object -com "X509Enrollment.CX500DistinguishedName.1"
  $name.Encode("CN=$CommonName", 0)

  $key = New-Object -com "X509Enrollment.CX509PrivateKey.1"
  $key.ProviderName = "Microsoft RSA SChannel Cryptographic Provider"
  $key.KeySpec = 1
  $key.Length = 2048
  $key.SecurityDescriptor = "D:PAI(A;;0xd01f01ff;;;SY)(A;;0xd01f01ff;;;BA)(A;;0x80120089;;;NS)"
  $key.MachineContext = 1
  $key.ExportPolicy = 1 # This is required to allow the private key to be exported
  $key.Create()

  $serverauthoid = New-Object -com "X509Enrollment.CObjectId.1"
  $serverauthoid.InitializeFromValue("1.3.6.1.5.5.7.3.1") # Server Authentication
  $ekuoids = New-Object -com "X509Enrollment.CObjectIds.1"
  $ekuoids.add($serverauthoid)
  $ekuext = New-Object -com "X509Enrollment.CX509ExtensionEnhancedKeyUsage.1"
  $ekuext.InitializeEncode($ekuoids)

  $cert = New-Object -com "X509Enrollment.CX509CertificateRequestCertificate.1"
  $cert.InitializeFromPrivateKey(2, $key, "")
  $cert.Subject = $name
  $cert.Issuer = $cert.Subject
  $cert.NotBefore = $StartDate
  $cert.NotAfter = $EndDate
  $cert.X509Extensions.Add($ekuext)
  $cert.Encode()

  $enrollment = New-Object -com "X509Enrollment.CX509Enrollment.1"
  $enrollment.InitializeFromRequest($cert)
  $certdata = $enrollment.CreateRequest(0)
  $enrollment.InstallResponse(2, $certdata, 0, "")
  return $true
}

function ExportPFXFile() {
  if ($CommonName.ToLower().StartsWith("cn=")) {
    # Remove CN from common name
    $CommonName = $CommonName.Substring(3)
  }
  if ($Password -eq $null) {
    $Password = Read-Host -Prompt "Enter Password to protect private key" -AsSecureString
  }
  $cert = Get-ChildItem -Path Cert:\LocalMachine\my | Where-Object { $_.Subject -eq "CN=$CommonName" }
  Export-PfxCertificate -Cert $cert -Password $Password -FilePath "$($CommonName).pfx"
  Export-Certificate -Cert $cert -Type CERT -FilePath "$CommonName.cer"
}

function RemoveCertsFromStore() {
  # Once the certificates have been been exported we can safely remove them from the store
  if ($CommonName.ToLower().StartsWith("cn=")) {
    # Remove CN from common name
    $CommonName = $CommonName.Substring(3)
  }
  $certs = Get-ChildItem -Path Cert:\LocalMachine\my | Where-Object { $_.Subject -eq "CN=$CommonName" }
  foreach ($c in $certs) {
    remove-item $c.PSPath
  }
}

if (CreateSelfSignedCertificate) {
  ExportPFXFile
  RemoveCertsFromStore
}