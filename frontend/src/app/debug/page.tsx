'use client';

import { auth } from '@/lib/auth';
import { useQuery } from '@apollo/client';
import { GET_ME } from '@/lib/graphql/queries';

export default function DebugPage() {
  const { data, loading, error } = useQuery(GET_ME);

  const accessToken = auth.getAccessToken();
  const refreshToken = auth.getRefreshToken();
  const userId = auth.getUserId();

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Debug Information</h1>

        {/* Tokens */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Stored Tokens</h2>
          <div className="space-y-2 text-sm">
            <div>
              <strong>Access Token:</strong>
              <pre className="bg-gray-100 p-2 rounded mt-1 overflow-x-auto">
                {accessToken || 'Not found'}
              </pre>
            </div>
            <div>
              <strong>Refresh Token:</strong>
              <pre className="bg-gray-100 p-2 rounded mt-1 overflow-x-auto">
                {refreshToken || 'Not found'}
              </pre>
            </div>
            <div>
              <strong>User ID:</strong>
              <pre className="bg-gray-100 p-2 rounded mt-1">
                {userId || 'Not found'}
              </pre>
            </div>
          </div>
        </div>

        {/* GraphQL Query Result */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">GET_ME Query Result</h2>
          {loading && <p>Loading...</p>}
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded">
              <strong>Error:</strong>
              <pre className="mt-2 text-xs overflow-x-auto">
                {JSON.stringify(error, null, 2)}
              </pre>
            </div>
          )}
          {data && (
            <div className="bg-green-50 border border-green-200 text-green-700 p-4 rounded">
              <strong>Success!</strong>
              <pre className="mt-2 text-xs overflow-x-auto">
                {JSON.stringify(data, null, 2)}
              </pre>
            </div>
          )}
        </div>

        {/* API URL */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Configuration</h2>
          <div>
            <strong>API URL:</strong>
            <pre className="bg-gray-100 p-2 rounded mt-1">
              {process.env.NEXT_PUBLIC_API_URL || 'Not configured'}
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}