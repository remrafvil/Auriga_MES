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

// En utils/auth.js
export const checkServerAuth = async () => {
  try {
    console.log('🔐 Checking server authentication...')
    
    const response = await fetch('/api/auth/check', {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    console.log('🔐 Auth check response status:', response.status)
    
    if (!response.ok) {
      console.log('❌ Auth check failed with status:', response.status)
      return { 
        authenticated: false, 
        message: `HTTP error! status: ${response.status}` 
      }
    }

    const data = await response.json()
    console.log('✅ Auth check response data:', data)
    
    return data
  } catch (error) {
    console.error('❌ Error checking server auth:', error)
    return { 
      authenticated: false, 
      message: error.message 
    }
  }
}

// Función mejorada para extraer organización del JWT
export const extractOrganizationFromJWT = (jwtToken) => {
  try {
    if (!jwtToken) {
      console.log('❌ No JWT token provided')
      return null
    }
    
    console.log('🔐 Extracting organization from JWT...')
    const parts = jwtToken.split('.')
    if (parts.length !== 3) {
      console.log('❌ Invalid JWT format')
      return null
    }
    
    const payload = JSON.parse(atob(parts[1]))
    console.log('✅ JWT payload:', payload)
    
    if (payload.organization) {
      console.log('✅ Organization found in JWT:', payload.organization)
      return payload.organization
    } else {
      console.log('❌ No organization found in JWT')
      return null
    }
  } catch (error) {
    console.error('❌ Error extracting organization from JWT:', error)
    return null
  }
}

// Función para obtener el token JWT de las cookies
export const getAuthToken = () => {
  const token = getCookie('auth_token')
  console.log('🔐 Auth token from cookies:', token ? 'Found' : 'Not found')
  return token
}

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
// En auth.js, mejora esta función:
export const fetchCompleteUserData = async () => {
  try {
    console.log('🔄 [fetchCompleteUserData] Obteniendo datos completos...');
    
    // Obtener datos del usuario
    const userResponse = await fetch('/api/users/me', {
      credentials: 'include'
    });
    
    if (!userResponse.ok) {
      throw new Error(`Error usuarios: ${userResponse.status}`);
    }
    
    const userData = await userResponse.json();
    console.log('✅ [fetchCompleteUserData] Datos recibidos:', userData);
    
    return userData;
  } catch (error) {
    console.error('❌ [fetchCompleteUserData] Error:', error);
    throw error;
  }
};
export const fetchUserData = () => apiRequest('/users/me');
export const fetchTokenInfo = () => apiRequest('/token-info');
export const fetchMyGroups = () => apiRequest('/my-groups');
export const fetchProtectedData = () => apiRequest('/protected-data');


