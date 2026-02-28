# Configuration
$loginUrl = 'http://localhost:8080/login'
$apiUrl = 'http://localhost:8080/runner'
$username = 'admin'
$password = 'admin'

function Get-BasicAuthHeader($user, $pass) {
    $pair = "$($user):$($pass)"
    # Use UTF8 encoding so non-ASCII characters (eg. ñ) are preserved
    $bytes = [System.Text.Encoding]::UTF8.GetBytes($pair)
    $b64 = [Convert]::ToBase64String($bytes)
    return "Basic $b64"
}

Write-Host "Obtenemos un Token con $loginUrl ..."
$authHeader = @{ Authorization = (Get-BasicAuthHeader $username $password) }
try {
    $token = Invoke-RestMethod -Uri $loginUrl -Method Post -Headers $authHeader -Body '' -ErrorAction Stop
} catch {
    Write-Error "No pudimos obtener el token: $_"
    exit 1
}

# sino trenemos un token terminamos aquí
if (-not $token) {
    Write-Error "No conseguimos el token. Terminamos."
    exit 1
}

Write-Host "Got token: $token"

$tokenHeader = @{ Token = $token; 'Content-Type' = 'application/json' }

# Creamos seis runners
$usuarios = @(
    @{ first_name = 'Usuario'; last_name = 'Uno'; age = 56; is_active = $true; country = 'España'; personal_best = ''; season_best = '' },
    @{ first_name = 'Usuario'; last_name = 'Dos'; age = 52; is_active = $true; country = 'España'; personal_best = ''; season_best = '' },
    @{ first_name = 'Usuario'; last_name = 'Tres'; age = 21; is_active = $true; country = 'España'; personal_best = ''; season_best = '' },
    @{ first_name = 'Usuario'; last_name = 'Cuatro'; age = 19; is_active = $true; country = 'España'; personal_best = ''; season_best = '' },
    @{ first_name = 'Usuario'; last_name = 'Cinco'; age = 17; is_active = $true; country = 'España'; personal_best = ''; season_best = '' },
    @{ first_name = 'Usuario'; last_name = 'Seis'; age = 14; is_active = $true; country = 'España'; personal_best = ''; season_best = '' }
)

Write-Host "Creamos $($usuarios.Count) runners con $apiUrl ..."

foreach ($u in $usuarios) {
    $body = $u | ConvertTo-Json -Depth 4
    try {
        $resp = Invoke-RestMethod -Uri $apiUrl -Method Post -Headers $tokenHeader -Body $body -ErrorAction Stop
        Write-Host "Creado el runner: $($u.last_name) -> Respuesta: $($resp)"
    } catch {
        Write-Warning "No se pudo crear el runner $($u.last_name): $_ $resp"
    }
}

Write-Host "Se han creado los payloads"

while ($true) {
    try {
        $res = Invoke-RestMethod -Uri $apiUrl -Method Get -Headers @{ Token = $token } -ErrorAction Stop
        $ts = (Get-Date).ToString('o')
        Write-Host "[$ts] GET ok (items: $((($res | Measure-Object).Count)))"
    } catch {
        $ts = (Get-Date).ToString('o')
        Write-Warning "[$ts] GET ko: $_"
    }

    # delay aleatorio
    $ms = Get-Random -Minimum 500 -Maximum 5000
    Start-Sleep -Milliseconds $ms
}
