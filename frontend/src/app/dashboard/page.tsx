'use client';

import { useQuery } from '@apollo/client';
// import DashboardLayout from '@/app/dashboard/components/DashboardLayout';
import DashboardLayout from '@/components/DashboardLayout';
import {
  GET_SESSION_STATS,
  GET_SECURITY_ALERTS,
  GET_ACTIVITY_SUMMARY,
} from '@/lib/graphql/queries';
import {
  FaDesktop,
  FaMobileAlt,
  FaShieldAlt,
  FaExclamationTriangle,
  FaCheckCircle,
  FaClock,
  FaMapMarkerAlt,
} from 'react-icons/fa';
import { format } from 'date-fns';

export default function DashboardPage() {
  const { data: statsData, loading: statsLoading } = useQuery(GET_SESSION_STATS);
  const { data: alertsData, loading: alertsLoading } = useQuery(GET_SECURITY_ALERTS, {
    variables: { includeResolved: false },
  });
  const { data: activityData, loading: activityLoading } = useQuery(GET_ACTIVITY_SUMMARY, {
    variables: { days: 7 },
  });

  const stats = statsData?.sessionStats;
  const alerts = alertsData?.securityAlerts?.alerts || [];
  const activity = activityData?.activitySummary;

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Welcome Section */}
        <div className="bg-gradient-to-r from-blue-600 to-indigo-700 rounded-xl shadow-lg p-6 text-white">
          <h2 className="text-2xl font-bold mb-2">Welcome Back!</h2>
          <p className="text-blue-100">
            Here's an overview of your account security and activity
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {/* Active Sessions */}
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Active Sessions</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">
                  {statsLoading ? '...' : stats?.activeSessions || 0}
                </p>
                <p className="text-xs text-gray-500 mt-1">
                  of {stats?.totalSessions || 0} total
                </p>
              </div>
              <div className="h-12 w-12 bg-blue-100 rounded-full flex items-center justify-center">
                <FaDesktop className="h-6 w-6 text-blue-600" />
              </div>
            </div>
          </div>

          {/* Total Devices */}
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Devices</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">
                  {statsLoading ? '...' : stats?.totalDevices || 0}
                </p>
                <p className="text-xs text-gray-500 mt-1">
                  {stats?.trustedDevices || 0} trusted
                </p>
              </div>
              <div className="h-12 w-12 bg-green-100 rounded-full flex items-center justify-center">
                <FaMobileAlt className="h-6 w-6 text-green-600" />
              </div>
            </div>
          </div>

          {/* Security Alerts */}
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Security Alerts</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">
                  {alertsLoading ? '...' : alerts.length}
                </p>
                <p className="text-xs text-gray-500 mt-1">unresolved</p>
              </div>
              <div className="h-12 w-12 bg-yellow-100 rounded-full flex items-center justify-center">
                <FaExclamationTriangle className="h-6 w-6 text-yellow-600" />
              </div>
            </div>
          </div>

          {/* Login Activity */}
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Logins</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">
                  {activityLoading ? '...' : activity?.totalLogins || 0}
                </p>
                <p className="text-xs text-gray-500 mt-1">last 7 days</p>
              </div>
              <div className="h-12 w-12 bg-purple-100 rounded-full flex items-center justify-center">
                <FaCheckCircle className="h-6 w-6 text-purple-600" />
              </div>
            </div>
          </div>
        </div>

        {/* Last Login & Recent Locations */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Last Login */}
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center mb-4">
              <FaClock className="h-5 w-5 text-gray-400 mr-2" />
              <h3 className="text-lg font-semibold text-gray-900">Last Login</h3>
            </div>
            {statsLoading ? (
              <div className="animate-pulse space-y-2">
                <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                <div className="h-4 bg-gray-200 rounded w-1/2"></div>
              </div>
            ) : (
              <div className="space-y-2">
                <p className="text-2xl font-bold text-gray-900">
                  {stats?.lastLogin
                    ? format(new Date(stats.lastLogin), 'PPpp')
                    : 'No recent login'}
                </p>
                {stats?.lastLoginLocation && (
                  <div className="flex items-center text-gray-600">
                    <FaMapMarkerAlt className="h-4 w-4 mr-2" />
                    <span>{stats.lastLoginLocation}</span>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Recent Locations */}
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center mb-4">
              <FaMapMarkerAlt className="h-5 w-5 text-gray-400 mr-2" />
              <h3 className="text-lg font-semibold text-gray-900">Recent Locations</h3>
            </div>
            {statsLoading ? (
              <div className="animate-pulse space-y-2">
                <div className="h-4 bg-gray-200 rounded w-full"></div>
                <div className="h-4 bg-gray-200 rounded w-5/6"></div>
                <div className="h-4 bg-gray-200 rounded w-4/6"></div>
              </div>
            ) : (
              <div className="space-y-2">
                {stats?.recentLocations && stats.recentLocations.length > 0 ? (
                  stats.recentLocations.map((location: string, index: number) => (
                    <div
                      key={index}
                      className="flex items-center py-2 px-3 bg-gray-50 rounded-lg"
                    >
                      <FaMapMarkerAlt className="h-4 w-4 text-blue-600 mr-3" />
                      <span className="text-gray-700">{location}</span>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-500 text-sm">No recent locations</p>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Security Alerts */}
        {alerts.length > 0 && (
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center">
                <FaExclamationTriangle className="h-5 w-5 text-yellow-500 mr-2" />
                <h3 className="text-lg font-semibold text-gray-900">
                  Recent Security Alerts
                </h3>
              </div>
              <a
                href="/dashboard/alerts"
                className="text-sm text-blue-600 hover:text-blue-800 font-medium"
              >
                View All →
              </a>
            </div>
            <div className="space-y-3">
              {alerts.slice(0, 3).map((alert: any) => (
                <div
                  key={alert.id}
                  className="flex items-start p-4 bg-yellow-50 border border-yellow-200 rounded-lg"
                >
                  <FaShieldAlt className="h-5 w-5 text-yellow-600 mr-3 mt-0.5" />
                  <div className="flex-1">
                    <div className="flex items-center justify-between">
                      <p className="font-medium text-gray-900">{alert.alertType}</p>
                      <span className="text-xs text-gray-500">
                        {format(new Date(alert.createdAt), 'MMM dd, yyyy')}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">{alert.description}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Activity Chart Placeholder */}
        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Login Activity (Last 7 Days)</h3>
          {activityLoading ? (
            <div className="h-64 flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : activity?.dailyActivity && activity.dailyActivity.length > 0 ? (
            <div className="space-y-3">
              {activity.dailyActivity.map((day: any) => (
                <div key={day.date} className="flex items-center">
                  <span className="text-sm text-gray-600 w-24">
                    {format(new Date(day.date), 'MMM dd')}
                  </span>
                  <div className="flex-1 flex items-center space-x-2">
                    <div className="flex-1 bg-gray-200 rounded-full h-2 overflow-hidden">
                      <div
                        className="bg-blue-600 h-full rounded-full"
                        style={{
                          width: `${Math.min((day.loginCount / Math.max(...activity.dailyActivity.map((d: any) => d.loginCount))) * 100, 100)}%`,
                        }}
                      ></div>
                    </div>
                    <span className="text-sm font-medium text-gray-900 w-8 text-right">
                      {day.loginCount}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-center py-8">No activity data available</p>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}