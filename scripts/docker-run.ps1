# Production run script for Windows
Write-Host "🚀 Starting Surplus Supper in production mode..." -ForegroundColor Green

# Build and start services
docker-compose up --build -d

Write-Host "✅ Services started successfully!" -ForegroundColor Green
Write-Host "📊 Application is running at: http://localhost:8080" -ForegroundColor Cyan
Write-Host "🗄️  Database is running at: localhost:5432" -ForegroundColor Cyan
Write-Host ""
Write-Host "🔧 Useful commands:" -ForegroundColor Yellow
Write-Host "  View logs: docker-compose logs -f"
Write-Host "  Stop services: docker-compose down"
Write-Host "  Rebuild: docker-compose up --build -d"
