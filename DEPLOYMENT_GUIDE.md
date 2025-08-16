# 🚀 Deployment Guide - Free Hosting for Portfolio

This guide will help you deploy your Surplus Supper project for free to showcase it in your resume.

## 🎯 **Recommended Approach: Vercel + Railway**

### **Frontend (Next.js) → Vercel**
### **Backend (Go) + Database → Railway**

---

## 📋 **Step 1: Prepare Your Code**

### 1.1 Update Environment Variables

**Frontend (`frontend-next/.env.local`):**
```env
NEXT_PUBLIC_API_URL=https://your-backend-url.railway.app
```

**Backend (`backend/.env`):**
```env
DATABASE_URL=postgresql://username:password@host:port/database
PORT=8080
JWT_SECRET=your-secret-key
CORS_ORIGIN=https://your-frontend-url.vercel.app
```

### 1.2 Update API Client

**File: `frontend-next/src/lib/auth.ts`**
```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
```

---

## 🚀 **Step 2: Deploy Backend to Railway**

### 2.1 Create Railway Account
1. Go to [railway.app](https://railway.app)
2. Sign up with GitHub
3. Get $5 free credit monthly

### 2.2 Deploy Backend
1. **Create New Project**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository

2. **Add Database**
   - Click "New" → "Database" → "PostgreSQL"
   - Railway will provide connection string

3. **Configure Backend Service**
   - Click "New" → "GitHub Repo"
   - Select your repo
   - Set root directory to `backend`
   - Railway will auto-detect Go

4. **Set Environment Variables**
   ```
   DATABASE_URL=postgresql://... (from Railway)
   PORT=8080
   JWT_SECRET=your-secret-key-here
   CORS_ORIGIN=https://your-frontend-url.vercel.app
   ```

5. **Deploy**
   - Railway will automatically build and deploy
   - Get your backend URL (e.g., `https://surplus-supper-backend.railway.app`)

### 2.3 Test Backend
```bash
curl https://your-backend-url.railway.app/health
```

---

## 🎨 **Step 3: Deploy Frontend to Vercel**

### 3.1 Create Vercel Account
1. Go to [vercel.com](https://vercel.com)
2. Sign up with GitHub
3. Free tier includes custom domains

### 3.2 Deploy Frontend
1. **Import Project**
   - Click "New Project"
   - Import your GitHub repository
   - Set root directory to `frontend-next`

2. **Configure Build Settings**
   - Framework Preset: Next.js
   - Build Command: `npm run build`
   - Output Directory: `.next`
   - Install Command: `npm install`

3. **Set Environment Variables**
   ```
   NEXT_PUBLIC_API_URL=https://your-backend-url.railway.app
   ```

4. **Deploy**
   - Vercel will automatically deploy
   - Get your frontend URL (e.g., `https://surplus-supper.vercel.app`)

### 3.3 Update Backend CORS
Update your backend CORS origin to include your Vercel URL:
```
CORS_ORIGIN=https://surplus-supper.vercel.app
```

---

## 🔧 **Alternative: Render (Free Tier)**

### Backend Deployment on Render
1. Go to [render.com](https://render.com)
2. Create account
3. **New Web Service**
   - Connect GitHub repo
   - Root directory: `backend`
   - Build Command: `go build -o backend`
   - Start Command: `./backend`
   - Environment: Go

4. **Add PostgreSQL Database**
   - New → PostgreSQL
   - Get connection string
   - Add to environment variables

---

## 🌐 **Step 4: Custom Domain (Optional)**

### Vercel Custom Domain
1. Go to your Vercel project
2. Settings → Domains
3. Add your domain
4. Configure DNS records

### Railway Custom Domain
1. Go to your Railway project
2. Settings → Domains
3. Add custom domain
4. Configure DNS

---

## 📊 **Step 5: Monitoring & Maintenance**

### Health Checks
- **Backend**: `https://your-backend-url/health`
- **Frontend**: Vercel provides built-in monitoring

### Logs
- **Railway**: Built-in logging dashboard
- **Vercel**: Function logs and analytics

---

## 💰 **Cost Breakdown (Free Tier)**

### Railway
- **Backend**: $5/month credit (sufficient for small projects)
- **Database**: Included in credit
- **Bandwidth**: Generous limits

### Vercel
- **Frontend**: Completely free
- **Bandwidth**: 100GB/month
- **Builds**: 6000 minutes/month

### Total Cost: $0-5/month

---

## 🎯 **Resume Showcase Tips**

### 1. **Live Demo Link**
```
🌐 Live Demo: https://surplus-supper.vercel.app
🔧 Backend API: https://surplus-supper-backend.railway.app
📚 GitHub: https://github.com/yourusername/surplus-supper
```

### 2. **Technical Highlights**
- **Full-Stack**: Go backend + Next.js frontend
- **Database**: PostgreSQL with migrations
- **Authentication**: JWT with bcrypt
- **Deployment**: Production-ready with CI/CD
- **Modern UI**: Glassmorphism design with animations

### 3. **Features to Highlight**
- ✅ User authentication system
- ✅ Location-based restaurant search
- ✅ Real-time geolocation
- ✅ Responsive design
- ✅ Modern UI/UX with animations
- ✅ Production deployment

---

## 🚨 **Troubleshooting**

### Common Issues

**1. CORS Errors**
```go
// Update backend CORS
cors.New(cors.Options{
    AllowedOrigins: []string{"https://your-frontend-url.vercel.app"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"*"},
})
```

**2. Database Connection**
- Check DATABASE_URL format
- Ensure database is running
- Verify network access

**3. Build Failures**
- Check Go version compatibility
- Verify all dependencies
- Review build logs

---

## 📈 **Performance Optimization**

### Frontend
- Enable Vercel's edge caching
- Optimize images with Next.js Image component
- Use dynamic imports for code splitting

### Backend
- Enable Railway's auto-scaling
- Use connection pooling for database
- Implement proper caching headers

---

## 🔒 **Security Considerations**

### Environment Variables
- Never commit secrets to GitHub
- Use Railway/Vercel environment variables
- Rotate JWT secrets regularly

### CORS Configuration
- Only allow your frontend domain
- Use HTTPS in production
- Implement proper authentication

---

## 🎉 **Success Checklist**

- [ ] Backend deployed and accessible
- [ ] Frontend deployed and working
- [ ] Database connected and migrations run
- [ ] Authentication working
- [ ] CORS configured properly
- [ ] Custom domain set up (optional)
- [ ] Health checks passing
- [ ] Performance optimized
- [ ] Security measures in place

---

**Your project is now live and ready for your resume! 🚀**

*Remember to update the URLs in this guide with your actual deployment URLs.*
