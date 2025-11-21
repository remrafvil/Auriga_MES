import { useState, useEffect } from 'react';
import { 
  isAuthenticated, 
  getCurrentUser, 
  checkServerAuth,
  fetchUserData 
} from '../utils/auth';

export const useAuth = () => {
  const [auth, setAuth] = useState({
    isAuthenticated: false,
    user: null,
    serverAuth: null,
    apiUserData: null,
    isLoading: true,
    error: null
  });

  useEffect(() => {
    const checkAuth = async () => {
      try {
        // Primero verificar autenticación en el servidor
        const serverAuth = await checkServerAuth();
        
        if (serverAuth.authenticated) {
          // Si el servidor dice que está autenticado, obtener datos completos
          const apiUserData = await fetchUserData();
          
          setAuth({
            isAuthenticated: true,
            user: serverAuth.user || getCurrentUser(),
            serverAuth,
            apiUserData,
            isLoading: false,
            error: null
          });
        } else {
          // No autenticado en el servidor
          setAuth({
            isAuthenticated: false,
            user: null,
            serverAuth,
            apiUserData: null,
            isLoading: false,
            error: serverAuth.message || 'Not authenticated'
          });
        }
      } catch (error) {
        console.error('Auth check error:', error);
        setAuth({
          isAuthenticated: false,
          user: null,
          serverAuth: null,
          apiUserData: null,
          isLoading: false,
          error: error.message
        });
      }
    };

    checkAuth();
  }, []);

  return auth;
};