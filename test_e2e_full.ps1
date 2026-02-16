# Helper function for requests
function Invoke-PizzaRequest {
    param (
        [string]$Method,
        [string]$Uri,
        [hashtable]$Body = @{},
        [string]$Token = ""
    )
    
    $Headers = @{ "Content-Type" = "application/json" }
    if ($Token) { $Headers["Authorization"] = "Bearer $Token" }
    
    try {
        $params = @{
            Method = $Method
            Uri = "http://localhost:8080/api/v1$Uri"
            Headers = $Headers
        }
        if ($Method -ne "GET" -and $Body.Keys.Count -gt 0) {
            $params["Body"] = ($Body | ConvertTo-Json -Depth 10)
        }
        
        $response = Invoke-RestMethod @params
        return $response
    } catch {
        Write-Host "Error calling $Uri" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        if ($_.Exception.Response) {
            $stream = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($stream)
            Write-Host $reader.ReadToEnd() -ForegroundColor Red
        }
        throw
    }
}

Write-Host "--- STARTING E2E TEST ---" -ForegroundColor Cyan

# 1. Register Users
Write-Host "`n1. Registering Users..." -ForegroundColor Yellow
$users = @(
    @{ email="manager@test.com"; pass="password123"; role="manager"; name="Manager Mike" },
    @{ email="chef@test.com";    pass="password123"; role="chef";    name="Chef Gordon" },
    @{ email="courier@test.com"; pass="password123"; role="courier"; name="Courier Sam" },
    @{ email="client@test.com";  pass="password123"; role="client";  name="Client Alice" }
)

$tokens = @{}
$userIds = @{}

foreach ($u in $users) {
    try {
        Invoke-PizzaRequest -Method POST -Uri "/auth/register" -Body @{ email=$u.email; password=$u.pass; name=$u.name; role=$u.role } | Out-Null
        Write-Host "Registered $($u.role): $($u.email)" -ForegroundColor Green
    } catch {
        Write-Host "User $($u.email) might already exist, trying login..." -ForegroundColor DarkGray
    }
    
    $login = Invoke-PizzaRequest -Method POST -Uri "/auth/login" -Body @{ email=$u.email; password=$u.pass }
    $tokens[$u.role] = $login.data.token
    $userIds[$u.role] = $login.data.user_id
    Write-Host "Logged in $($u.role) (ID: $($login.data.user_id))" -ForegroundColor Green
}

# 2. Create Product (Manager)
Write-Host "`n2. Creating Product..." -ForegroundColor Yellow
$prodBody = @{
    name = "Pepperoni Feast"
    description = "Double pepperoni"
    price = 15.50
    category_id = 1
    image_url = "http://img.com/pep.png"
}
$product = Invoke-PizzaRequest -Method POST -Uri "/catalog/products" -Body $prodBody -Token $tokens["manager"]
$productId = $product.data.id
Write-Host "Created Product Response: $($product | ConvertTo-Json -Depth 5)" -ForegroundColor Cyan
Write-Host "Created Product: $($product.data.name) ($productId)" -ForegroundColor Green

# 3. Create Order (Client)
Write-Host "`n3. Creating Order..." -ForegroundColor Yellow
$orderBody = @{
    items = @( @{ product_id = $productId; quantity = 2 } )
    address = @{ city="New York"; street="5th Ave"; house="10"; apartment="1" }
}
$order = Invoke-PizzaRequest -Method POST -Uri "/orders" -Body $orderBody -Token $tokens["client"]
$orderId = $order.data.order_id
Write-Host "Created Order Response: $($order | ConvertTo-Json -Depth 5)" -ForegroundColor Cyan
Write-Host "Created Order: $orderId (Status: $($order.data.status))" -ForegroundColor Green

# 4. Pay Order (Client)
Write-Host "`n4. Paying Order..." -ForegroundColor Yellow
Invoke-PizzaRequest -Method POST -Uri "/orders/$orderId/pay" -Token $tokens["client"] | Out-Null
Write-Host "Order Paid" -ForegroundColor Green

# 5. Kitchen Flow (Chef)
Write-Host "`n5. Kitchen Flow..." -ForegroundColor Yellow
Start-Sleep -Seconds 1 # Wait for sync

$tickets = Invoke-PizzaRequest -Method GET -Uri "/kitchen/tickets" -Token $tokens["chef"]
$myTicket = $tickets.data.tickets | Where-Object { $_.order_id -eq $orderId }

if (!$myTicket) { Write-Error "Ticket not found in kitchen!" }
Write-Host "Found Ticket: $($myTicket.id)" -ForegroundColor Green

# Start Cooking
Invoke-PizzaRequest -Method PATCH -Uri "/kitchen/tickets/$($myTicket.id)/status" -Body @{ status="cooking" } -Token $tokens["chef"] | Out-Null
Write-Host "Status -> Cooking" -ForegroundColor Green

# Mark Ready
Invoke-PizzaRequest -Method PATCH -Uri "/kitchen/tickets/$($myTicket.id)/status" -Body @{ status="ready" } -Token $tokens["chef"] | Out-Null
Write-Host "Status -> Ready" -ForegroundColor Green

# 6. Logistics Flow (Courier)
Write-Host "`n6. Logistics Flow..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

$deliveries = Invoke-PizzaRequest -Method GET -Uri "/logistics/deliveries" -Token $tokens["courier"]
$myDelivery = $deliveries.data.deliveries | Where-Object { $_.order_id -eq $orderId }

if (!$myDelivery) { Write-Error "Delivery not found!" }
Write-Host "Found Delivery for Order: $($myDelivery.order_id)" -ForegroundColor Green

# Assign Courier
Invoke-PizzaRequest -Method POST -Uri "/logistics/orders/$orderId/assign" -Body @{ courier_id = $userIds["courier"] } -Token $tokens["courier"] | Out-Null
Write-Host "Assigned Courier" -ForegroundColor Green

# Start Delivery
Invoke-PizzaRequest -Method PATCH -Uri "/logistics/orders/$orderId/status" -Body @{ status="delivering" } -Token $tokens["courier"] | Out-Null
Write-Host "Status -> Delivering" -ForegroundColor Green

# Complete Delivery
Invoke-PizzaRequest -Method PATCH -Uri "/logistics/orders/$orderId/status" -Body @{ status="completed" } -Token $tokens["courier"] | Out-Null
Write-Host "Status -> Completed" -ForegroundColor Green

# 7. Check Analytics (Manager)
Write-Host "`n7. Checking Analytics..." -ForegroundColor Yellow
# Since we used "GENERAL_STORE" in the code for now, we might not see per-manager KPI yet, 
# but let's check order status or general analytics if available.
# Actually, let's verify the order is completed in the system.
$finalOrder = Invoke-PizzaRequest -Method GET -Uri "/orders/$orderId" -Token $tokens["manager"]
if ($finalOrder.data.status -eq "completed") {
    Write-Host "Final Order Status Verified: COMPLETED" -ForegroundColor Green
} else {
    Write-Error "Final status mismatch: $($finalOrder.data.status)"
}

Write-Host "`n--- E2E TEST PASSED ---" -ForegroundColor Cyan
