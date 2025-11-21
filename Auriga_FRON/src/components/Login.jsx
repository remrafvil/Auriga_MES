import React, { useEffect } from 'react'

const Login = () => {
  useEffect(() => {
    const sessionActive = document.cookie.includes('session_active=true')
    if (sessionActive) {
      window.location.href = '/dashboard'
    }
  }, [])

  const handleLogin = () => {
    // Usar ruta relativa - Vite proxy redirigirá a la API Go
    window.location.href = '/auth/login'
  }

  return (
    <div style={styles.container}>
      <div style={styles.loginCard}>
        <div style={styles.header}>
          <h1 style={styles.title}>🔐 Iniciar Sesión</h1>
          <p style={styles.subtitle}>
            Para acceder a la aplicación, inicia sesión con tu cuenta de Authentik
          </p>
        </div>
        
        <button 
          onClick={handleLogin}
          style={styles.loginButton}
        >
          Iniciar Sesión con Authentik
        </button>

        <div style={styles.info}>
          <p><strong>Flujo:</strong> React → API Go (8081) → Authentik → API Go → React</p>
        </div>
      </div>
    </div>
  )
}

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '100vh',
    backgroundColor: '#f5f5f5',
    padding: '20px',
  },
  loginCard: {
    background: 'white',
    padding: '40px',
    borderRadius: '12px',
    boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
    textAlign: 'center',
    maxWidth: '450px',
    width: '100%',
  },
  header: {
    marginBottom: '30px',
  },
  title: {
    fontSize: '28px',
    color: '#333',
    marginBottom: '10px',
  },
  subtitle: {
    fontSize: '16px',
    color: '#666',
    lineHeight: '1.5',
  },
  loginButton: {
    background: '#007bff',
    color: 'white',
    border: 'none',
    padding: '15px 30px',
    borderRadius: '8px',
    fontSize: '16px',
    cursor: 'pointer',
    width: '100%',
    fontWeight: '600',
    marginBottom: '20px',
  },
  info: {
    padding: '15px',
    backgroundColor: '#f8f9fa',
    borderRadius: '6px',
    fontSize: '12px',
    color: '#555',
  },
}

export default Login
