import FingerprintJS from '@fingerprintjs/fingerprintjs';

export interface DeviceInfo {
  deviceFingerprint: string;
  deviceName: string;
  deviceType: string;
  os: string;
  browser: string;
  ipAddress: string;
  userAgent: string;
  locationCountry?: string;
  locationCity?: string;
  latitude?: number;
  longitude?: number;
}

// Initialize FingerprintJS
const fpPromise = FingerprintJS.load();

// Get device fingerprint
export const getDeviceFingerprint = async (): Promise<string> => {
  const fp = await fpPromise;
  const result = await fp.get();
  return result.visitorId;
};

// Detect device type
const getDeviceType = (): string => {
  const ua = navigator.userAgent;
  if (/(tablet|ipad|playbook|silk)|(android(?!.*mobi))/i.test(ua)) {
    return 'tablet';
  }
  if (/Mobile|Android|iP(hone|od)|IEMobile|BlackBerry|Kindle|Silk-Accelerated|(hpw|web)OS|Opera M(obi|ini)/.test(ua)) {
    return 'mobile';
  }
  return 'desktop';
};

// Detect OS
const getOS = (): string => {
  const ua = navigator.userAgent;
  if (ua.indexOf('Win') !== -1) return 'Windows';
  if (ua.indexOf('Mac') !== -1) return 'macOS';
  if (ua.indexOf('Linux') !== -1) return 'Linux';
  if (ua.indexOf('Android') !== -1) return 'Android';
  if (ua.indexOf('iOS') !== -1) return 'iOS';
  return 'Unknown';
};

// Detect Browser
const getBrowser = (): string => {
  const ua = navigator.userAgent;
  if (ua.indexOf('Firefox') !== -1) return 'Firefox';
  if (ua.indexOf('Chrome') !== -1) return 'Chrome';
  if (ua.indexOf('Safari') !== -1) return 'Safari';
  if (ua.indexOf('Edge') !== -1) return 'Edge';
  if (ua.indexOf('Opera') !== -1 || ua.indexOf('OPR') !== -1) return 'Opera';
  return 'Unknown';
};

// Get device name
const getDeviceName = (): string => {
  const os = getOS();
  const browser = getBrowser();
  return `${os} - ${browser}`;
};

// Get IP address (we'll use a placeholder - backend will get real IP)
const getIPAddress = (): string => {
  return '0.0.0.0'; // Backend will capture the real IP
};

// Get complete device info
export const getDeviceInfo = async (): Promise<DeviceInfo> => {
  const fingerprint = await getDeviceFingerprint();
  
  return {
    deviceFingerprint: fingerprint,
    deviceName: getDeviceName(),
    deviceType: getDeviceType(),
    os: getOS(),
    browser: getBrowser(),
    ipAddress: getIPAddress(),
    userAgent: navigator.userAgent,
  };
};

// Format device info for display
export const formatDeviceInfo = (device: any): string => {
  if (device.deviceName) return device.deviceName;
  if (device.os && device.browser) return `${device.os} - ${device.browser}`;
  return 'Unknown Device';
};