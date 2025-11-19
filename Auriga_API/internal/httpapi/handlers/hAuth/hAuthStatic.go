package hAuth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *handler) loginPageHandler(c echo.Context) error {
	html := `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - Mi App</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background-color: #f5f5f5;
            margin: 0;
        }
        .login-container {
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.1);
            text-align: center;
        }
        .login-btn {
            background: #007bff;
            color: white;
            border: none;
            padding: 15px 30px;
            border-radius: 5px;
            font-size: 16px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
        .login-btn:hover {
            background: #0056b3;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>🔐 Iniciar Sesión</h1>
        <p>Para acceder a la aplicación, inicia sesión con tu cuenta de Authentik</p>
        <a href="/auth/login" class="login-btn">Iniciar Sesión con Authentik</a>
    </div>
</body>
</html>
    `
	return c.HTML(http.StatusOK, html)
}

func (h *handler) dashboardHandler(c echo.Context) error {
	// DEBUG: Forzar obtener la organización del contexto nuevamente
	claims, err := getUserClaimsFromContext(c)
	if err != nil {
		h.logger.Warn("No claims in dashboard", zap.Error(err))
	} else {
		h.logger.Debug("Claims in dashboard",
			zap.Any("organization", claims["organization"]),
			zap.String("organization_type", fmt.Sprintf("%T", claims["organization"])))
	}

	userID, userEmail, userName, userGroups, organization := getUserInfoFromContext(c)

	h.logger.Debug("Serving dashboard",
		zap.String("user_id", userID),
		zap.String("email", userEmail),
		zap.String("name", userName),
		zap.Strings("groups", userGroups),
		zap.Any("organization_raw", organization))

	// Obtener información COMPLETA de factories y roles usando las nuevas funciones
	groupedRoles := getUserRolesGrouped(organization)
	allRoles := getAllUserRoles(organization)
	orgSummary := getOrganizationSummary(organization)

	// OBTENER INFORMACIÓN DE DEBUG
	debugInfo := debugOrganizationStructure(organization)
	debugJSON := "{}"
	if len(debugInfo) > 0 {
		jsonBytes, _ := json.Marshal(debugInfo)
		debugJSON = string(jsonBytes)
	}

	// Convertir a JSON para el JavaScript
	groupedRolesJSON := "[]"
	if len(groupedRoles) > 0 {
		jsonBytes, _ := json.Marshal(groupedRoles)
		groupedRolesJSON = string(jsonBytes)
	}

	allRolesJSON := "[]"
	if len(allRoles) > 0 {
		jsonBytes, _ := json.Marshal(allRoles)
		allRolesJSON = string(jsonBytes)
	}

	orgSummaryJSON := "{}"
	if len(orgSummary) > 0 {
		jsonBytes, _ := json.Marshal(orgSummary)
		orgSummaryJSON = string(jsonBytes)
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard - Mi App</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .user-info, .organization-info, .groups-info, .api-section, .debug-section {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 5px;
            margin: 15px 0;
            border-left: 4px solid #007bff;
        }
        .organization-info {
            border-left-color: #28a745;
        }
        .groups-info {
            border-left-color: #ffc107;
        }
        .api-section {
            border-left-color: #6f42c1;
        }
        .debug-section {
            border-left-color: #6c757d;
        }
        .factory-info {
            background: #fff3e0;
            padding: 12px;
            margin: 8px 0;
            border-radius: 4px;
            border-left: 3px solid #fd7e14;
        }
        .group-badge {
            background: #2196f3;
            color: white;
            padding: 6px 12px;
            border-radius: 15px;
            font-size: 12px;
            margin: 3px;
            display: inline-block;
            font-weight: 500;
        }
        .role-badge {
            background: #4caf50;
            color: white;
            padding: 4px 10px;
            border-radius: 12px;
            font-size: 11px;
            margin: 2px;
            display: inline-block;
        }
        .admin-section {
            background: #ffebee;
            padding: 20px;
            border-radius: 5px;
            margin: 15px 0;
            border-left: 4px solid #f44336;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin: 15px 0;
        }
        .stat-card {
            background: white;
            padding: 15px;
            border-radius: 5px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            border: 1px solid #e9ecef;
        }
        .stat-number {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
        }
        .stat-label {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
        button {
            background: #007bff;
            color: white;
            border: none;
            padding: 10px 16px;
            border-radius: 4px;
            cursor: pointer;
            margin: 5px;
            font-size: 14px;
            transition: background 0.3s;
        }
        button:hover {
            background: #0056b3;
        }
        .logout-btn {
            background: #dc3545;
        }
        .logout-btn:hover {
            background: #c82333;
        }
        .debug-toggle {
            background: #6c757d;
            color: white;
            border: none;
            padding: 8px 12px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 12px;
            margin: 5px 0;
        }
        .debug-toggle:hover {
            background: #5a6268;
        }
        .api-result {
            background: #e9ecef;
            padding: 15px;
            border-radius: 4px;
            margin-top: 10px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            white-space: pre-wrap;
            max-height: 400px;
            overflow-y: auto;
        }
        .debug-content {
            display: none;
            background: #e9ecef;
            padding: 15px;
            border-radius: 4px;
            margin-top: 10px;
            font-family: 'Courier New', monospace;
            font-size: 11px;
            white-space: pre-wrap;
            max-height: 400px;
            overflow-y: auto;
            border: 1px solid #ced4da;
        }
        h1, h2, h3 {
            color: #333;
            margin-top: 0;
        }
        .section-title {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
        .section-title h3 {
            margin: 0;
            display: flex;
            align-items: center;
        }
        .section-title h3:before {
            content: "";
            display: inline-block;
            width: 20px;
            height: 20px;
            margin-right: 8px;
        }
        .user-info .section-title h3:before {
            content: "👤";
        }
        .organization-info .section-title h3:before {
            content: "🏢";
        }
        .groups-info .section-title h3:before {
            content: "👥";
        }
        .api-section .section-title h3:before {
            content: "🔧";
        }
        .admin-section .section-title h3:before {
            content: "⚙️";
        }
        .debug-section .section-title h3:before {
            content: "🐛";
        }
        .warning-message {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .info-message {
            background: #d1ecf1;
            border: 1px solid #bee5eb;
            color: #0c5460;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .success-message {
            background: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .employee-details {
            background: #e8f5e8;
            padding: 15px;
            border-radius: 5px;
            margin: 10px 0;
            border-left: 4px solid #28a745;
        }
        .detail-row {
            display: flex;
            justify-content: space-between;
            margin: 5px 0;
            padding: 5px 0;
            border-bottom: 1px solid #dee2e6;
        }
        .detail-label {
            font-weight: bold;
            color: #495057;
        }
        .detail-value {
            color: #28a745;
            font-weight: 500;
        }
        .factory-structure {
            background: #e3f2fd;
            padding: 15px;
            border-radius: 5px;
            margin: 10px 0;
            border-left: 4px solid #2196f3;
        }
        .factory-item {
            background: #bbdefb;
            padding: 10px;
            margin: 8px 0;
            border-radius: 4px;
            border-left: 3px solid #1976d2;
        }
        .department-item {
            background: #e1f5fe;
            padding: 8px;
            margin: 5px 0 5px 15px;
            border-radius: 3px;
            border-left: 2px solid #0288d1;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎉 ¡Bienvenido a tu Dashboard!</h1>
        
        <!-- Sección de Debug -->
        <div class="debug-section">
            <div class="section-title">
                <h3>Información de Debug</h3>
            </div>
            <button class="debug-toggle" onclick="toggleDebug()">🔍 Mostrar/Ocultar Información de Debug</button>
            <div id="debug-content" class="debug-content">
                <!-- La información de debug se cargará aquí -->
            </div>
        </div>

        <!-- Sección de estadísticas -->
        <div class="stats-grid" id="stats-section">
            <!-- Las estadísticas se cargarán con JavaScript -->
        </div>
        
        <div class="user-info">
            <div class="section-title">
                <h3>Información del Usuario</h3>
            </div>
            <div id="user-display">Cargando información del usuario...</div>
        </div>

        <div class="organization-info">
            <div class="section-title">
                <h3>Información de la Organización</h3>
            </div>
            <div id="organization-display">Cargando información de la organización...</div>
        </div>

        <div class="groups-info">
            <div class="section-title">
                <h3>Tus Grupos de Acceso</h3>
            </div>
            <div id="groups-display">Cargando grupos...</div>
        </div>

        <div id="admin-section" class="admin-section" style="display: none;">
            <div class="section-title">
                <h3>Panel de Administración</h3>
            </div>
            <button onclick="loadAdminStats()">📊 Estadísticas del Sistema</button>
            <button onclick="loadGrafanaData()">📈 Datos de Grafana</button>
        </div>

        <div style="text-align: center; margin: 20px 0;">
            <button class="logout-btn" onclick="logout()">🚪 Cerrar Sesión</button>
        </div>

        <div class="api-section">
            <div class="section-title">
                <h3>Datos de la API</h3>
            </div>
            <button onclick="loadProfile()">📊 Cargar Perfil</button>
            <button onclick="loadUserData()">👤 Cargar Datos de Usuario</button>
            <button onclick="loadMyGroups()">👥 Ver Mis Grupos y Roles</button>
            <button onclick="loadTokenInfo()">⏰ Info del Token</button>
            <button onclick="loadDebugToken()">🐛 Debug del Token</button>
            <div id="api-result" class="api-result">Los resultados de la API aparecerán aquí...</div>
        </div>
    </div>

    <script>
        // Datos completos de organización desde el servidor
        const groupedRoles = %s;
        const allRoles = %s;
        const orgSummary = %s;
        const debugInfo = %s;

        // Mostrar información de debug
        function displayDebugInfo() {
            const debugContent = document.getElementById('debug-content');
            if (debugInfo) {
                debugContent.innerHTML = JSON.stringify(debugInfo, null, 2);
            } else {
                debugContent.innerHTML = 'No hay información de debug disponible';
            }
            
            // También mostrar en consola para debugging
            console.log('=== DEBUG ORGANIZATION INFO ===');
            console.log('Grouped Roles:', groupedRoles);
            console.log('All Roles:', allRoles);
            console.log('Org Summary:', orgSummary);
            console.log('Debug Info:', debugInfo);
            console.log('===============================');
        }

        function toggleDebug() {
            const debugContent = document.getElementById('debug-content');
            if (debugContent.style.display === 'none') {
                debugContent.style.display = 'block';
                displayDebugInfo();
            } else {
                debugContent.style.display = 'none';
            }
        }

        // Mostrar estadísticas
        function displayStats() {
            const statsSection = document.getElementById('stats-section');
            
            // Verificar si tenemos datos de organización
            if (orgSummary && orgSummary.stats) {
                const stats = orgSummary.stats;
                let statsHTML = '';
                
                statsHTML += '<div class="stat-card">';
                statsHTML += '<div class="stat-number">' + (stats.factories_count || 0) + '</div>';
                statsHTML += '<div class="stat-label">Fábricas</div>';
                statsHTML += '</div>';
                
                statsHTML += '<div class="stat-card">';
                statsHTML += '<div class="stat-number">' + (stats.departments_count || 0) + '</div>';
                statsHTML += '<div class="stat-label">Departamentos</div>';
                statsHTML += '</div>';
                
                statsHTML += '<div class="stat-card">';
                statsHTML += '<div class="stat-number">' + (stats.total_roles || 0) + '</div>';
                statsHTML += '<div class="stat-label">Roles Totales</div>';
                statsHTML += '</div>';
                
                statsSection.innerHTML = statsHTML;
            } else if (debugInfo && debugInfo.factories_count !== undefined) {
                // Usar datos de debug si están disponibles
                let statsHTML = '';
                
                statsHTML += '<div class="stat-card">';
                statsHTML += '<div class="stat-number">' + (debugInfo.factories_count || 0) + '</div>';
                statsHTML += '<div class="stat-label">Fábricas</div>';
                statsHTML += '</div>';
                
                statsHTML += '<div class="stat-card">';
                statsHTML += '<div class="stat-number">' + (getDepartmentsCount(debugInfo) || 0) + '</div>';
                statsHTML += '<div class="stat-label">Departamentos</div>';
                statsHTML += '</div>';
                
                statsHTML += '<div class="stat-card">';
                statsHTML += '<div class="stat-number">' + (getTotalRolesCount(debugInfo) || 0) + '</div>';
                statsHTML += '<div class="stat-label">Roles Totales</div>';
                statsHTML += '</div>';
                
                statsSection.innerHTML = statsHTML;
            } else {
                statsSection.innerHTML = 
                    '<div class="stat-card"><div class="stat-number">0</div><div class="stat-label">Fábricas</div></div>' +
                    '<div class="stat-card"><div class="stat-number">0</div><div class="stat-label">Departamentos</div></div>' +
                    '<div class="stat-card"><div class="stat-number">0</div><div class="stat-label">Roles Totales</div></div>';
            }
        }

        // Helper para contar departamentos desde debugInfo
        function getDepartmentsCount(debugInfo) {
            if (!debugInfo.factory_details) return 0;
            let count = 0;
            for (const factory in debugInfo.factory_details) {
                if (debugInfo.factory_details[factory].departments_count) {
                    count += debugInfo.factory_details[factory].departments_count;
                }
            }
            return count;
        }

        // Helper para contar roles totales desde debugInfo
        function getTotalRolesCount(debugInfo) {
            if (!debugInfo.factory_details) return 0;
            let count = 0;
            for (const factory in debugInfo.factory_details) {
                const factoryDetail = debugInfo.factory_details[factory];
                if (factoryDetail.department_details) {
                    for (const dept in factoryDetail.department_details) {
                        const deptDetail = factoryDetail.department_details[dept];
                        if (deptDetail.roles_count) {
                            count += deptDetail.roles_count;
                        }
                    }
                }
            }
            return count;
        }

        // Mostrar información del usuario
        function displayUserInfo() {
            const userCookie = getCookie('user_data');
            if (userCookie) {
                try {
                    const userData = atob(userCookie);
                    const parts = userData.split('|');
                    const id = parts[0] || '';
                    const name = parts[1] || '';
                    const email = parts[2] || '';
                    
                    document.getElementById('user-display').innerHTML = 
                        '<strong>Nombre:</strong> ' + name + '<br>' +
                        '<strong>Email:</strong> ' + email + '<br>' +
                        '<strong>Authentik ID:</strong> ' + (id.length > 20 ? id.substring(0, 20) + '...' : id);
                } catch (e) {
                    console.error('Error decoding user data:', e);
                    document.getElementById('user-display').innerHTML = 'Error al cargar información del usuario';
                }
            } else {
                document.getElementById('user-display').innerHTML = 'No se encontró información del usuario';
            }
        }

        // Mostrar información de la organización MEJORADA con nuevos campos
        function displayOrganizationInfo() {
            const orgDisplay = document.getElementById('organization-display');
            
            if (groupedRoles && groupedRoles.length > 0) {
                let orgHTML = '<strong>Organización:</strong> ' + (orgSummary.name || 'Coexpan') + '<br>';
                
                // NUEVOS CAMPOS - Mostrar detalles del empleado
                orgHTML += '<div class="employee-details">';
                orgHTML += '<h4>📋 Detalles del Empleado</h4>';
                
                if (orgSummary.workday_id) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Workday ID:</span><span class="detail-value">' + orgSummary.workday_id + '</span></div>';
                }
                if (orgSummary.idn) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">IDN:</span><span class="detail-value">' + orgSummary.idn + '</span></div>';
                }
                if (orgSummary.status) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Estado:</span><span class="detail-value">' + orgSummary.status + '</span></div>';
                }
                if (orgSummary.departamento) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Departamento:</span><span class="detail-value">' + orgSummary.departamento + '</span></div>';
                }
                if (orgSummary.cargo) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Cargo:</span><span class="detail-value">' + orgSummary.cargo + '</span></div>';
                }
                if (orgSummary.tipo_empleado) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Tipo Empleado:</span><span class="detail-value">' + orgSummary.tipo_empleado + '</span></div>';
                }
                orgHTML += '</div>';
                
                // Mostrar estructura de fábricas y departamentos
                orgHTML += '<div class="factory-structure">';
                orgHTML += '<h4>🏭 Estructura de Fábricas y Roles</h4>';
                
                groupedRoles.forEach(function(factory) {
                    orgHTML += '<div class="factory-item">';
                    orgHTML += '<strong>🏭 Fábrica:</strong> ' + factory.factory;
                    orgHTML += '<div class="department-item">';
                    orgHTML += '<strong>📋 Departamento:</strong> ' + factory.department + '<br>';
                    orgHTML += '<strong>🎯 Roles:</strong> ';
                    if (factory.roles && factory.roles.length > 0) {
                        factory.roles.forEach(function(role) {
                            orgHTML += '<span class="role-badge">' + role + '</span>';
                        });
                    } else {
                        orgHTML += '<em>Sin roles asignados</em>';
                    }
                    orgHTML += '</div>';
                    orgHTML += '</div>';
                });
                orgHTML += '</div>';
                
                orgDisplay.innerHTML = orgHTML;
            } else if (debugInfo && debugInfo.organization_name === "Coexpan") {
                // CASO ESPECÍFICO: Organización Coexpan pero sin fábricas asignadas
                let orgHTML = '<strong>Organización:</strong> Coexpan<br>';
                
                // NUEVOS CAMPOS - Mostrar detalles del empleado
                orgHTML += '<div class="employee-details">';
                orgHTML += '<h4>📋 Detalles del Empleado</h4>';
                
                if (debugInfo.workday_id) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Workday ID:</span><span class="detail-value">' + debugInfo.workday_id + '</span></div>';
                }
                if (debugInfo.idn) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">IDN:</span><span class="detail-value">' + debugInfo.idn + '</span></div>';
                }
                if (debugInfo.status) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Estado:</span><span class="detail-value">' + debugInfo.status + '</span></div>';
                }
                if (debugInfo.departamento) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Departamento:</span><span class="detail-value">' + debugInfo.departamento + '</span></div>';
                }
                if (debugInfo.cargo) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Cargo:</span><span class="detail-value">' + debugInfo.cargo + '</span></div>';
                }
                if (debugInfo.tipo_empleado) {
                    orgHTML += '<div class="detail-row"><span class="detail-label">Tipo Empleado:</span><span class="detail-value">' + debugInfo.tipo_empleado + '</span></div>';
                }
                orgHTML += '</div>';
                
                if (debugInfo.factories_count > 0) {
                    orgHTML += '<div class="info-message">';
                    orgHTML += '🏢 Tienes ' + debugInfo.factories_count + ' fábrica(s) asignada(s)<br>';
                    orgHTML += '📋 Pero no tienes roles específicos en departamentos<br>';
                    orgHTML += '💡 Contacta al administrador para asignarte roles específicos';
                    orgHTML += '</div>';
                } else {
                    orgHTML += '<div class="info-message">';
                    orgHTML += '🏢 Eres parte de Coexpan<br>';
                    orgHTML += '📋 Pero no tienes fábricas específicas asignadas<br>';
                    orgHTML += '💡 Contacta al administrador para asignarte fábricas y roles';
                    orgHTML += '</div>';
                }
                
                orgDisplay.innerHTML = orgHTML;
            } else if (debugInfo && debugInfo.organization_name) {
                // Mostrar información básica si hay organización pero no roles
                let orgHTML = '<strong>Organización:</strong> ' + debugInfo.organization_name + '<br>';
                
                // NUEVOS CAMPOS
                if (debugInfo.workday_id) {
                    orgHTML += '<strong>Workday ID:</strong> ' + debugInfo.workday_id + '<br>';
                }
                if (debugInfo.idn) {
                    orgHTML += '<strong>IDN:</strong> ' + debugInfo.idn + '<br>';
                }
                
                if (debugInfo.factories_count > 0) {
                    orgHTML += '<strong>Fábricas detectadas:</strong> ' + debugInfo.factories_count + '<br>';
                    orgHTML += '<div class="warning-message">';
                    orgHTML += '⚠️ Se detectó organización pero no se pudieron extraer los roles. ';
                    orgHTML += '<button class="debug-toggle" onclick="toggleDebug()">Ver detalles de debug</button>';
                    orgHTML += '</div>';
                } else {
                    orgHTML += '<div class="warning-message">';
                    orgHTML += '⚠️ Organización detectada pero sin estructura de fábricas. ';
                    orgHTML += '<button class="debug-toggle" onclick="toggleDebug()">Ver detalles de debug</button>';
                    orgHTML += '</div>';
                }
                
                orgDisplay.innerHTML = orgHTML;
            } else {
                orgDisplay.innerHTML = 
                    '<div class="warning-message">' +
                    '❌ No se encontró información de organización específica en el token JWT.<br>' +
                    'Esto puede significar que:<br>' +
                    '• El token no incluye la información de organización<br>' +
                    '• La estructura del token es diferente a la esperada<br>' +
                    '• Hay un problema en la extracción de datos<br>' +
                    '<button class="debug-toggle" onclick="toggleDebug()">Ver detalles de debug</button>' +
                    '</div>';
            }
        }

        // Cargar grupos del usuario
        async function loadMyGroups() {
            try {
                showLoading('Cargando grupos y roles...');
                const response = await fetch('/api/my-groups', {
                    credentials: 'include'
                });
                
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#e8f5e8';
                    
                    // Actualizar display de grupos
                    displayGroups(data.groups);
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Mostrar grupos en la interfaz
        function displayGroups(groups) {
            const groupsDisplay = document.getElementById('groups-display');
            if (groups && groups.length > 0) {
                let groupsHTML = '';
                groups.forEach(function(group) {
                    groupsHTML += '<span class="group-badge">' + group + '</span>';
                });
                groupsDisplay.innerHTML = groupsHTML;
                
                // Mostrar sección admin si tiene grupos de admin
                const hasAdmin = groups.includes('authentik Admins');
                const hasGrafana = groups.includes('Grafana Admins') || groups.includes('Grafana Editors');
                
                if (hasAdmin || hasGrafana) {
                    document.getElementById('admin-section').style.display = 'block';
                }
            } else {
                groupsDisplay.innerHTML = 'No tienes grupos asignados';
            }
        }

        // Cargar estadísticas de admin
        async function loadAdminStats() {
            try {
                showLoading('Cargando estadísticas de administrador...');
                const response = await fetch('/api/admin/stats', {
                    credentials: 'include'
                });
                
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#fff3e0';
                } else if (response.status === 403) {
                    result.innerHTML = '❌ Acceso denegado: No tienes permisos de administrador';
                    result.style.background = '#f8d7da';
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Cargar datos de Grafana
        async function loadGrafanaData() {
            try {
                showLoading('Cargando datos de Grafana...');
                const response = await fetch('/api/grafana/dashboards', {
                    credentials: 'include'
                });
                
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#e8f5e8';
                } else if (response.status === 403) {
                    result.innerHTML = '❌ Acceso denegado: No tienes permisos de Grafana';
                    result.style.background = '#f8d7da';
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Helper para leer cookies
        function getCookie(name) {
            const value = '; ' + document.cookie;
            const parts = value.split('; ' + name + '=');
            if (parts.length === 2) return parts.pop().split(';').shift();
            return null;
        }

        // Cerrar sesión
        async function logout() {
            if (!confirm('¿Estás seguro de que quieres cerrar sesión?')) {
                return;
            }
            
            try {
                showLoading('Cerrando sesión...');
                await fetch('/auth/logout', {
                    method: 'POST',
                    credentials: 'include'
                });
            } catch (error) {
                console.error('Logout error:', error);
            } finally {
                // Limpiar cookies locales
                document.cookie = "session_active=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                document.cookie = "user_data=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                document.cookie = "auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                
                // Redirigir al login
                window.location.href = '/';
            }
        }

        // Cargar perfil
        async function loadProfile() {
            try {
                showLoading('Cargando perfil...');
                const response = await fetch('/api/profile', {
                    credentials: 'include'
                });
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#e8f5e8';
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Cargar datos de usuario
        async function loadUserData() {
            try {
                showLoading('Cargando datos de usuario...');
                const response = await fetch('/api/users/me', {
                    credentials: 'include'
                });
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#e8f5e8';
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Cargar información del token
        async function loadTokenInfo() {
            try {
                showLoading('Cargando información del token...');
                const response = await fetch('/api/token-info', {
                    credentials: 'include'
                });
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#e0f7fa';
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Cargar debug del token
        async function loadDebugToken() {
            try {
                showLoading('Analizando token...');
                const response = await fetch('/api/token-debug', {
                    credentials: 'include'
                });
                const result = document.getElementById('api-result');
                if (response.ok) {
                    const data = await response.json();
                    result.innerHTML = JSON.stringify(data, null, 2);
                    result.style.background = '#f3e5f5';
                } else {
                    result.innerHTML = 'Error: ' + response.status + ' - ' + response.statusText;
                    result.style.background = '#f8d7da';
                }
            } catch (error) {
                document.getElementById('api-result').innerHTML = 'Error de conexión: ' + error;
                result.style.background = '#f8d7da';
            }
        }

        // Mostrar estado de carga
        function showLoading(message) {
            const result = document.getElementById('api-result');
            result.innerHTML = '<em>' + message + '</em>';
            result.style.background = '#fff3cd';
        }

        // Verificar sesión al cargar
        function checkSession() {
            const sessionActive = getCookie('session_active');
            if (!sessionActive) {
                window.location.href = '/login';
                return;
            }
            displayUserInfo();
            displayStats();
            displayOrganizationInfo();
            displayDebugInfo(); // Mostrar debug inicialmente en consola
            loadMyGroups(); // Cargar grupos automáticamente
        }

        // Inicializar cuando se carga la página
        document.addEventListener('DOMContentLoaded', checkSession);
    </script>
</body>
</html>
`, groupedRolesJSON, allRolesJSON, orgSummaryJSON, debugJSON)
	return c.HTML(http.StatusOK, html)
}
