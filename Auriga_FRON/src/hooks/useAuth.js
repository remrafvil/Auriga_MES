import { useState, useEffect } from 'react';
import { 
  isAuthenticated, 
  getCurrentUser, 
  checkServerAuth,
  fetchCompleteUserData,
  extractOrganizationFromJWT,
  getCookie
} from '../utils/auth';

export const useAuth = () => {
  const [authState, setAuthState] = useState({
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
        console.log('🔍 [useAuth] Iniciando verificación...');
        
        // Verificar autenticación primero
        const serverAuth = await checkServerAuth();
        console.log('🔐 [useAuth] Auth check:', serverAuth);
        
        if (serverAuth.authenticated) {
          console.log('🎉 [useAuth] Usuario autenticado, obteniendo datos...');
          
          try {
            // Obtener datos completos
            const completeUserData = await fetchCompleteUserData();
            console.log('📦 [useAuth] Datos completos recibidos:', completeUserData);
            
            // ✅ CORRECCIÓN: Extraer datos de from_context y organization_structure
            const userData = completeUserData.from_context || completeUserData.user || serverAuth.user;
            const organizationData = completeUserData.organization_structure || completeUserData.organization;
            
            console.log('👤 [useAuth] User data extraído:', userData);
            console.log('🏢 [useAuth] Organization data extraído:', organizationData);
            
            setAuthState({
              isAuthenticated: true,
              user: userData,
              serverAuth,
              apiUserData: {
                ...completeUserData,
                user: userData,
                organization: organizationData
              },
              isLoading: false,
              error: null
            });
          } catch (apiError) {
            console.error('❌ [useAuth] Error obteniendo datos:', apiError);
            // Usar datos del servidor como fallback
            setAuthState({
              isAuthenticated: true,
              user: serverAuth.user,
              serverAuth,
              apiUserData: null,
              isLoading: false,
              error: 'Error cargando datos adicionales'
            });
          }
        } else {
          setAuthState({
            isAuthenticated: false,
            user: null,
            serverAuth,
            apiUserData: null,
            isLoading: false,
            error: serverAuth.message
          });
        }
      } catch (error) {
        console.error('💥 [useAuth] Error:', error);
        setAuthState({
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

  return authState;
};