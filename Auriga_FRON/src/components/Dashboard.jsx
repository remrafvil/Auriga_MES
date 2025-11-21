import React from 'react'
import { useAuth } from '../hooks/useAuth'
import { logout } from '../utils/auth'

const Dashboard = () => {
  const { auth, isLoading, error } = useAuth()

  const handleLogout = () => {
    logout()
  }

  if (isLoading) {
    return (
      <div style={{ padding: '20px', textAlign: 'center' }}>
        <div style={{
          border: '4px solid #f3f3f3',
          borderTop: '4px solid #007bff',
          borderRadius: '50%',
          width: '40px',
          height: '40px',
          animation: 'spin 1s linear infinite',
          margin: '0 auto 20px'
        }}></div>
        <p>Cargando...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div style={{ padding: '20px' }}>
        <div style={{ 
          background: '#f8d7da', 
          color: '#721c24', 
          padding: '15px', 
          borderRadius: '5px',
          marginBottom: '20px'
        }}>
          <strong>Error:</strong> {error}
        </div>
        <button onClick={() => window.location.href = '/login'}>
          Ir al Login
        </button>
      </div>
    )
  }

  return (
    <div style={{ padding: '20px' }}>
      <h1>Dashboard</h1>
      
      {auth.user && (
        <div style={{ 
          background: '#e8f5e8', 
          padding: '15px', 
          borderRadius: '5px',
          marginBottom: '20px'
        }}>
          <h2>Bienvenido, {auth.user.name}!</h2>
          <p><strong>Email:</strong> {auth.user.email}</p>
          <p><strong>ID:</strong> {auth.user.id}</p>
        </div>
      )}

      {auth.apiUserData && (
        <div style={{ 
          background: '#e3f2fd', 
          padding: '15px', 
          borderRadius: '5px',
          marginBottom: '20px'
        }}>
          <h3>Información de la API</h3>
          <pre style={{ fontSize: '12px', overflow: 'auto' }}>
            {JSON.stringify(auth.apiUserData, null, 2)}
          </pre>
        </div>
      )}

      <div style={{ marginTop: '20px' }}>
        <button 
          onClick={handleLogout}
          style={{
            background: '#dc3545',
            color: 'white',
            border: 'none',
            padding: '10px 20px',
            borderRadius: '5px',
            cursor: 'pointer'
          }}
        >
          Cerrar Sesión
        </button>
      </div>

      <style>{`
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  )
}

export default Dashboard