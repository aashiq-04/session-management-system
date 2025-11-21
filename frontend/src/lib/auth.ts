import Cookies from 'js-cookie';

const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';
const USER_ID_KEY = 'userId';

export const auth = {
  // Save tokens after login/register
  setTokens: (accessToken: string, refreshToken: string, userId: string) => {
    Cookies.set(ACCESS_TOKEN_KEY, accessToken, { expires: 7 }); // 7 days
    Cookies.set(REFRESH_TOKEN_KEY, refreshToken, { expires: 7 });
    Cookies.set(USER_ID_KEY, userId, { expires: 7 });
  },

  // Get access token
  getAccessToken: (): string | undefined => {
    return Cookies.get(ACCESS_TOKEN_KEY);
  },

  // Get refresh token
  getRefreshToken: (): string | undefined => {
    return Cookies.get(REFRESH_TOKEN_KEY);
  },

  // Get user ID
  getUserId: (): string | undefined => {
    return Cookies.get(USER_ID_KEY);
  },

  // Clear all tokens (logout)
  clearTokens: () => {
    Cookies.remove(ACCESS_TOKEN_KEY);
    Cookies.remove(REFRESH_TOKEN_KEY);
    Cookies.remove(USER_ID_KEY);
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    return !!Cookies.get(ACCESS_TOKEN_KEY);
  },
};