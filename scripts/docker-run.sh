#!/bin/bash

# Production run script
echo "🚀 Starting Surplus Supper in production mode..."

# Build and start services
docker-compose up --build -d

echo "✅ Services started successfully!"
echo "📊 Application is running at: http://localhost:8080"
echo "🗄️  Database is running at: localhost:5432"
echo ""
echo "�� Useful commands:"
echo "  View logs: docker-compose logs -f"
echo "  Stop services: docker-compose down"
echo "  Rebuild: docker-compose up --build -d"
