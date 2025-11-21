// Utilidades para manejar autenticación y cookies
const API_BASE_URL = ''; // Usar rutas relativas por el proxy

export const getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
  return null;
};

export const getUserSession = () => {
  try {
    const cookie = getCookie('user_data');
    if (cookie) {
      const decoded = atob(cookie);
      return JSON.parse(decoded);
    }
  } catch (error) {
    console.error('Error reading user session:', error);
  }
  return null;
};

export const isAuthenticated = () => {
  const sessionActive = getCookie('session_active') === 'true';
  const authToken = getCookie('auth_token');
  return !!(sessionActive && authToken);
};

export const getCurrentUser = () => {
  const session = getUserSession();
  return session || null;
};

export const checkServerAuth = async () => {
  try {
    const response = await fetch('/api/auth/check', {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error checking server auth:', error);
    return { authenticated: false };
  }
};

export const logout = async () => {
  try {
    await fetch('/auth/logout', {
      method: 'POST',
      credentials: 'include'
    });
  } catch (error) {
    console.error('Logout error:', error);
  } finally {
    // Limpiar cookies del lado del cliente
    document.cookie = 'user_data=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    document.cookie = 'session_active=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    
    // Redirigir al login
    window.location.href = '/login';
  }
};

// API calls
export const apiRequest = async (endpoint, options = {}) => {
  const defaultOptions = {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  try {
    const response = await fetch(`/api${endpoint}`, {
      ...defaultOptions,
      ...options,
    });

    if (response.status === 401) {
      // No autenticado
      window.location.href = '/login';
      return null;
    }

    if (response.status === 403) {
      throw new Error('No tienes permisos para acceder a este recurso');
    }

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
};

export const fetchUserData = () => apiRequest('/users/me');
export const fetchTokenInfo = () => apiRequest('/token-info');
export const fetchMyGroups = () => apiRequest('/my-groups');
export const fetchProtectedData = () => apiRequest('/protected-data');