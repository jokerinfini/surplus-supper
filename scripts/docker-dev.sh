#!/bin/bash

# Development run script
echo "🔧 Starting Surplus Supper in development mode..."

# Build and start services with hot reload
docker-compose -f docker-compose.dev.yml up --build -d

echo "✅ Development services started successfully!"
echo "📊 Application is running at: http://localhost:8080"
echo "🗄️  Database is running at: localhost:5432"
echo "🔄 Hot reload is enabled - changes will auto-restart the server"
echo ""
echo "�� Useful commands:"
echo "  View logs: docker-compose -f docker-compose.dev.yml logs -f"
echo "  Stop services: docker-compose -f docker-compose.dev.yml down"
echo "  Rebuild: docker-compose -f docker-compose.dev.yml up --build -d"
```

