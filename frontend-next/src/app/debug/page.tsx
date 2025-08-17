'use client'

export default function DebugPage() {
  const rawApiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  const apiUrl = rawApiUrl.endsWith('/api') ? rawApiUrl.slice(0, -4) : rawApiUrl
  
  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Debug Information</h1>
      <div className="space-y-4">
        <div>
          <strong>Raw NEXT_PUBLIC_API_URL:</strong> {rawApiUrl}
        </div>
        <div>
          <strong>Processed API URL:</strong> {apiUrl}
        </div>
        <div>
          <strong>Full API URL for restaurants:</strong> {apiUrl}/api/restaurants
        </div>
        <div>
          <strong>Full API URL for auth:</strong> {apiUrl}/api/auth/login
        </div>
        <div>
          <strong>Environment:</strong> {process.env.NODE_ENV}
        </div>
      </div>
    </div>
  )
}
