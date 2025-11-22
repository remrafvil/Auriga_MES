import React, { useEffect, useState } from 'react'
import { checkServerAuth } from '../utils/auth'
import './ProtectedRoute.css'

const ProtectedRoute = ({ children }) => {
  const [authStatus, setAuthStatus] = useState('checking') // 'checking', 'authenticated', 'unauthenticated'

  useEffect(() => {
    const verifyAuth = async () => {
      try {
        const auth = await checkServerAuth()
        console.log('ProtectedRoute auth check:', auth)
        
        if (auth.authenticated) {
          setAuthStatus('authenticated')
        } else {
          setAuthStatus('unauthenticated')
          window.location.href = '/login'
        }
      } catch (error) {
        console.error('Auth verification failed in ProtectedRoute:', error)
        setAuthStatus('unauthenticated')
        window.location.href = '/login'
      }
    }

    verifyAuth()
  }, [])

  if (authStatus === 'checking') {
    return (
      <div className="protected-route-loading">
        <div className="loading-spinner"></div>
        <p>Verificando autenticación...</p>
      </div>
    )
  }

  if (authStatus === 'unauthenticated') {
    return null // Ya se redirigió a login
  }

  return children
}

export default ProtectedRoute