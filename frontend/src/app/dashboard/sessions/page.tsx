'use client';

import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import DashboardLayout from '@/components/DashboardLayout';
import { GET_SESSIONS } from '@/lib/graphql/queries';
import { REVOKE_SESSION, REVOKE_ALL_SESSIONS } from '@/lib/graphql/mutations';
import {
  FaDesktop,
  FaMobileAlt,
  FaTabletAlt,
  FaMapMarkerAlt,
  FaCheckCircle,
  FaTimesCircle,
  FaTrash,
  FaSignOutAlt,
} from 'react-icons/fa';
import { format } from 'date-fns';
import type { Session } from '@/types';

export default function SessionsPage() {
  const [includeInactive, setIncludeInactive] = useState(false);
  const [revoking, setRevoking] = useState<string | null>(null);

  const { data, loading, refetch } = useQuery(GET_SESSIONS, {
    variables: { includeInactive },
  });

  const [revokeSession] = useMutation(REVOKE_SESSION);
  const [revokeAllSessions] = useMutation(REVOKE_ALL_SESSIONS);

  const sessions = data?.sessions?.sessions || [];
  const activeCount = data?.sessions?.activeCount || 0;
  const totalCount = data?.sessions?.totalCount || 0;

  const handleRevokeSession = async (sessionId: string) => {
    if (!confirm('Are you sure you want to revoke this session?')) return;

    setRevoking(sessionId);
    try {
      await revokeSession({ variables: { sessionId } });
      refetch();
    } catch (error) {
      console.error('Failed to revoke session:', error);
      alert('Failed to revoke session');
    } finally {
      setRevoking(null);
    }
  };

  const handleRevokeAll = async () => {
    if (!confirm('Are you sure you want to logout from all other devices? This will end all active sessions except your current one.')) {
      return;
    }

    try {
      await revokeAllSessions({ variables: { exceptCurrent: true } });
      refetch();
      alert('All other sessions have been revoked');
    } catch (error) {
      console.error('Failed to revoke all sessions:', error);
      alert('Failed to revoke all sessions');
    }
  };

  const getDeviceIcon = (deviceType: string) => {
    switch (deviceType.toLowerCase()) {
      case 'mobile':
        return FaMobileAlt;
      case 'tablet':
        return FaTabletAlt;
      default:
        return FaDesktop;
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">Active Sessions</h2>
              <p className="text-gray-600 mt-1">
                Manage your login sessions across all devices
              </p>
            </div>
            <div className="mt-4 sm:mt-0 flex items-center space-x-4">
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={includeInactive}
                  onChange={(e) => setIncludeInactive(e.target.checked)}
                  className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 h-4 w-4"
                />
                <span className="ml-2 text-sm text-gray-700">Show inactive</span>
              </label>
              {activeCount > 1 && (
                <button
                  onClick={handleRevokeAll}
                  className="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 transition flex items-center"
                >
                  <FaSignOutAlt className="mr-2" />
                  Logout All Devices
                </button>
              )}
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <p className="text-sm font-medium text-gray-600">Active Sessions</p>
            <p className="text-3xl font-bold text-green-600 mt-2">{activeCount}</p>
          </div>
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <p className="text-sm font-medium text-gray-600">Total Sessions</p>
            <p className="text-3xl font-bold text-gray-900 mt-2">{totalCount}</p>
          </div>
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <p className="text-sm font-medium text-gray-600">Inactive Sessions</p>
            <p className="text-3xl font-bold text-gray-400 mt-2">{totalCount - activeCount}</p>
          </div>
        </div>

        {/* Sessions List */}
        {loading ? (
          <div className="bg-white rounded-xl shadow-md p-12 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading sessions...</p>
          </div>
        ) : sessions.length === 0 ? (
          <div className="bg-white rounded-xl shadow-md p-12 text-center">
            <FaDesktop className="h-16 w-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-600">No sessions found</p>
          </div>
        ) : (
          <div className="space-y-4">
            {sessions.map((session: Session) => {
              const DeviceIcon = getDeviceIcon(session.deviceType);
              const isActive = session.isActive;

              return (
                <div
                  key={session.id}
                  className={`bg-white rounded-xl shadow-md p-6 border ${
                    isActive ? 'border-green-200' : 'border-gray-200'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start space-x-4">
                      <div
                        className={`h-12 w-12 rounded-full flex items-center justify-center ${
                          isActive ? 'bg-green-100' : 'bg-gray-100'
                        }`}
                      >
                        <DeviceIcon
                          className={`h-6 w-6 ${
                            isActive ? 'text-green-600' : 'text-gray-400'
                          }`}
                        />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center space-x-3 mb-2">
                          <h3 className="text-lg font-semibold text-gray-900">
                            {session.deviceName}
                          </h3>
                          {isActive ? (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                              <FaCheckCircle className="mr-1 h-3 w-3" />
                              Active
                            </span>
                          ) : (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                              <FaTimesCircle className="mr-1 h-3 w-3" />
                              Inactive
                            </span>
                          )}
                          {session.isCurrent && (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                              Current
                            </span>
                          )}
                        </div>

                        <div className="space-y-2 text-sm text-gray-600">
                          <div className="flex items-center">
                            <span className="font-medium w-32">IP Address:</span>
                            <span>{session.ipAddress}</span>
                          </div>
                          {(session.locationCity || session.locationCountry) && (
                            <div className="flex items-center">
                              <FaMapMarkerAlt className="mr-2 h-4 w-4 text-gray-400" />
                              <span className="font-medium w-30">Location:</span>
                              <span className="ml-2">
                                {session.locationCity && session.locationCountry
                                  ? `${session.locationCity}, ${session.locationCountry}`
                                  : session.locationCountry || 'Unknown'}
                              </span>
                            </div>
                          )}
                          <div className="flex items-center">
                            <span className="font-medium w-32">First Seen:</span>
                            <span>{format(new Date(session.createdAt), 'PPpp')}</span>
                          </div>
                          <div className="flex items-center">
                            <span className="font-medium w-32">Last Active:</span>
                            <span>{format(new Date(session.lastSeenAt), 'PPpp')}</span>
                          </div>
                          {session.userAgent && (
                            <div className="flex items-start">
                              <span className="font-medium w-32">User Agent:</span>
                              <span className="flex-1 text-xs">{session.userAgent}</span>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    {isActive && !session.isCurrent && (
                      <button
                        onClick={() => handleRevokeSession(session.id)}
                        disabled={revoking === session.id}
                        className="ml-4 px-4 py-2 bg-red-100 text-red-700 text-sm font-medium rounded-lg hover:bg-red-200 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
                      >
                        <FaTrash className="mr-2" />
                        {revoking === session.id ? 'Revoking...' : 'Revoke'}
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}