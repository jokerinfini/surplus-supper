# 🚀 Deployment Checklist - Surplus Supper

## ✅ **Step 1: PostgreSQL Database (COMPLETED)**
- [x] Added PostgreSQL service to Railway
- [x] Database is running and accessible

## 🔧 **Step 2: Backend Service Setup**

### **2.1 Add Backend Service**
- [ ] Go to Railway project dashboard
- [ ] Click **"New"** → **"GitHub Repo"**
- [ ] Select your `surplus-supper` repository
- [ ] Set **Root Directory** to `backend`
- [ ] Railway will auto-detect Go project

### **2.2 Configure Environment Variables**
Once backend service is added, go to **Variables** tab and add:

```env
DATABASE_URL=postgresql://username:password@host:port/database
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-here
CORS_ORIGIN=https://your-frontend-url.vercel.app
```

**How to get DATABASE_URL:**
1. Go to your **PostgreSQL** service in Railway
2. Click **"Connect"** tab
3. Copy the **PostgreSQL Connection URL**
4. Paste it as `DATABASE_URL`

**Generate JWT_SECRET:**
```bash
# Run this in terminal to generate a secure secret
openssl rand -base64 32
```

### **2.3 Deploy Backend**
- [ ] Railway will automatically build and deploy
- [ ] Wait for deployment to complete
- [ ] Copy the backend URL (e.g., `https://surplus-supper-backend.railway.app`)

### **2.4 Test Backend**
```bash
curl https://your-backend-url.railway.app/health
```
Expected response: `{"status": "healthy", "restaurant_count": 5}`

## 🎨 **Step 3: Frontend Deployment (Vercel)**

### **3.1 Create Vercel Account**
- [ ] Go to [vercel.com](https://vercel.com)
- [ ] Sign up with GitHub account
- [ ] Click **"New Project"**

### **3.2 Import Repository**
- [ ] Select your `surplus-supper` repository
- [ ] Set **Root Directory** to `frontend-next`
- [ ] Vercel will auto-detect Next.js

### **3.3 Configure Build Settings**
- [ ] **Framework Preset**: Next.js
- [ ] **Build Command**: `npm run build`
- [ ] **Output Directory**: `.next`
- [ ] **Install Command**: `npm install`

### **3.4 Add Environment Variable**
- [ ] Add environment variable:
  ```
  NEXT_PUBLIC_API_URL=https://your-backend-url.railway.app
  ```
- [ ] Replace `your-backend-url.railway.app` with your actual Railway backend URL

### **3.5 Deploy Frontend**
- [ ] Click **"Deploy"**
- [ ] Wait for deployment to complete
- [ ] Copy the frontend URL (e.g., `https://surplus-supper.vercel.app`)

## 🔄 **Step 4: Update CORS Configuration**

### **4.1 Update Backend CORS**
- [ ] Go back to Railway backend service
- [ ] Update `CORS_ORIGIN` environment variable:
  ```
  CORS_ORIGIN=https://your-frontend-url.vercel.app
  ```
- [ ] Replace with your actual Vercel frontend URL

### **4.2 Redeploy Backend**
- [ ] Railway will automatically redeploy with new environment variable

## 🧪 **Step 5: Testing**

### **5.1 Test Backend API**
```bash
# Health check
curl https://your-backend-url.railway.app/health

# Test restaurants endpoint
curl https://your-backend-url.railway.app/api/restaurants
```

### **5.2 Test Frontend**
- [ ] Visit your Vercel frontend URL
- [ ] Test user registration/login
- [ ] Test restaurant search
- [ ] Test location-based features

### **5.3 Test Authentication**
- [ ] Register a new user
- [ ] Login with credentials
- [ ] Verify JWT token is stored
- [ ] Test protected routes

## 🌐 **Step 6: Custom Domain (Optional)**

### **6.1 Vercel Custom Domain**
- [ ] Go to Vercel project settings
- [ ] Click **"Domains"**
- [ ] Add your custom domain
- [ ] Configure DNS records

### **6.2 Railway Custom Domain**
- [ ] Go to Railway project settings
- [ ] Click **"Domains"**
- [ ] Add custom domain
- [ ] Configure DNS records

## 📊 **Step 7: Monitoring**

### **7.1 Health Checks**
- [ ] Backend health: `https://your-backend-url/health`
- [ ] Frontend: Vercel provides built-in monitoring

### **7.2 Logs**
- [ ] Railway logs: Built-in logging dashboard
- [ ] Vercel logs: Function logs and analytics

## 🎯 **Step 8: Resume Links**

Once everything is working, update these URLs:

```
🌐 Live Demo: https://your-frontend-url.vercel.app
🔧 Backend API: https://your-backend-url.railway.app
📚 GitHub: https://github.com/yourusername/surplus-supper
```

## 🚨 **Troubleshooting**

### **Common Issues:**

**1. CORS Errors**
- Check that `CORS_ORIGIN` is set correctly
- Ensure frontend URL is exact (including https://)

**2. Database Connection**
- Verify `DATABASE_URL` format
- Check Railway PostgreSQL service is running

**3. Build Failures**
- Check Go version compatibility
- Verify all dependencies in go.mod

**4. Frontend Build Issues**
- Check Node.js version
- Verify all npm dependencies

## 🎉 **Success Indicators**

- [ ] Backend health check returns success
- [ ] Frontend loads without errors
- [ ] User registration/login works
- [ ] Restaurant search works
- [ ] Location-based features work
- [ ] No CORS errors in browser console
- [ ] All environment variables are set correctly

---

**Your project will be live and ready for your resume! 🚀**
