'use client';

import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import DashboardLayout from '@/components/DashboardLayout';
import { GET_ME } from '@/lib/graphql/queries';
import { ENABLE_MFA, VERIFY_MFA } from '@/lib/graphql/mutations';
import {
  FaShieldAlt,
  FaQrcode,
  FaKey,
  FaCheckCircle,
  FaExclamationTriangle,
  FaCopy,
  FaUser,
  FaEnvelope,
  FaCalendar,
} from 'react-icons/fa';
import { format } from 'date-fns';

export default function SettingsPage() {
  const [mfaStep, setMfaStep] = useState<'initial' | 'setup' | 'verify'>('initial');
  const [mfaData, setMfaData] = useState<any>(null);
  const [verificationCode, setVerificationCode] = useState('');
  const [error, setError] = useState('');
  const [copiedCode, setCopiedCode] = useState<string | null>(null);

  const { data, refetch } = useQuery(GET_ME);
  const [enableMFA] = useMutation(ENABLE_MFA);
  const [verifyMFA] = useMutation(VERIFY_MFA);

  const user = data?.me;

  const handleEnableMFA = async () => {
    setError('');
    try {
      const { data } = await enableMFA();
      if (data.enableMFA.success) {
        setMfaData(data.enableMFA);
        setMfaStep('setup');
      } else {
        setError(data.enableMFA.message);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to enable MFA');
    }
  };

  const handleVerifyMFA = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      const { data } = await verifyMFA({ variables: { code: verificationCode } });
      if (data.verifyMFA.success) {
        setMfaStep('verify');
        refetch();
        setTimeout(() => {
          setMfaStep('initial');
          setMfaData(null);
          setVerificationCode('');
        }, 3000);
      } else {
        setError(data.verifyMFA.message);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to verify MFA code');
    }
  };

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    setCopiedCode(label);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* User Profile Section */}
        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center mb-6">
            <FaUser className="h-5 w-5 text-gray-400 mr-2" />
            <h2 className="text-2xl font-bold text-gray-900">User Profile</h2>
          </div>

          <div className="space-y-4">
            <div className="flex items-start">
              <div className="h-16 w-16 rounded-full bg-blue-600 flex items-center justify-center text-white text-2xl font-bold">
                {user?.fullName?.charAt(0).toUpperCase()}
              </div>
              <div className="ml-4 flex-1">
                <h3 className="text-xl font-semibold text-gray-900">{user?.fullName}</h3>
                <div className="mt-2 space-y-2 text-sm text-gray-600">
                  <div className="flex items-center">
                    <FaEnvelope className="mr-2 h-4 w-4 text-gray-400" />
                    <span>{user?.email}</span>
                  </div>
                  <div className="flex items-center">
                    <FaCalendar className="mr-2 h-4 w-4 text-gray-400" />
                    <span>
                      Member since {user?.createdAt && format(new Date(user.createdAt), 'MMMM yyyy')}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <FaShieldAlt className="mr-2 h-4 w-4 text-gray-400" />
                    <span>
                      Account Status:{' '}
                      <span className={user?.isActive ? 'text-green-600 font-medium' : 'text-red-600 font-medium'}>
                        {user?.isActive ? 'Active' : 'Inactive'}
                      </span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* MFA Section */}
        <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
          <div className="flex items-center mb-6">
            <FaShieldAlt className="h-5 w-5 text-gray-400 mr-2" />
            <h2 className="text-2xl font-bold text-gray-900">Two-Factor Authentication</h2>
          </div>

          {user?.mfaEnabled ? (
            <div className="bg-green-50 border border-green-200 rounded-lg p-6">
              <div className="flex items-start">
                <FaCheckCircle className="h-6 w-6 text-green-600 mr-3 mt-1" />
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-green-900 mb-2">
                    MFA is Enabled
                  </h3>
                  <p className="text-green-800">
                    Your account is protected with two-factor authentication. You'll need to
                    provide a verification code from your authenticator app when you sign in.
                  </p>
                </div>
              </div>
            </div>
          ) : (
            <>
              {mfaStep === 'initial' && (
                <div className="space-y-6">
                  <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-6">
                    <div className="flex items-start">
                      <FaExclamationTriangle className="h-6 w-6 text-yellow-600 mr-3 mt-1" />
                      <div className="flex-1">
                        <h3 className="text-lg font-semibold text-yellow-900 mb-2">
                          MFA Not Enabled
                        </h3>
                        <p className="text-yellow-800 mb-4">
                          Two-factor authentication adds an extra layer of security to your
                          account. We highly recommend enabling it to protect against
                          unauthorized access.
                        </p>
                        <button
                          onClick={handleEnableMFA}
                          className="px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition flex items-center"
                        >
                          <FaShieldAlt className="mr-2" />
                          Enable Two-Factor Authentication
                        </button>
                      </div>
                    </div>
                  </div>

                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
                    <h4 className="font-semibold text-blue-900 mb-3">How it works:</h4>
                    <ol className="list-decimal list-inside space-y-2 text-sm text-blue-800">
                      <li>Download an authenticator app (Google Authenticator, Authy, etc.)</li>
                      <li>Scan the QR code we provide with your app</li>
                      <li>Enter the 6-digit code from your app to verify</li>
                      <li>Save your backup codes in a secure location</li>
                      <li>Use the code from your app every time you sign in</li>
                    </ol>
                  </div>
                </div>
              )}

              {mfaStep === 'setup' && mfaData && (
                <div className="space-y-6">
                  {error && (
                    <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
                      {error}
                    </div>
                  )}

                  {/* QR Code */}
                  <div className="bg-gray-50 rounded-lg p-6 text-center">
                    <div className="flex items-center justify-center mb-4">
                      <FaQrcode className="h-6 w-6 text-gray-600 mr-2" />
                      <h3 className="text-lg font-semibold text-gray-900">
                        Scan QR Code
                      </h3>
                    </div>
                    <div className="bg-white p-4 rounded-lg inline-block border-4 border-gray-200">
                      <img
                        src={mfaData.qrCodeUrl}
                        alt="MFA QR Code"
                        className="w-64 h-64"
                      />
                    </div>
                    <p className="text-sm text-gray-600 mt-4">
                      Scan this QR code with your authenticator app
                    </p>
                  </div>

                  {/* Manual Entry */}
                  <div className="bg-gray-50 rounded-lg p-6">
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center">
                        <FaKey className="h-5 w-5 text-gray-600 mr-2" />
                        <h4 className="font-semibold text-gray-900">Manual Entry</h4>
                      </div>
                      <button
                        onClick={() => copyToClipboard(mfaData.secret, 'secret')}
                        className="text-sm text-blue-600 hover:text-blue-800 flex items-center"
                      >
                        <FaCopy className="mr-1" />
                        {copiedCode === 'secret' ? 'Copied!' : 'Copy'}
                      </button>
                    </div>
                    <p className="text-xs text-gray-600 mb-2">
                      If you can't scan the QR code, enter this key manually:
                    </p>
                    <code className="block bg-white p-3 rounded border border-gray-200 text-sm font-mono break-all">
                      {mfaData.secret}
                    </code>
                  </div>

                  {/* Backup Codes */}
                  <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-6">
                    <h4 className="font-semibold text-yellow-900 mb-3 flex items-center">
                      <FaExclamationTriangle className="mr-2" />
                      Backup Codes
                    </h4>
                    <p className="text-sm text-yellow-800 mb-3">
                      Save these backup codes in a secure location. You can use them to access
                      your account if you lose your device.
                    </p>
                    <div className="grid grid-cols-2 gap-2 bg-white p-4 rounded border border-yellow-300">
                      {mfaData.backupCodes?.map((code: string, index: number) => (
                        <div
                          key={index}
                          className="flex items-center justify-between bg-gray-50 p-2 rounded"
                        >
                          <code className="text-sm font-mono">{code}</code>
                          <button
                            onClick={() => copyToClipboard(code, code)}
                            className="text-xs text-blue-600 hover:text-blue-800"
                          >
                            {copiedCode === code ? '✓' : <FaCopy />}
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Verification */}
                  <form onSubmit={handleVerifyMFA} className="bg-gray-50 rounded-lg p-6">
                    <h4 className="font-semibold text-gray-900 mb-4">
                      Verify Setup
                    </h4>
                    <p className="text-sm text-gray-600 mb-4">
                      Enter the 6-digit code from your authenticator app to complete setup:
                    </p>
                    <div className="flex space-x-4">
                      <input
                        type="text"
                        value={verificationCode}
                        onChange={(e) => setVerificationCode(e.target.value)}
                        placeholder="000000"
                        maxLength={6}
                        className="flex-1 px-4 py-3 border border-gray-300 rounded-lg text-center text-2xl tracking-widest focus:outline-none focus:ring-2 focus:ring-blue-500"
                        required
                      />
                      <button
                        type="submit"
                        className="px-8 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition"
                      >
                        Verify
                      </button>
                    </div>
                  </form>
                </div>
              )}

              {mfaStep === 'verify' && (
                <div className="bg-green-50 border border-green-200 rounded-lg p-6">
                  <div className="flex items-start">
                    <FaCheckCircle className="h-8 w-8 text-green-600 mr-4 mt-1" />
                    <div>
                      <h3 className="text-xl font-semibold text-green-900 mb-2">
                        MFA Successfully Enabled!
                      </h3>
                      <p className="text-green-800">
                        Your account is now protected with two-factor authentication. You'll
                        need to provide a code from your authenticator app when you sign in.
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
}