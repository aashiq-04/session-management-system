# Frontend Setup Guide - Session Management System

## Complete File Structure

Ensure your frontend directory has this exact structure:

```
frontend/
├── src/
│   ├── app/
│   │   ├── dashboard/
│   │   │   ├── alerts/
│   │   │   │   └── page.tsx
│   │   │   ├── devices/
│   │   │   │   └── page.tsx
│   │   │   ├── logs/
│   │   │   │   └── page.tsx
│   │   │   ├── sessions/
│   │   │   │   └── page.tsx
│   │   │   ├── settings/
│   │   │   │   └── page.tsx
│   │   │   └── page.tsx
│   │   ├── login/
│   │   │   └── page.tsx
│   │   ├── register/
│   │   │   └── page.tsx
│   │   ├── globals.css
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── components/
│   │   ├── ApolloWrapper.tsx
│   │   └── DashboardLayout.tsx
│   ├── lib/
│   │   ├── graphql/
│   │   │   ├── mutations.ts
│   │   │   └── queries.ts
│   │   ├── apollo-client.ts
│   │   ├── auth.ts
│   │   └── device.ts
│   └── types/
│       └── index.ts
├── public/
├── .env.local
├── next.config.js
├── package.json
├── postcss.config.js
├── tailwind.config.ts
└── tsconfig.json
```

## Installation Steps

### 1. Navigate to Frontend Directory

```bash
cd frontend
```

### 2. Install Dependencies

```bash
npm install
```

This will install all required packages:
- Next.js 14
- React 18
- TypeScript
- Tailwind CSS
- Apollo Client (GraphQL)
- React Icons
- Date-fns
- FingerprintJS
- Recharts (for charts)

### 3. Create Environment File

Create `.env.local` file in the frontend directory:

```bash
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/graphql" > .env.local
```

### 4. Verify All Files Are Created

Check that all the following files exist:

**Core App Files:**
- [ ] src/app/page.tsx
- [ ] src/app/layout.tsx
- [ ] src/app/globals.css
- [ ] src/app/login/page.tsx
- [ ] src/app/register/page.tsx

**Dashboard Files:**
- [ ] src/app/dashboard/page.tsx
- [ ] src/app/dashboard/sessions/page.tsx
- [ ] src/app/dashboard/devices/page.tsx
- [ ] src/app/dashboard/alerts/page.tsx
- [ ] src/app/dashboard/logs/page.tsx
- [ ] src/app/dashboard/settings/page.tsx

**Components:**
- [ ] src/components/ApolloWrapper.tsx
- [ ] src/components/DashboardLayout.tsx

**Library Files:**
- [ ] src/lib/apollo-client.ts
- [ ] src/lib/auth.ts
- [ ] src/lib/device.ts
- [ ] src/lib/graphql/mutations.ts
- [ ] src/lib/graphql/queries.ts

**Types:**
- [ ] src/types/index.ts

**Configuration Files:**
- [ ] package.json
- [ ] tsconfig.json
- [ ] next.config.js
- [ ] tailwind.config.ts
- [ ] postcss.config.js
- [ ] .env.local

### 5. Run Development Server

```bash
npm run dev
```

The application should now be running at: **http://localhost:3000**

## Testing the Application

### 1. Ensure Backend is Running

Make sure all backend services are running:
- Auth Service: localhost:50051
- Session Service: localhost:50052
- Audit Service: localhost:50053
- GraphQL Gateway: localhost:8080

### 2. Test User Flow

1. **Register a new account:**
   - Navigate to http://localhost:3000
   - Should redirect to login page
   - Click "Create one here"
   - Fill in registration form
   - Should automatically log you in and redirect to dashboard

2. **Test Login:**
   - Logout from dashboard
   - Navigate to login page
   - Enter credentials
   - Should redirect to dashboard

3. **Test Dashboard Features:**
   - **Dashboard Overview:** View stats, recent activity, security alerts
   - **Sessions:** See active sessions, revoke sessions
   - **Devices:** View trusted devices, trust new devices
   - **Security Alerts:** View and resolve security alerts
   - **Audit Logs:** View complete activity history with filters
   - **Settings:** Enable MFA (two-factor authentication)

## Common Issues and Solutions

### Issue 1: "Cannot find module '@/lib/...'"

**Solution:** Ensure tsconfig.json has the correct path mapping:
```json
{
  "compilerOptions": {
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

### Issue 2: GraphQL Connection Error

**Solution:** 
1. Verify GraphQL Gateway is running on port 8080
2. Check .env.local has correct API URL
3. Ensure CORS is enabled in gateway

### Issue 3: Device Fingerprinting Not Working

**Solution:**
- Ensure @fingerprintjs/fingerprintjs is installed
- Check browser console for errors
- Try in incognito mode to reset fingerprint

### Issue 4: Tailwind Styles Not Loading

**Solution:**
1. Verify tailwind.config.ts content paths are correct
2. Restart dev server: `npm run dev`
3. Clear .next folder: `rm -rf .next`

## Production Build

To create a production build:

```bash
npm run build
npm start
```

## Environment Variables for Production

For production deployment, update `.env.local`:

```env
NEXT_PUBLIC_API_URL=https://your-api-domain.com/graphql
```

## Features Implemented

✅ User Authentication (Login/Register)
✅ JWT Token Management
✅ Device Fingerprinting
✅ Session Management
✅ Multi-Device Tracking
✅ Security Alerts
✅ Audit Logs with Filtering
✅ MFA Setup (TOTP)
✅ Responsive Design
✅ Real-time Stats
✅ Activity Charts
✅ Session Revocation
✅ Device Trust Management

## Next Steps

1. **Customize Styling:** Modify Tailwind configuration and colors
2. **Add Charts:** Implement more detailed charts with Recharts
3. **Real-time Updates:** Add WebSocket support for live updates
4. **Mobile App:** Build React Native app using same GraphQL API
5. **Advanced Analytics:** Add more detailed compliance reports

## Support

For issues or questions:
- Check backend logs for API errors
- Verify all services are running
- Review browser console for frontend errors
- Check network tab for failed API calls