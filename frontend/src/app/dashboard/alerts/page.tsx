'use client';

import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import DashboardLayout from '@/components/DashboardLayout';
import { GET_SECURITY_ALERTS } from '@/lib/graphql/queries';
import { RESOLVE_SECURITY_ALERT } from '@/lib/graphql/mutations';
import {
  FaExclamationTriangle,
  FaExclamationCircle,
  FaInfoCircle,
  FaCheckCircle,
  FaMapMarkerAlt,
  FaCheck,
} from 'react-icons/fa';
import { format } from 'date-fns';
import type { SecurityAlert } from '@/types';

export default function SecurityAlertsPage() {
  const [includeResolved, setIncludeResolved] = useState(false);
  const [selectedSeverity, setSelectedSeverity] = useState<string>('');
  const [resolvingAlert, setResolvingAlert] = useState<string | null>(null);

  const { data, loading, refetch } = useQuery(GET_SECURITY_ALERTS, {
    variables: {
      includeResolved,
      severity: selectedSeverity || undefined,
    },
  });

  const [resolveAlert] = useMutation(RESOLVE_SECURITY_ALERT);

  const alerts = data?.securityAlerts?.alerts || [];
  const totalCount = data?.securityAlerts?.totalCount || 0;
  const unresolvedCount = data?.securityAlerts?.unresolvedCount || 0;

  const handleResolveAlert = async (alertId: string) => {
    if (!confirm('Mark this alert as resolved?')) return;

    setResolvingAlert(alertId);
    try {
      await resolveAlert({ variables: { alertId } });
      refetch();
    } catch (error) {
      console.error('Failed to resolve alert:', error);
      alert('Failed to resolve alert');
    } finally {
      setResolvingAlert(null);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'low':
        return 'bg-blue-100 text-blue-800 border-blue-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getSeverityIcon = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'critical':
      case 'high':
        return FaExclamationTriangle;
      case 'medium':
        return FaExclamationCircle;
      case 'low':
        return FaInfoCircle;
      default:
        return FaInfoCircle;
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 className="text-2xl font-bold text-gray-900">Security Alerts</h2>
              <p className="text-gray-600 mt-1">
                Monitor and respond to security events
              </p>
            </div>
            <div className="mt-4 sm:mt-0 flex items-center space-x-4">
              <select
                value={selectedSeverity}
                onChange={(e) => setSelectedSeverity(e.target.value)}
                className="border border-gray-300 text-black rounded-lg px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 "
              >
                <option value="">All Severities</option>
                <option value="critical">Critical</option>
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={includeResolved}
                  onChange={(e) => setIncludeResolved(e.target.checked)}
                  className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 h-4 w-4"
                />
                <span className="ml-2 text-sm text-gray-700">Show resolved</span>
              </label>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Alerts</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">{totalCount}</p>
              </div>
              <FaExclamationTriangle className="h-8 w-8 text-gray-400" />
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Unresolved</p>
                <p className="text-3xl font-bold text-red-600 mt-2">{unresolvedCount}</p>
              </div>
              <FaExclamationCircle className="h-8 w-8 text-red-400" />
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Resolved</p>
                <p className="text-3xl font-bold text-green-600 mt-2">
                  {totalCount - unresolvedCount}
                </p>
              </div>
              <FaCheckCircle className="h-8 w-8 text-green-400" />
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Resolution Rate</p>
                <p className="text-3xl font-bold text-blue-600 mt-2">
                  {totalCount > 0 ? Math.round(((totalCount - unresolvedCount) / totalCount) * 100) : 0}%
                </p>
              </div>
              <FaInfoCircle className="h-8 w-8 text-blue-400" />
            </div>
          </div>
        </div>

        {/* Alerts List */}
        {loading ? (
          <div className="bg-white rounded-xl shadow-md p-12 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading alerts...</p>
          </div>
        ) : alerts.length === 0 ? (
          <div className="bg-white rounded-xl shadow-md p-12 text-center">
            <FaCheckCircle className="h-16 w-16 text-green-400 mx-auto mb-4" />
            <p className="text-xl font-semibold text-gray-900 mb-2">All Clear!</p>
            <p className="text-gray-600">No security alerts found</p>
          </div>
        ) : (
          <div className="space-y-4">
            {alerts.map((alert: SecurityAlert) => {
              const SeverityIcon = getSeverityIcon(alert.severity);
              const isResolved = alert.isResolved;

              return (
                <div
                  key={alert.id}
                  className={`bg-white rounded-xl shadow-md p-6 border-l-4 ${
                    isResolved
                      ? 'border-green-500 opacity-75'
                      : alert.severity === 'critical'
                      ? 'border-red-500'
                      : alert.severity === 'high'
                      ? 'border-orange-500'
                      : alert.severity === 'medium'
                      ? 'border-yellow-500'
                      : 'border-blue-500'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start space-x-4 flex-1">
                      <div
                        className={`h-12 w-12 rounded-full flex items-center justify-center ${
                          isResolved
                            ? 'bg-green-100'
                            : alert.severity === 'critical'
                            ? 'bg-red-100'
                            : alert.severity === 'high'
                            ? 'bg-orange-100'
                            : alert.severity === 'medium'
                            ? 'bg-yellow-100'
                            : 'bg-blue-100'
                        }`}
                      >
                        {isResolved ? (
                          <FaCheckCircle className="h-6 w-6 text-green-600" />
                        ) : (
                          <SeverityIcon
                            className={`h-6 w-6 ${
                              alert.severity === 'critical'
                                ? 'text-red-600'
                                : alert.severity === 'high'
                                ? 'text-orange-600'
                                : alert.severity === 'medium'
                                ? 'text-yellow-600'
                                : 'text-blue-600'
                            }`}
                          />
                        )}
                      </div>

                      <div className="flex-1">
                        <div className="flex items-center space-x-3 mb-2">
                          <h3 className="text-lg font-semibold text-gray-900">
                            {alert.alertType.replace(/_/g, ' ').toUpperCase()}
                          </h3>
                          <span
                            className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium ${getSeverityColor(
                              alert.severity
                            )}`}
                          >
                            {alert.severity.toUpperCase()}
                          </span>
                          {isResolved && (
                            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                              <FaCheckCircle className="mr-1 h-3 w-3" />
                              Resolved
                            </span>
                          )}
                        </div>

                        <p className="text-gray-700 mb-4">{alert.description}</p>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                          <div className="space-y-2">
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">Created:</span>
                              <span className="text-gray-600">
                                {format(new Date(alert.createdAt), 'PPpp')}
                              </span>
                            </div>
                            {isResolved && alert.resolvedAt && (
                              <div className="flex items-start">
                                <span className="font-medium text-gray-700 w-32">Resolved:</span>
                                <span className="text-gray-600">
                                  {format(new Date(alert.resolvedAt), 'PPpp')}
                                </span>
                              </div>
                            )}
                          </div>

                          <div className="space-y-2">
                            {alert.ipAddress && (
                              <div className="flex items-start">
                                <span className="font-medium text-gray-700 w-32">IP Address:</span>
                                <span className="text-gray-600">{alert.ipAddress}</span>
                              </div>
                            )}
                            {(alert.locationCity || alert.locationCountry) && (
                              <div className="flex items-start">
                                <FaMapMarkerAlt className="mr-2 mt-1 h-4 w-4 text-gray-400" />
                                <span className="font-medium text-gray-700 w-28">Location:</span>
                                <span className="text-gray-600">
                                  {alert.locationCity && alert.locationCountry
                                    ? `${alert.locationCity}, ${alert.locationCountry}`
                                    : alert.locationCountry || 'Unknown'}
                                </span>
                              </div>
                            )}
                          </div>
                        </div>

                        {alert.metadata && (
                          <div className="mt-4 p-3 bg-gray-50 rounded-lg">
                            <p className="text-xs font-medium text-gray-700 mb-1">
                              Additional Details
                            </p>
                            <pre className="text-xs text-gray-600 overflow-x-auto">
                              {JSON.stringify(JSON.parse(alert.metadata), null, 2)}
                            </pre>
                          </div>
                        )}
                      </div>
                    </div>

                    {!isResolved && (
                      <button
                        onClick={() => handleResolveAlert(alert.id)}
                        disabled={resolvingAlert === alert.id}
                        className="ml-4 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center whitespace-nowrap"
                      >
                        <FaCheck className="mr-2" />
                        {resolvingAlert === alert.id ? 'Resolving...' : 'Resolve'}
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Info Box */}
        <div className="bg-yellow-50 border border-yellow-200 rounded-xl p-6">
          <div className="flex items-start">
            <FaExclamationTriangle className="h-6 w-6 text-yellow-600 mr-3 mt-1" />
            <div>
              <h4 className="font-semibold text-yellow-900 mb-2">About Security Alerts</h4>
              <p className="text-sm text-yellow-800 leading-relaxed">
                Security alerts are automatically generated when our system detects unusual or
                potentially malicious activity on your account. Review each alert carefully and
                mark them as resolved once you've investigated. If you see any alerts you don't
                recognize, change your password immediately and enable MFA.
              </p>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}