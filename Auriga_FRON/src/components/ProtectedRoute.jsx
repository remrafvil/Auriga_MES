import React, { useEffect, useState } from 'react'
import { checkServerAuth } from '../utils/auth'

const ProtectedRoute = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(null)

  useEffect(() => {
    const verifyAuth = async () => {
      try {
        const auth = await checkServerAuth()
        setIsAuthenticated(auth.authenticated)
        
        if (!auth.authenticated) {
          window.location.href = '/login'
        }
      } catch (error) {
        console.error('Auth verification failed:', error)
        setIsAuthenticated(false)
        window.location.href = '/login'
      }
    }

    verifyAuth()
  }, [])

  if (isAuthenticated === null) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh' 
      }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{
            border: '4px solid #f3f3f3',
            borderTop: '4px solid #007bff',
            borderRadius: '50%',
            width: '40px',
            height: '40px',
            animation: 'spin 1s linear infinite',
            margin: '0 auto 20px'
          }}></div>
          <p>Verificando autenticación...</p>
        </div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return null // Ya se redirigió a login
  }

  return children
}

export default ProtectedRoute