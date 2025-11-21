package hAuth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Función helper segura para obtener claims del contexto - MEJORADA CON DEBUG
func getUserClaimsFromContext(c echo.Context) (map[string]interface{}, error) {
	claimsInterface := c.Get("user_claims")
	if claimsInterface == nil {
		return nil, fmt.Errorf("no user claims found in context")
	}

	// IMPORTANTE: jwt.MapClaims es ya un map[string]interface{}, no necesitamos convertirlo
	if claimsMap, ok := claimsInterface.(jwt.MapClaims); ok {
		return map[string]interface{}(claimsMap), nil
	}

	// Fallback: intentar como map[string]interface{} regular
	if claimsMap, ok := claimsInterface.(map[string]interface{}); ok {
		return claimsMap, nil
	}

	return nil, fmt.Errorf("invalid claims type in context, got: %T", claimsInterface)
}

// Función helper mejorada para obtener info del usuario + organización - CON NUEVOS CAMPOS
func getUserInfoFromContext(c echo.Context) (userID, email, name string, groups []string, organization map[string]interface{}) {
	claims, err := getUserClaimsFromContext(c)
	if err != nil {
		return "", "", "", nil, nil
	}

	// Extraer información básica
	userID, _ = claims["sub"].(string)
	email, _ = claims["email"].(string)
	name, _ = claims["name"].(string)

	// ✅ Extraer grupos del token
	if groupsInterface, ok := claims["groups"].([]interface{}); ok {
		for _, group := range groupsInterface {
			if groupStr, ok := group.(string); ok {
				groups = append(groups, groupStr)
			}
		}
	}

	// ✅ Extraer información completa de la organización - CON MÁS ROBUSTEZ
	if orgInterface, exists := claims["organization"]; exists {
		if orgInterface == nil {
			return userID, email, name, groups, nil
		}

		// Intentar diferentes formas de castear la organización
		if orgMap, ok := orgInterface.(map[string]interface{}); ok {
			organization = orgMap
		} else if orgMapInterface, ok := orgInterface.(map[interface{}]interface{}); ok {
			organization = convertMap(orgMapInterface)
		}
	}

	return
}

// Función helper para convertir map[interface{}]interface{} a map[string]interface{}
func convertMap(original map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range original {
		if strKey, ok := key.(string); ok {
			// Convertir valores recursivamente si es necesario
			if nestedMap, ok := value.(map[interface{}]interface{}); ok {
				result[strKey] = convertMap(nestedMap)
			} else if nestedSlice, ok := value.([]interface{}); ok {
				// Convertir slices si es necesario
				result[strKey] = convertSlice(nestedSlice)
			} else {
				result[strKey] = value
			}
		}
	}
	return result
}

// Función helper para convertir slices
func convertSlice(original []interface{}) []interface{} {
	result := make([]interface{}, len(original))
	for i, item := range original {
		if nestedMap, ok := item.(map[interface{}]interface{}); ok {
			result[i] = convertMap(nestedMap)
		} else {
			result[i] = item
		}
	}
	return result
}

// Función helper para obtener la estructura completa de fábricas
func getUserFactories(organization map[string]interface{}) map[string]interface{} {
	if organization == nil {
		return nil
	}

	if factories, ok := organization["factories"].(map[string]interface{}); ok {
		return factories
	}

	return nil
}

// Función helper para obtener el nombre de la organización
func getOrganizationName(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	if name, ok := organization["name"].(string); ok {
		return name
	}

	return ""
}

// NUEVAS FUNCIONES PARA LOS CAMPOS ADICIONALES

// Función helper para obtener workday_id
func getWorkdayID(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	if workdayID, ok := organization["workday_id"].(string); ok {
		return workdayID
	}

	return ""
}

// Función helper para obtener id_card
func getIDCard(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	// Ahora el campo se llama "idn" en el token
	if idn, ok := organization["idn"].(string); ok {
		return idn
	}

	return ""
}

