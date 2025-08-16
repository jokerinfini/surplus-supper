#!/bin/bash

# 🚀 Surplus Supper Deployment Script
# This script helps you deploy your project to free hosting platforms

echo "🍽️ Surplus Supper Deployment Script"
echo "=================================="

# Check if git is initialized
if [ ! -d ".git" ]; then
    echo "❌ Git repository not found. Please initialize git first:"
    echo "   git init"
    echo "   git add ."
    echo "   git commit -m 'Initial commit'"
    echo "   git remote add origin https://github.com/YOUR_USERNAME/surplus-supper.git"
    echo "   git push -u origin main"
    exit 1
fi

# Check if changes are committed
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  You have uncommitted changes. Please commit them first:"
    echo "   git add ."
    echo "   git commit -m 'Prepare for deployment'"
    echo "   git push"
    exit 1
fi

echo "✅ Git repository is ready"

# Create .env.local for frontend
echo "📝 Creating frontend environment file..."
cat > frontend-next/.env.local << EOF
NEXT_PUBLIC_API_URL=https://your-backend-url.railway.app
EOF

echo "✅ Created frontend-next/.env.local"
echo "⚠️  Remember to update NEXT_PUBLIC_API_URL with your actual backend URL"

# Create .env for backend
echo "📝 Creating backend environment file..."
cat > backend/.env << EOF
DATABASE_URL=postgresql://username:password@host:port/database
PORT=8080
JWT_SECRET=your-secret-key-here
CORS_ORIGIN=https://your-frontend-url.vercel.app
EOF

echo "✅ Created backend/.env"
echo "⚠️  Remember to update the values with your actual deployment URLs"

echo ""
echo "🎯 Next Steps:"
echo "=============="
echo ""
echo "1. 🚀 Deploy Backend to Railway:"
echo "   - Go to https://railway.app"
echo "   - Create account and new project"
echo "   - Add PostgreSQL database"
echo "   - Deploy backend service"
echo "   - Copy the backend URL"
echo ""
echo "2. 🎨 Deploy Frontend to Vercel:"
echo "   - Go to https://vercel.com"
echo "   - Create account and new project"
echo "   - Import your GitHub repository"
echo "   - Set root directory to 'frontend-next'"
echo "   - Update NEXT_PUBLIC_API_URL environment variable"
echo ""
echo "3. 🔧 Update Environment Variables:"
echo "   - Update backend/.env with your Railway database URL"
echo "   - Update frontend-next/.env.local with your Railway backend URL"
echo "   - Update CORS_ORIGIN in backend with your Vercel frontend URL"
echo ""
echo "4. 🧪 Test Your Deployment:"
echo "   - Test backend: curl https://your-backend-url/health"
echo "   - Test frontend: Visit your Vercel URL"
echo "   - Test authentication flow"
echo ""
echo "📚 For detailed instructions, see DEPLOYMENT_GUIDE.md"
echo ""
echo "🎉 Good luck with your deployment!"
