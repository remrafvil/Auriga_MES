import React, { useState, useEffect } from 'react'
import { useAuth } from '../hooks/useAuth'
import { logout, fetchUserData, fetchTokenInfo, fetchMyGroups, fetchProtectedData, extractOrganizationFromJWT, getCookie } from '../utils/auth'
import './Dashboard.css'

const Dashboard = () => {
  const { isAuthenticated, user, apiUserData, isLoading, error } = useAuth()
  const [apiResult, setApiResult] = useState('Los resultados de la API aparecerán aquí...')
  const [debugInfo, setDebugInfo] = useState(null)
  const [showDebug, setShowDebug] = useState(false)
  const [organizationData, setOrganizationData] = useState(null)

  console.log('📊 Dashboard render - Auth state:', { isAuthenticated, user, apiUserData, isLoading, error })

  // Función para procesar datos de fábricas - VERSIÓN MEJORADA
  const processOrganizationData = (orgData) => {
    if (!orgData) {
      console.log('❌ No organization data to process')
      return null
    }

    console.log('🔄 Processing organization data:', orgData)

    // ✅ CORRECCIÓN: Manejar tanto array como objeto de fábricas
    let factories = orgData.factories
    let groupedRoles = []
    let totalRoles = 0
    let departmentsCount = 0
    let factoryNames = []

    if (Array.isArray(factories)) {
      // Si factories es un array (formato actual)
      console.log('✅ Factories es array, convirtiendo a objeto...')
      
      const factoriesObj = {}
      groupedRoles = []
      
      factories.forEach(item => {
        if (!factoriesObj[item.factory]) {
          factoriesObj[item.factory] = { departments: {} }
          factoryNames.push(item.factory)
        }
        
        if (!factoriesObj[item.factory].departments[item.department]) {
          factoriesObj[item.factory].departments[item.department] = { roles: [] }
          departmentsCount++
        }
        
        if (item.roles && Array.isArray(item.roles)) {
          factoriesObj[item.factory].departments[item.department].roles = item.roles
          groupedRoles.push({
            factory: item.factory,
            department: item.department,
            roles: item.roles
          })
          totalRoles += item.roles.length
        }
      })
      
      factories = factoriesObj
    } else if (factories && typeof factories === 'object') {
      // Si factories ya es un objeto (formato esperado)
      console.log('✅ Factories es objeto, usando directamente...')
      factoryNames = Object.keys(factories)
      
      factoryNames.forEach(factoryName => {
        const factory = factories[factoryName]
        if (factory.departments && typeof factory.departments === 'object') {
          const departmentNames = Object.keys(factory.departments)
          departmentsCount += departmentNames.length
          
          departmentNames.forEach(deptName => {
            const department = factory.departments[deptName]
            if (department.roles && Array.isArray(department.roles)) {
              groupedRoles.push({
                factory: factoryName,
                department: deptName,
                roles: department.roles
              })
              totalRoles += department.roles.length
            }
          })
        }
      })
    }

    const processedData = {
      name: orgData.name || 'Coexpan',
      workday_id: orgData.workday_id,
      idn: orgData.idn,
      status: orgData.status,
      departamento: orgData.departamento,
      cargo: orgData.cargo,
      tipo_empleado: orgData.tipo_empleado,
      factories: factories,
      grouped_roles: groupedRoles,
      stats: {
        factories_count: factoryNames.length,
        departments_count: departmentsCount,
        total_roles: totalRoles
      }
    }

    console.log('✅ Processed organization data:', processedData)
    return processedData
  }

  // DEBUG DETALLADO
  useEffect(() => {
    console.log('🔄 Dashboard useEffect - apiUserData changed:', apiUserData)
    
    if (apiUserData) {
      console.log('=== 🎯 DEBUG DETALLADO ===')
      console.log('1. apiUserData completo:', apiUserData)
      console.log('2. Tiene propiedad user?:', !!apiUserData.user)
      console.log('3. Tiene propiedad organization?:', !!apiUserData.organization)
      console.log('4. User object:', apiUserData.user)
      console.log('5. Organization object:', apiUserData.organization)
      console.log('6. Tipo de apiUserData:', typeof apiUserData)
      console.log('7. Keys de apiUserData:', Object.keys(apiUserData))
      
      let orgData = apiUserData.organization
      console.log('📦 Organization data from API:', orgData)
      
      // Si no hay datos de organización, intentar extraer del JWT
      if (!orgData) {
        console.log('🔄 No organization data from API, checking JWT...')
        const authToken = getCookie('auth_token')
        if (authToken) {
          orgData = extractOrganizationFromJWT(authToken)
          console.log('📦 Organization data from JWT:', orgData)
        }
      }
      
      // Procesar los datos de organización
      const processedOrgData = processOrganizationData(orgData)
      console.log('🔄 Processed organization data:', processedOrgData)
      setOrganizationData(processedOrgData)
      
      // Configurar información de debug
      setDebugInfo({
        rawOrganization: orgData,
        processedOrganization: processedOrgData,
        apiUserData: apiUserData,
        user: user,
        authToken: getCookie('auth_token') ? 'Present' : 'Missing',
        debugDetails: {
          hasUser: !!apiUserData.user,
          hasOrganization: !!apiUserData.organization,
          userKeys: apiUserData.user ? Object.keys(apiUserData.user) : [],
          organizationKeys: orgData ? Object.keys(orgData) : []
        }
      })
    } else {
      console.log('❌ No apiUserData available')
    }
  }, [apiUserData, user])

  const handleLogout = () => {
    console.log('🚪 Logging out...')
    logout()
  }

  const showLoading = (message) => {
    setApiResult(message)
  }

  const loadProfile = async () => {
    try {
      showLoading('Cargando perfil...')
      const data = await fetchUserData()
      setApiResult(JSON.stringify(data, null, 2))
    } catch (error) {
      setApiResult(`Error: ${error.message}`)
    }
  }

  const loadUserData = async () => {
    try {
      showLoading('Cargando datos de usuario...')
      const data = await fetchUserData()
      setApiResult(JSON.stringify(data, null, 2))
    } catch (error) {
      setApiResult(`Error: ${error.message}`)
    }
  }

  const loadMyGroups = async () => {
    try {
      showLoading('Cargando grupos y roles...')
      const data = await fetchMyGroups()
      setApiResult(JSON.stringify(data, null, 2))
      
      if (data.organization_structure) {
        const processedOrgData = processOrganizationData(data.organization_structure)
        setOrganizationData(processedOrgData)
      }
    } catch (error) {
      setApiResult(`Error: ${error.message}`)
    }
  }

  const loadTokenInfo = async () => {
    try {
      showLoading('Cargando información del token...')
      const data = await fetchTokenInfo()
      setApiResult(JSON.stringify(data, null, 2))
    } catch (error) {
      setApiResult(`Error: ${error.message}`)
    }
  }

  const loadProtectedData = async () => {
    try {
      showLoading('Cargando datos protegidos...')
      const data = await fetchProtectedData()
      setApiResult(JSON.stringify(data, null, 2))
    } catch (error) {
      setApiResult(`Error: ${error.message}`)
    }
  }

  const loadAdminStats = async () => {
    try {
      showLoading('Cargando estadísticas de administrador...')
      const response = await fetch('/api/admin/stats', {
        credentials: 'include'
      })
      
      if (response.ok) {
        const data = await response.json()
        setApiResult(JSON.stringify(data, null, 2))
      } else if (response.status === 403) {
        setApiResult('❌ Acceso denegado: No tienes permisos de administrador')
      } else {
        setApiResult(`Error: ${response.status} - ${response.statusText}`)
      }
    } catch (error) {
      setApiResult(`Error de conexión: ${error.message}`)
    }
  }

  const reloadUserData = async () => {
    try {
      showLoading('Recargando datos...')
      const response = await fetch('/api/users/me', {
        credentials: 'include'
      });
      
      if (response.ok) {
        const data = await response.json();
        setApiResult(JSON.stringify(data, null, 2));
        console.log('✅ Datos recargados:', data);
        
        // Procesar organización nuevamente
        const orgData = data.organization || data.organization_structure;
        const processedOrgData = processOrganizationData(orgData);
        setOrganizationData(processedOrgData);
      } else {
        setApiResult(`Error: ${response.status}`);
      }
    } catch (error) {
      setApiResult(`Error: ${error.message}`);
    }
  };

  const toggleDebug = () => {
    setShowDebug(!showDebug)
  }

  if (isLoading) {
    return (
      <div className="dashboard-loading">
        <div className="loading-spinner"></div>
        <p>Cargando...</p>
      </div>
    )
  }

  if (error || !isAuthenticated) {
    return (
      <div className="dashboard-container">
        <div className="error-message">
          <strong>Error de autenticación:</strong> {error || 'No autenticado'}
          <br />
          <small>Verifica la consola para más detalles</small>
        </div>
        <button 
          onClick={() => window.location.href = '/login'}
          className="primary-button"
        >
          Ir al Login
        </button>
      </div>
    )
  }

  const hasAdminAccess = user?.groups?.includes('authentik Admins') || false
  const hasGrafanaAccess = user?.groups?.includes('Grafana Admins') || 
                          user?.groups?.includes('Grafana Editors') || false

  // Función para renderizar fábricas
  const renderFactories = () => {
    if (organizationData?.grouped_roles && organizationData.grouped_roles.length > 0) {
      return (
        <div className="factory-structure">
          <h4>🏭 Estructura de Fábricas y Roles</h4>
          {organizationData.grouped_roles.map((factory, index) => (
            <div key={index} className="factory-item">
              <strong>🏭 Fábrica:</strong> {factory.factory}
              <div className="department-item">
                <strong>📋 Departamento:</strong> {factory.department}<br/>
                <strong>🎯 Roles:</strong> 
                {factory.roles && factory.roles.length > 0 ? (
                  factory.roles.map((role, roleIndex) => (
                    <span key={roleIndex} className="role-badge">{role}</span>
                  ))
                ) : (
                  <em>Sin roles asignados</em>
                )}
              </div>
            </div>
          ))}
        </div>
      )
    } else if (organizationData?.factories) {
      return (
        <div className="factory-structure">
          <h4>🏭 Estructura de Fábricas</h4>
          {Object.entries(organizationData.factories).map(([factoryName, factoryData], index) => (
            <div key={index} className="factory-item">
              <strong>🏭 Fábrica:</strong> {factoryName}
              {factoryData.departments && Object.entries(factoryData.departments).map(([deptName, deptData], deptIndex) => (
                <div key={deptIndex} className="department-item">
                  <strong>📋 Departamento:</strong> {deptName}<br/>
                  <strong>🎯 Roles:</strong> 
                  {deptData.roles && deptData.roles.length > 0 ? (
                    deptData.roles.map((role, roleIndex) => (
                      <span key={roleIndex} className="role-badge">{role}</span>
                    ))
                  ) : (
                    <em>Sin roles asignados</em>
                  )}
                </div>
              ))}
            </div>
          ))}
        </div>
      )
    } else {
      return (
        <div className="info-message">
          🏢 Eres parte de Coexpan<br/>
          📋 Pero no tienes fábricas específicas asignadas<br/>
          💡 Contacta al administrador para asignarte fábricas y roles
          <br/>
          <button onClick={loadMyGroups} className="primary-button" style={{ marginTop: '10px' }}>
            Intentar Cargar Fábricas
          </button>
        </div>
      )
    }
  }

  return (
    <div className="dashboard-container">
      <h1>🎉 ¡Bienvenido a tu Dashboard!</h1>
      
      {/* Sección de Debug */}
      <div className="debug-section">
        <div className="section-title">
          <h3>🐛 Información de Debug</h3>
        </div>
        <button className="debug-toggle" onClick={toggleDebug}>
          🔍 {showDebug ? 'Ocultar' : 'Mostrar'} Información de Debug
        </button>
        {showDebug && debugInfo && (
          <div className="debug-content">
            <h4>Estado de Autenticación:</h4>
            <pre>{JSON.stringify({
              isAuthenticated,
              user,
              hasUserData: !!user,
              hasApiUserData: !!apiUserData,
              hasOrganizationData: !!organizationData,
              authToken: debugInfo.authToken
            }, null, 2)}</pre>
            
            <h4>Datos del Usuario:</h4>
            <pre>{JSON.stringify(user, null, 2)}</pre>
            
            <h4>Datos de la API:</h4>
            <pre>{JSON.stringify(apiUserData, null, 2)}</pre>
            
            <h4>Organización Procesada:</h4>
            <pre>{JSON.stringify(organizationData, null, 2)}</pre>

            <h4>Debug Details:</h4>
            <pre>{JSON.stringify(debugInfo.debugDetails, null, 2)}</pre>
          </div>
        )}
      </div>

      {/* Sección de estadísticas */}
      <div className="stats-grid">
        {organizationData?.stats ? (
          <>
            <div className="stat-card">
              <div className="stat-number">{organizationData.stats.factories_count || 0}</div>
              <div className="stat-label">Fábricas</div>
            </div>
            <div className="stat-card">
              <div className="stat-number">{organizationData.stats.departments_count || 0}</div>
              <div className="stat-label">Departamentos</div>
            </div>
            <div className="stat-card">
              <div className="stat-number">{organizationData.stats.total_roles || 0}</div>
              <div className="stat-label">Roles Totales</div>
            </div>
          </>
        ) : (
          <>
            <div className="stat-card">
              <div className="stat-number">0</div>
              <div className="stat-label">Fábricas</div>
            </div>
            <div className="stat-card">
              <div className="stat-number">0</div>
              <div className="stat-label">Departamentos</div>
            </div>
            <div className="stat-card">
              <div className="stat-number">0</div>
              <div className="stat-label">Roles Totales</div>
            </div>
          </>
        )}
      </div>
      
      {/* Información del Usuario */}
      <div className="user-info">
        <div className="section-title">
          <h3>👤 Información del Usuario</h3>
        </div>
        {user ? (
          <div>
            <p><strong>Nombre:</strong> {user.name || 'No disponible'}</p>
            <p><strong>Email:</strong> {user.email || 'No disponible'}</p>
            <p><strong>Authentik ID:</strong> {user.id ? (user.id.length > 20 ? user.id.substring(0, 20) + '...' : user.id) : 'No disponible'}</p>
            {user.groups && user.groups.length > 0 && (
              <div style={{ marginTop: '10px' }}>
                <strong>Grupos:</strong>
                <div style={{ marginTop: '5px' }}>
                  {user.groups.map((group, index) => (
                    <span key={index} className="group-badge">{group}</span>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="warning-message">
            ❌ No se pudo cargar la información del usuario
            <br/>
            <button onClick={loadUserData} className="primary-button" style={{ marginTop: '10px' }}>
              Cargar Datos del Usuario
            </button>
          </div>
        )}
      </div>

      {/* Información de la Organización */}
      <div className="organization-info">
        <div className="section-title">
          <h3>🏢 Información de la Organización</h3>
        </div>
        {organizationData ? (
          <div>
            <p><strong>Organización:</strong> {organizationData.name || 'Coexpan'}</p>
            
            {/* Detalles del Empleado */}
            {(organizationData.workday_id || organizationData.idn || organizationData.status || organizationData.departamento || organizationData.cargo || organizationData.tipo_empleado) && (
              <div className="employee-details">
                <h4>📋 Detalles del Empleado</h4>
                {organizationData.workday_id && (
                  <div className="detail-row">
                    <span className="detail-label">Workday ID:</span>
                    <span className="detail-value">{organizationData.workday_id}</span>
                  </div>
                )}
                {organizationData.idn && (
                  <div className="detail-row">
                    <span className="detail-label">IDN:</span>
                    <span className="detail-value">{organizationData.idn}</span>
                  </div>
                )}
                {organizationData.status && (
                  <div className="detail-row">
                    <span className="detail-label">Estado:</span>
                    <span className="detail-value">{organizationData.status}</span>
                  </div>
                )}
                {organizationData.departamento && (
                  <div className="detail-row">
                    <span className="detail-label">Departamento:</span>
                    <span className="detail-value">{organizationData.departamento}</span>
                  </div>
                )}
                {organizationData.cargo && (
                  <div className="detail-row">
                    <span className="detail-label">Cargo:</span>
                    <span className="detail-value">{organizationData.cargo}</span>
                  </div>
                )}
                {organizationData.tipo_empleado && (
                  <div className="detail-row">
                    <span className="detail-label">Tipo Empleado:</span>
                    <span className="detail-value">{organizationData.tipo_empleado}</span>
                  </div>
                )}
              </div>
            )}

            {/* Estructura de Fábricas */}
            {renderFactories()}
          </div>
        ) : (
          <div className="warning-message">
            ❌ No se encontró información de organización específica.
            <br/>
            <button onClick={loadMyGroups} className="primary-button" style={{ marginTop: '10px' }}>
              Cargar Información de Organización
            </button>
          </div>
        )}
      </div>

      {/* Panel de Administración */}
      {(hasAdminAccess || hasGrafanaAccess) && (
        <div className="admin-section">
          <div className="section-title">
            <h3>⚙️ Panel de Administración</h3>
          </div>
          <button onClick={loadAdminStats} className="primary-button">
            📊 Estadísticas del Sistema
          </button>
        </div>
      )}

      {/* Botón de Logout */}
      <div className="logout-section">
        <button onClick={handleLogout} className="logout-button">
          🚪 Cerrar Sesión
        </button>
      </div>

      {/* Sección de API */}
      <div className="api-section">
        <div className="section-title">
          <h3>🔧 Datos de la API</h3>
        </div>
        <div className="api-buttons">
          <button onClick={loadProfile} className="primary-button">📊 Cargar Perfil</button>
          <button onClick={loadUserData} className="primary-button">👤 Cargar Datos de Usuario</button>
          <button onClick={reloadUserData} className="primary-button">🔄 Recargar Datos</button>
          <button onClick={loadMyGroups} className="primary-button">👥 Ver Mis Grupos y Roles</button>
          <button onClick={loadTokenInfo} className="primary-button">⏰ Info del Token</button>
          <button onClick={loadProtectedData} className="primary-button">🔒 Datos Protegidos</button>
        </div>
        
        <div className="api-result">
          <pre>{apiResult}</pre>
        </div>
      </div>
    </div>
  )
}

export default Dashboard