# Development run script for Windows
Write-Host "🔧 Starting Surplus Supper in development mode..." -ForegroundColor Green

# Build and start services with hot reload
docker-compose -f docker-compose.dev.yml up --build -d

Write-Host "✅ Development services started successfully!" -ForegroundColor Green
Write-Host "📊 Application is running at: http://localhost:8080" -ForegroundColor Cyan
Write-Host "🗄️  Database is running at: localhost:5432" -ForegroundColor Cyan
Write-Host "🔄 Hot reload is enabled - changes will auto-restart the server" -ForegroundColor Yellow
Write-Host ""
Write-Host "🔧 Useful commands:" -ForegroundColor Yellow
Write-Host "  View logs: docker-compose -f docker-compose.dev.yml logs -f"
Write-Host "  Stop services: docker-compose -f docker-compose.dev.yml down"
Write-Host "  Rebuild: docker-compose -f docker-compose.dev.yml up --build -d"
