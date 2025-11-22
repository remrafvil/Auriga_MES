import React, { useEffect, useState } from 'react'
import { checkServerAuth } from '../utils/auth'
import './Login.css' // Crearemos este archivo

const Login = () => {
  const [checkingAuth, setCheckingAuth] = useState(true)

  useEffect(() => {
    const verifyAuth = async () => {
      try {
        const auth = await checkServerAuth()
        console.log('Auth check result:', auth)
        
        if (auth.authenticated) {
          console.log('User is authenticated, redirecting to dashboard')
          window.location.href = '/dashboard'
          return
        }
        
        // Si no está autenticado, mostrar la página de login
        setCheckingAuth(false)
      } catch (error) {
        console.error('Auth verification failed:', error)
        setCheckingAuth(false)
      }
    }

    verifyAuth()
  }, [])

  const handleLogin = () => {
    // Redirigir al endpoint de login de la API
    window.location.href = '/auth/login'
  }

  if (checkingAuth) {
    return (
      <div className="login-container">
        <div className="login-card">
          <div className="loading">
            <div className="spinner"></div>
            <p>Verificando autenticación...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="login-container">
      <div className="login-card">
        <div className="header">
          <h1 className="title">🔐 Iniciar Sesión</h1>
          <p className="subtitle">
            Para acceder a la aplicación, inicia sesión con tu cuenta de Authentik
          </p>
        </div>
        
        <button 
          onClick={handleLogin}
          className="login-button"
        >
          Iniciar Sesión con Authentik
        </button>

        <div className="info">
          <p><strong>Flujo:</strong> React → API Go (8081) → Authentik → API Go → React</p>
          <p><strong>Estado:</strong> No autenticado - Puede iniciar sesión</p>
        </div>
      </div>
    </div>
  )
}

export default Login