// Función helper para obtener status
func getEmployeeStatus(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	if status, ok := organization["status"].(string); ok {
		return status
	}

	return ""
}

// Función helper para obtener departamento
func getDepartment(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	if dept, ok := organization["departamento"].(string); ok {
		return dept
	}

	return ""
}

// Función helper para obtener cargo
func getPosition(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	if position, ok := organization["cargo"].(string); ok {
		return position
	}

	return ""
}

// Función helper para obtener tipo_empleado
func getEmployeeType(organization map[string]interface{}) string {
	if organization == nil {
		return ""
	}

	if empType, ok := organization["tipo_empleado"].(string); ok {
		return empType
	}

	return ""
}

// Función helper para obtener TODOS los roles del usuario en formato plano
func getAllUserRoles(organization map[string]interface{}) []map[string]interface{} {
	var allRoles []map[string]interface{}

	factories := getUserFactories(organization)
	if factories == nil {
		return allRoles
	}

	for factoryName, factoryData := range factories {
		if factoryMap, ok := factoryData.(map[string]interface{}); ok {
			if deptData, ok := factoryMap["departments"].(map[string]interface{}); ok {
				for deptName, deptInfo := range deptData {
					if deptMap, ok := deptInfo.(map[string]interface{}); ok {
						if rolesInterface, ok := deptMap["roles"].([]interface{}); ok {
							for _, role := range rolesInterface {
								if roleStr, ok := role.(string); ok {
									allRoles = append(allRoles, map[string]interface{}{
										"factory":    factoryName,
										"department": deptName,
										"role":       roleStr,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return allRoles
}

// Función helper para obtener roles agrupados por factory y department
func getUserRolesGrouped(organization map[string]interface{}) []map[string]interface{} {
	var groupedRoles []map[string]interface{}

	if organization == nil {
		return groupedRoles
	}

	factories := getUserFactories(organization)
	if factories == nil {
		return groupedRoles
	}

	for factoryName, factoryData := range factories {
		if factoryMap, ok := factoryData.(map[string]interface{}); ok {
			if deptData, ok := factoryMap["departments"].(map[string]interface{}); ok {
				for deptName, deptInfo := range deptData {
					if deptMap, ok := deptInfo.(map[string]interface{}); ok {
						if rolesInterface, ok := deptMap["roles"].([]interface{}); ok {
							var roles []string
							for _, role := range rolesInterface {
								if roleStr, ok := role.(string); ok {
									roles = append(roles, roleStr)
								}
							}

							if len(roles) > 0 {
								groupedRoles = append(groupedRoles, map[string]interface{}{
									"factory":    factoryName,
									"department": deptName,
									"roles":      roles,
								})
							}
						}
					}
				}
			}
		}
	}

	return groupedRoles
}

// Función helper para obtener roles específicos del usuario por factory y department
func getUserRoles(organization map[string]interface{}, factory, department string) []string {
	var roles []string

	factories := getUserFactories(organization)
	if factories == nil {
		return roles
	}

	if factoryData, ok := factories[factory].(map[string]interface{}); ok {
		if deptData, ok := factoryData["departments"].(map[string]interface{}); ok {
			if deptInfo, ok := deptData[department].(map[string]interface{}); ok {
				if rolesInterface, ok := deptInfo["roles"].([]interface{}); ok {
					for _, role := range rolesInterface {
						if roleStr, ok := role.(string); ok {
							roles = append(roles, roleStr)
						}
					}
				}
			}
		}
	}

	return roles
}

// Función helper para obtener todas las fábricas a las que tiene acceso el usuario
func getUserFactoryNames(organization map[string]interface{}) []string {
	var factoryNames []string

	factories := getUserFactories(organization)
	if factories == nil {
		return factoryNames
	}

	for factoryName := range factories {
		factoryNames = append(factoryNames, factoryName)
	}

	return factoryNames
}

// Función helper para obtener todos los departamentos de una fábrica específica
func getDepartmentsForFactory(organization map[string]interface{}, factory string) []string {
	var departments []string

	factories := getUserFactories(organization)
	if factories == nil {
		return departments
	}

	if factoryData, ok := factories[factory].(map[string]interface{}); ok {
		if deptData, ok := factoryData["departments"].(map[string]interface{}); ok {
			for deptName := range deptData {
				departments = append(departments, deptName)
			}
		}
	}

	return departments
}

// Función helper para verificar si el usuario tiene un rol específico
func hasOrganizationRole(organization map[string]interface{}, factory, department, role string) bool {
	roles := getUserRoles(organization, factory, department)
	for _, userRole := range roles {
		if userRole == role {
			return true
		}
	}
	return false
}

// Función helper para verificar si el usuario tiene acceso a una fábrica
func hasFactoryAccess(organization map[string]interface{}, factory string) bool {
	factories := getUserFactoryNames(organization)
	for _, userFactory := range factories {
		if userFactory == factory {
			return true
		}
	}
	return false
}

// Función para obtener resumen completo de la organización del usuario - ACTUALIZADA
func getOrganizationSummary(organization map[string]interface{}) map[string]interface{} {
	summary := make(map[string]interface{})

	if organization == nil {
		return summary
	}

	// Información básica
	summary["name"] = getOrganizationName(organization)

	// CAMPOS ACTUALIZADOS - usar idn en lugar de id_card
	summary["workday_id"] = getWorkdayID(organization)
	summary["idn"] = getIDCard(organization) // Ahora getIDCard retorna el campo "idn"
	summary["status"] = getEmployeeStatus(organization)
	summary["departamento"] = getDepartment(organization)
	summary["cargo"] = getPosition(organization)
	summary["tipo_empleado"] = getEmployeeType(organization)

	// Fábricas y acceso
	factories := getUserFactories(organization)
	summary["factories_count"] = len(factories)
	summary["factory_names"] = getUserFactoryNames(organization)

	// Roles y permisos
	allRoles := getAllUserRoles(organization)
	summary["total_roles"] = len(allRoles)
	summary["roles_detail"] = allRoles

	// Agrupado por factory/department
	summary["grouped_roles"] = getUserRolesGrouped(organization)

	// Estructura completa (opcional, para debugging)
	summary["full_structure"] = organization

	return summary
}

// Función alternativa más robusta para debug - ACTUALIZADA
func debugOrganizationStructure(organization map[string]interface{}) map[string]interface{} {
	debugInfo := make(map[string]interface{})

	if organization == nil {
		debugInfo["error"] = "Organization is nil"
		return debugInfo
	}

	debugInfo["raw_organization"] = organization

	// Verificar estructura nivel por nivel
	if name, ok := organization["name"].(string); ok {
		debugInfo["organization_name"] = name
	} else {
		debugInfo["organization_name_error"] = "No name found or wrong type"
	}

	// CAMPOS ACTUALIZADOS - usar idn
	debugInfo["workday_id"] = getWorkdayID(organization)
	debugInfo["idn"] = getIDCard(organization) // Ahora usa idn
	debugInfo["status"] = getEmployeeStatus(organization)
	debugInfo["departamento"] = getDepartment(organization)
	debugInfo["cargo"] = getPosition(organization)
	debugInfo["tipo_empleado"] = getEmployeeType(organization)

	if factories, ok := organization["factories"].(map[string]interface{}); ok {
		debugInfo["factories_count"] = len(factories)
		debugInfo["factory_names"] = getKeys(factories)

		factoryDetails := make(map[string]interface{})
		for factoryName, factoryData := range factories {
			factoryDetail := make(map[string]interface{})

			if factoryMap, ok := factoryData.(map[string]interface{}); ok {
				if deptData, ok := factoryMap["departments"].(map[string]interface{}); ok {
					factoryDetail["departments_count"] = len(deptData)
					factoryDetail["department_names"] = getKeys(deptData)

					departmentDetails := make(map[string]interface{})
					for deptName, deptInfo := range deptData {
						deptDetail := make(map[string]interface{})

						if deptMap, ok := deptInfo.(map[string]interface{}); ok {
							if rolesInterface, ok := deptMap["roles"].([]interface{}); ok {
								deptDetail["roles_count"] = len(rolesInterface)
								deptDetail["roles"] = rolesInterface
							} else {
								deptDetail["roles_error"] = "No roles found or wrong type"
							}
						} else {
							deptDetail["department_error"] = "Department data is not a map"
						}
						departmentDetails[deptName] = deptDetail
					}
					factoryDetail["department_details"] = departmentDetails
				} else {
					factoryDetail["departments_error"] = "No departments found or wrong type"
				}
			} else {
				factoryDetail["factory_error"] = "Factory data is not a map"
			}
			factoryDetails[factoryName] = factoryDetail
		}
		debugInfo["factory_details"] = factoryDetails
	} else {
		debugInfo["factories_error"] = "No factories found or wrong type"
	}

	return debugInfo
}

// Función helper para obtener keys de un map
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Función helper para verificar si el usuario tiene fábricas asignadas
func hasFactoryAssignments(organization map[string]interface{}) bool {
	if organization == nil {
		return false
	}

	factories := getUserFactories(organization)
	if factories == nil {
		return false
	}

	// Verificar si hay al menos una fábrica con departamentos y roles
	for _, factoryData := range factories {
		if factoryMap, ok := factoryData.(map[string]interface{}); ok {
			if deptData, ok := factoryMap["departments"].(map[string]interface{}); ok && len(deptData) > 0 {
				// Verificar si hay al menos un departamento con roles
				for _, deptInfo := range deptData {
					if deptMap, ok := deptInfo.(map[string]interface{}); ok {
						if rolesInterface, ok := deptMap["roles"].([]interface{}); ok && len(rolesInterface) > 0 {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

// Función helper para obtener mensaje de estado de la organización
func getOrganizationStatusMessage(organization map[string]interface{}) string {
	if organization == nil {
		return "No se detectó información de organización"
	}

	factories := getUserFactories(organization)
	if factories == nil || len(factories) == 0 {
		return "Eres parte de Coexpan pero no tienes fábricas asignadas. Contacta al administrador."
	}

	allRoles := getAllUserRoles(organization)
	if len(allRoles) == 0 {
		return "Tienes fábricas asignadas pero no roles específicos. Contacta al administrador."
	}

	return "Tienes acceso completo a la organización"
}

func (h *handler) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *handler) loginHandler(c echo.Context) error {
	// Verificar si ya está autenticado mirando la cookie
	cookie, err := c.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		// Validar el token
		token, err := h.service.ValidateToken(cookie.Value)
		if err == nil && token.Valid {
			// Si ya tiene token válido, redirigir al frontend
			frontendURL := "http://192.168.122.211:5826/dashboard"
			return c.Redirect(http.StatusFound, frontendURL)
		}
	}

	state, err := generateState()
	if err != nil {
		h.logger.Error("Failed to generate state", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}

	authURL := h.service.GetAuthURL(state)
	return c.Redirect(http.StatusFound, authURL)
}

func (h *handler) authCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		h.logger.Warn("OAuth callback missing code parameter")
		return echo.NewHTTPError(http.StatusBadRequest, "Missing code parameter")
	}

	state := c.QueryParam("state")
	if state == "" {
		h.logger.Warn("OAuth callback missing state parameter")
		return echo.NewHTTPError(http.StatusBadRequest, "Missing state parameter")
	}

	// Intercambiar código por tokens
	token, err := h.service.ExchangeCode(c.Request().Context(), code)
	if err != nil {
		h.logger.Error("Failed to exchange code for token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Failed to exchange code: "+err.Error())
	}

	// Obtener información del usuario
	userInfo, err := h.service.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		h.logger.Error("Failed to get user info", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Failed to get user info: "+err.Error())
	}

	// Sincronizar empleado con la base de datos
	employee, err := h.service.SyncUser(userInfo)
	if err != nil {
		h.logger.Error("Failed to sync employee",
			zap.Error(err),
			zap.String("user_id", userInfo.Sub))
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Failed to sync employee: "+err.Error())
	}

	// Calcular expiración de la cookie
	tokenExpiry := token.Expiry
	maxAge := int(tokenExpiry.Sub(time.Now()).Seconds())
	if maxAge > 24*3600 {
		maxAge = 24 * 3600
	}

	// ✅ CONFIGURACIÓN CORREGIDA DE COOKIES - Sin dominio específico
	// 1. Cookie principal con el token JWT (HTTP-only para seguridad)
	authCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // ✅ false para desarrollo HTTP
		SameSite: http.SameSiteLaxMode, // ✅ Cambiado a Lax
		MaxAge:   maxAge,
		Expires:  tokenExpiry,
	}
	c.SetCookie(authCookie)

	// 2. Cookie con información básica del empleado (accesible desde React)
	employeeName := employee.FirstName + " " + employee.LastName
	userData := map[string]interface{}{
		"id":    employee.AuthentikID,
		"name":  employeeName,
		"email": employee.Email,
	}

	userDataJSON, _ := json.Marshal(userData)
	employeeCookie := &http.Cookie{
		Name:     "user_data",
		Value:    base64.StdEncoding.EncodeToString(userDataJSON),
		Path:     "/",
		HttpOnly: false, // ✅ Accesible desde JavaScript
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // ✅ Cambiado a Lax
		MaxAge:   maxAge,
		Expires:  tokenExpiry,
	}
	c.SetCookie(employeeCookie)

	// 3. Cookie de sesión activa
	sessionCookie := &http.Cookie{
		Name:     "session_active",
		Value:    "true",
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // ✅ Cambiado a Lax
		MaxAge:   maxAge,
		Expires:  tokenExpiry,
	}
	c.SetCookie(sessionCookie)

	h.logger.Info("User authenticated and cookies set",
		zap.String("authentik_id", userInfo.Sub),
		zap.Uint("employee_id", employee.ID),
		zap.String("email", employee.Email))

	// ✅ REDIRIGIR AL FRONTEND REACT
	frontendURL := "http://192.168.122.211:5826/dashboard"
	return c.Redirect(http.StatusFound, frontendURL)
}

func (h *handler) homeHandler(c echo.Context) error {
	// Verificar si hay cookie de autenticación
	cookie, err := c.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return c.Redirect(http.StatusFound, "/dashboard")
	}

	return h.loginPageHandler(c)
}

// Handler para verificar autenticación desde React
func (h *handler) authCheckHandler(c echo.Context) error {
	// Verificar si hay token en cookie
	cookie, err := c.Cookie("auth_token")
	if err != nil || cookie.Value == "" {
		h.logger.Debug("No auth_token cookie found",
			zap.Error(err),
			zap.String("path", c.Path()))
		return c.JSON(http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "No authentication token found",
		})
	}

	// Validar el token
	token, err := h.service.ValidateToken(cookie.Value)
	if err != nil || !token.Valid {
		h.logger.Debug("Invalid token",
			zap.Error(err),
			zap.String("path", c.Path()))
		return c.JSON(http.StatusOK, map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid token: " + err.Error(),
		})
	}

	// Obtener información del usuario desde el contexto
	userID, userEmail, userName, userGroups, organization := getUserInfoFromContext(c)

	h.logger.Debug("User authenticated successfully",
		zap.String("user_id", userID),
		zap.String("path", c.Path()))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"id":     userID,
			"email":  userEmail,
			"name":   userName,
			"groups": userGroups,
		},
		"organization": getOrganizationSummary(organization),
	})
}
