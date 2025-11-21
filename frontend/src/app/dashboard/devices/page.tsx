'use client';

import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import DashboardLayout from '@/components/DashboardLayout';
import { GET_DEVICES } from '@/lib/graphql/queries';
import { TRUST_DEVICE } from '@/lib/graphql/mutations';
import {
  FaDesktop,
  FaMobileAlt,
  FaTabletAlt,
  FaShieldAlt,
  FaCheckCircle,
  FaClock,
} from 'react-icons/fa';
import { format } from 'date-fns';
import type { Device } from '@/types';

export default function DevicesPage() {
  const [trustingDevice, setTrustingDevice] = useState<string | null>(null);

  const { data, loading, refetch } = useQuery(GET_DEVICES);
  const [trustDeviceMutation] = useMutation(TRUST_DEVICE);

  const devices = data?.devices?.devices || [];
  const totalCount = data?.devices?.totalCount || 0;
  const trustedCount = data?.devices?.trustedCount || 0;

  const handleTrustDevice = async (deviceId: string) => {
    if (!confirm('Are you sure you want to mark this device as trusted? Trusted devices have additional privileges.')) {
      return;
    }

    setTrustingDevice(deviceId);
    try {
      await trustDeviceMutation({ variables: { deviceId } });
      refetch();
    } catch (error) {
      console.error('Failed to trust device:', error);
      alert('Failed to trust device');
    } finally {
      setTrustingDevice(null);
    }
  };

  const getDeviceIcon = (deviceType: string) => {
    switch (deviceType?.toLowerCase()) {
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
              <h2 className="text-2xl font-bold text-gray-900">Trusted Devices</h2>
              <p className="text-gray-600 mt-1">
                Manage devices that have accessed your account
              </p>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Devices</p>
                <p className="text-3xl font-bold text-gray-900 mt-2">{totalCount}</p>
              </div>
              <div className="h-12 w-12 bg-blue-100 rounded-full flex items-center justify-center">
                <FaDesktop className="h-6 w-6 text-blue-600" />
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Trusted Devices</p>
                <p className="text-3xl font-bold text-green-600 mt-2">{trustedCount}</p>
              </div>
              <div className="h-12 w-12 bg-green-100 rounded-full flex items-center justify-center">
                <FaShieldAlt className="h-6 w-6 text-green-600" />
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Untrusted Devices</p>
                <p className="text-3xl font-bold text-yellow-600 mt-2">
                  {totalCount - trustedCount}
                </p>
              </div>
              <div className="h-12 w-12 bg-yellow-100 rounded-full flex items-center justify-center">
                <FaClock className="h-6 w-6 text-yellow-600" />
              </div>
            </div>
          </div>
        </div>

        {/* Devices List */}
        {loading ? (
          <div className="bg-white rounded-xl shadow-md p-12 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading devices...</p>
          </div>
        ) : devices.length === 0 ? (
          <div className="bg-white rounded-xl shadow-md p-12 text-center">
            <FaDesktop className="h-16 w-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-600">No devices found</p>
          </div>
        ) : (
          <div className="space-y-4">
            {devices.map((device: Device) => {
              const DeviceIcon = getDeviceIcon(device.deviceType);
              const isTrusted = device.isTrusted;

              return (
                <div
                  key={device.id}
                  className={`bg-white rounded-xl shadow-md p-6 border ${
                    isTrusted ? 'border-green-200 bg-green-50' : 'border-gray-200'
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start space-x-4 flex-1">
                      <div
                        className={`h-14 w-14 rounded-full flex items-center justify-center ${
                          isTrusted ? 'bg-green-100' : 'bg-gray-100'
                        }`}
                      >
                        <DeviceIcon
                          className={`h-7 w-7 ${
                            isTrusted ? 'text-green-600' : 'text-gray-500'
                          }`}
                        />
                      </div>

                      <div className="flex-1">
                        <div className="flex items-center space-x-3 mb-3">
                          <h3 className="text-lg font-semibold text-gray-900">
                            {device.deviceName || 'Unknown Device'}
                          </h3>
                          {isTrusted ? (
                            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                              <FaShieldAlt className="mr-1 h-3 w-3" />
                              Trusted
                            </span>
                          ) : (
                            <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                              <FaClock className="mr-1 h-3 w-3" />
                              Untrusted
                            </span>
                          )}
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                          <div className="space-y-2">
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">Device Type:</span>
                              <span className="text-gray-600 capitalize">
                                {device.deviceType || 'Unknown'}
                              </span>
                            </div>
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">Operating System:</span>
                              <span className="text-gray-600">{device.os || 'Unknown'}</span>
                            </div>
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">Browser:</span>
                              <span className="text-gray-600">{device.browser || 'Unknown'}</span>
                            </div>
                          </div>

                          <div className="space-y-2">
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">First Seen:</span>
                              <span className="text-gray-600">
                                {format(new Date(device.firstSeenAt), 'PPp')}
                              </span>
                            </div>
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">Last Seen:</span>
                              <span className="text-gray-600">
                                {format(new Date(device.lastSeenAt), 'PPp')}
                              </span>
                            </div>
                            <div className="flex items-start">
                              <span className="font-medium text-gray-700 w-32">Sessions:</span>
                              <span className="text-gray-600">
                                {device.sessionCount || 0} total
                              </span>
                            </div>
                          </div>
                        </div>

                        <div className="mt-4 p-3 bg-gray-50 rounded-lg">
                          <p className="text-xs font-medium text-gray-700 mb-1">
                            Device Fingerprint
                          </p>
                          <p className="text-xs text-gray-500 font-mono break-all">
                            {device.deviceFingerprint}
                          </p>
                        </div>
                      </div>
                    </div>

                    {!isTrusted && (
                      <button
                        onClick={() => handleTrustDevice(device.id)}
                        disabled={trustingDevice === device.id}
                        className="ml-4 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center whitespace-nowrap"
                      >
                        <FaShieldAlt className="mr-2" />
                        {trustingDevice === device.id ? 'Trusting...' : 'Trust Device'}
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Info Box */}
        <div className="bg-blue-50 border border-blue-200 rounded-xl p-6">
          <div className="flex items-start">
            <FaShieldAlt className="h-6 w-6 text-blue-600 mr-3 mt-1" />
            <div>
              <h4 className="font-semibold text-blue-900 mb-2">About Trusted Devices</h4>
              <p className="text-sm text-blue-800 leading-relaxed">
                Marking a device as trusted helps us recognize your regular devices and provides
                enhanced security. Trusted devices may bypass certain security checks and will be
                flagged if unusual activity is detected. You can revoke trust at any time.
              </p>
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}