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

// Initialize FingerprintJS with custom options for stability
const fpPromise = FingerprintJS.load();

// Get device fingerprint with fallback to localStorage
export const getDeviceFingerprint = async (): Promise<string> => {
  try {
    // Check if we already have a fingerprint stored
    const storedFingerprint = localStorage.getItem('deviceFingerprint');
    
    // If we have a stored fingerprint and it's recent (less than 30 days old)
    const storedTime = localStorage.getItem('deviceFingerprintTime');
    if (storedFingerprint && storedTime) {
      const daysSinceStored = (Date.now() - parseInt(storedTime)) / (1000 * 60 * 60 * 24);
      if (daysSinceStored < 30) {
        console.log('Using stored device fingerprint:', storedFingerprint);
        return storedFingerprint;
      }
    }
    
    // Generate new fingerprint
    const fp = await fpPromise;
    const result = await fp.get();
    const fingerprint = result.visitorId;
    
    // Store for future use
    localStorage.setItem('deviceFingerprint', fingerprint);
    localStorage.setItem('deviceFingerprintTime', Date.now().toString());
    
    console.log('Generated new device fingerprint:', fingerprint);
    return fingerprint;
  } catch (error) {
    console.error('Fingerprinting failed, using fallback:', error);
    
    // Fallback: Create a stable fingerprint from browser characteristics
    const fallbackData = {
      userAgent: navigator.userAgent,
      language: navigator.language,
      platform: navigator.platform,
      screenResolution: `${screen.width}x${screen.height}`,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    };
    
    // Create hash of fallback data
    const fallbackString = JSON.stringify(fallbackData);
    const fallbackFingerprint = btoa(fallbackString).substring(0, 32);
    
    // Store fallback
    localStorage.setItem('deviceFingerprint', fallbackFingerprint);
    localStorage.setItem('deviceFingerprintTime', Date.now().toString());
    
    return fallbackFingerprint;
  }
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

// Get location data from IP
const getLocationData = async (): Promise<{
  ip: string;
  country?: string;
  city?: string;
  latitude?: number;
  longitude?: number;
}> => {
  try {
    // Using ipapi.co - free tier allows 1000 requests/day
    const response = await fetch('https://ipapi.co/json/');
    const data = await response.json();
    
    return {
      ip: data.ip || '0.0.0.0',
      country: data.country_name,
      city: data.city,
      latitude: data.latitude,
      longitude: data.longitude,
    };
  } catch (error) {
    console.error('Failed to get location data:', error);
    return {
      ip: '0.0.0.0',
      country: undefined,
      city: undefined,
      latitude: undefined,
      longitude: undefined,
    };
  }
};

// Get complete device info
export const getDeviceInfo = async (): Promise<DeviceInfo> => {
  const fingerprint = await getDeviceFingerprint();
  const locationData = await getLocationData();
  
  return {
    deviceFingerprint: fingerprint,
    deviceName: getDeviceName(),
    deviceType: getDeviceType(),
    os: getOS(),
    browser: getBrowser(),
    ipAddress: locationData.ip,
    userAgent: navigator.userAgent,
    locationCountry: locationData.country,
    locationCity: locationData.city,
    latitude: locationData.latitude,
    longitude: locationData.longitude,
  };
};

// Format device info for display
export const formatDeviceInfo = (device: any): string => {
  if (device.deviceName) return device.deviceName;
  if (device.os && device.browser) return `${device.os} - ${device.browser}`;
  return 'Unknown Device';
};