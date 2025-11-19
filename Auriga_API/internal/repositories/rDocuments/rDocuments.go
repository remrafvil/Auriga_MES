package rDocuments

// lo llamaremos repositories

import (
	"fmt"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rDocuments/rD_Types"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

type Repository interface {
	// CrearMultipleDocuments(db *gorm.DB, documentos []*rModels.MrDocuments) ([]*rModels.MrDocuments, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// CrearMultipleDocuments crea múltiples documentos dentro de una transacción
func CrearMultipleDocuments(db *gorm.DB, documentosData []rD_Types.DocumentsCreationData) ([]rModels.MrDocuments, error) {
	// Validar que hay documentos para crear
	if len(documentosData) == 0 {
		return []rModels.MrDocuments{}, nil
	}

	var documentosCreados []rModels.MrDocuments

	// Ejecutar en una transacción
	err := db.Transaction(func(tx *gorm.DB) error {
		// Convertir DocumentsCreationData a MrDocuments
		documentos := make([]rModels.MrDocuments, len(documentosData))
		now := time.Now()

		for i, data := range documentosData {
			documentos[i] = rModels.MrDocuments{
				Nombre:      data.Nombre,
				Descripcion: data.Descripcion,
				URL:         data.URL,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
		}

		// Crear los documentos en la base de datos
		if err := tx.Create(&documentos).Error; err != nil {
			return err
		}

		documentosCreados = documentos
		return nil
	})

	if err != nil {
		return nil, err
	}

	return documentosCreados, nil
}

// EntityConstraint interface constraint para tipos que pueden tener documentos
type EntityConstraint interface {
	rModels.MrProduct | rModels.MrAsset
}

// ==================== FUNCIONES GENÉRICAS PRINCIPALES ====================

// BuscarYAsignarDocumentos busca la entidad por campo/valor y asigna documentos existentes
func BuscarYAsignarDocumentos[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
	nombresDocumentos []string,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var entity T

		// Buscar la entidad
		if err := tx.Where(fieldName+" = ?", fieldValue).First(&entity).Error; err != nil {
			return fmt.Errorf("error buscando entidad: %w", err)
		}

		// Buscar documentos
		var documentos []rModels.MrDocuments
		if err := tx.Where("nombre IN ?", nombresDocumentos).Find(&documentos).Error; err != nil {
			return fmt.Errorf("error buscando documentos: %w", err)
		}

		// Verificar que se encontraron todos los documentos
		if len(documentos) != len(nombresDocumentos) {
			return fmt.Errorf("no se encontraron todos los documentos solicitados")
		}

		// Asignar documentos
		if err := tx.Model(&entity).Association("Documents").Append(&documentos); err != nil {
			return fmt.Errorf("error asignando documentos: %w", err)
		}

		return nil
	})
}

// CrearYAsignarDocumentos busca la entidad y crea/asigna nuevos documentos
func CrearYAsignarDocumentos[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
	documentosData []rD_Types.DocumentsCreationData,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var entity T

		// Buscar la entidad
		if err := tx.Where(fieldName+" = ?", fieldValue).First(&entity).Error; err != nil {
			return fmt.Errorf("error buscando entidad: %w", err)
		}

		// Crear los documentos
		documentosCreados, err := crearDocumentosEnTx(tx, documentosData)
		if err != nil {
			return fmt.Errorf("error creando documentos: %w", err)
		}

		// Asignar documentos
		if err := tx.Model(&entity).Association("Documents").Append(&documentosCreados); err != nil {
			return fmt.Errorf("error asignando documentos: %w", err)
		}

		return nil
	})
}

// ReemplazarDocumentos busca la entidad y reemplaza todos sus documentos
func ReemplazarDocumentos[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
	documentosData []rD_Types.DocumentsCreationData,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var entity T

		// Buscar la entidad
		if err := tx.Where(fieldName+" = ?", fieldValue).First(&entity).Error; err != nil {
			return fmt.Errorf("error buscando entidad: %w", err)
		}

		// Limpiar documentos existentes
		if err := tx.Model(&entity).Association("Documents").Clear(); err != nil {
			return fmt.Errorf("error limpiando documentos existentes: %w", err)
		}

		// Crear y asignar nuevos documentos si se proporcionan
		if len(documentosData) > 0 {
			documentosCreados, err := crearDocumentosEnTx(tx, documentosData)
			if err != nil {
				return fmt.Errorf("error creando documentos: %w", err)
			}

			if err := tx.Model(&entity).Association("Documents").Append(&documentosCreados); err != nil {
				return fmt.Errorf("error asignando nuevos documentos: %w", err)
			}
		}

		return nil
	})
}

// AsignarDocumentosConVerificacion asigna documentos y retorna los no encontrados
func AsignarDocumentosConVerificacion[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
	nombresDocumentos []string,
) ([]string, error) {
	var documentosNoEncontrados []string

	err := db.Transaction(func(tx *gorm.DB) error {
		var entity T

		// Buscar la entidad
		if err := tx.Where(fieldName+" = ?", fieldValue).First(&entity).Error; err != nil {
			return fmt.Errorf("error buscando entidad: %w", err)
		}

		// Buscar documentos existentes
		var documentos []rModels.MrDocuments
		if err := tx.Where("nombre IN ?", nombresDocumentos).Find(&documentos).Error; err != nil {
			return fmt.Errorf("error buscando documentos: %w", err)
		}

		// Verificar qué documentos no se encontraron
		documentosEncontrados := make(map[string]bool)
		for _, doc := range documentos {
			documentosEncontrados[doc.Nombre] = true
		}

		documentosNoEncontrados = []string{}
		for _, nombreDoc := range nombresDocumentos {
			if !documentosEncontrados[nombreDoc] {
				documentosNoEncontrados = append(documentosNoEncontrados, nombreDoc)
			}
		}

		// Si no se encontraron todos los documentos, retornar error
		if len(documentosNoEncontrados) > 0 {
			return fmt.Errorf("documentos no encontrados: %v", documentosNoEncontrados)
		}

		// Asignar documentos
		if err := tx.Model(&entity).Association("Documents").Append(&documentos); err != nil {
			return fmt.Errorf("error asignando documentos: %w", err)
		}

		return nil
	})

	return documentosNoEncontrados, err
}

// AsignarDocumentosAMultiplesEntidades asigna documentos a múltiples entidades
func AsignarDocumentosAMultiplesEntidades[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValues []interface{},
	nombresDocumentos []string,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var entities []T

		// Buscar las entidades
		if err := tx.Where(fieldName+" IN ?", fieldValues).Find(&entities).Error; err != nil {
			return fmt.Errorf("error buscando entidades: %w", err)
		}

		// Verificar que se encontraron todas las entidades
		if len(entities) != len(fieldValues) {
			return fmt.Errorf("no se encontraron todas las entidades solicitadas")
		}

		// Buscar los documentos
		var documentos []rModels.MrDocuments
		if err := tx.Where("nombre IN ?", nombresDocumentos).Find(&documentos).Error; err != nil {
			return fmt.Errorf("error buscando documentos: %w", err)
		}

		// Asignar documentos a cada entidad
		for i := range entities {
			if err := tx.Model(&entities[i]).Association("Documents").Append(&documentos); err != nil {
				return fmt.Errorf("error asignando documentos a entidad %d: %w", i, err)
			}
		}

		return nil
	})
}

// ==================== FUNCIONES ESPECÍFICAS WRAPPER ====================

// Para Productos
func AsignarDocumentosAProductoPorNombre(db *gorm.DB, nombre string, documentos []string) error {
	return BuscarYAsignarDocumentos[rModels.MrProduct](db, "name", nombre, documentos)
}

func AsignarDocumentosAProductoPorID(db *gorm.DB, id uint, documentos []string) error {
	return BuscarYAsignarDocumentos[rModels.MrProduct](db, "id", id, documentos)
}

func CrearYAsignarDocumentosAProductoPorNombre(db *gorm.DB, nombre string, documentosData []rD_Types.DocumentsCreationData) error {
	return CrearYAsignarDocumentos[rModels.MrProduct](db, "name", nombre, documentosData)
}

func ReemplazarDocumentosDeProductoPorNombre(db *gorm.DB, nombre string, documentosData []rD_Types.DocumentsCreationData) error {
	return ReemplazarDocumentos[rModels.MrProduct](db, "name", nombre, documentosData)
}

// Para Activos
func AsignarDocumentosAActivoPorTechCode(db *gorm.DB, techCode uint, documentos []string) error {
	return BuscarYAsignarDocumentos[rModels.MrAsset](db, "tech_code", techCode, documentos)
}

func AsignarDocumentosAActivoPorID(db *gorm.DB, id uint, documentos []string) error {
	return BuscarYAsignarDocumentos[rModels.MrAsset](db, "id", id, documentos)
}

func CrearYAsignarDocumentosAActivoPorTechCode(db *gorm.DB, techCode uint, documentosData []rD_Types.DocumentsCreationData) error {
	return CrearYAsignarDocumentos[rModels.MrAsset](db, "tech_code", techCode, documentosData)
}

func ReemplazarDocumentosDeActivoPorTechCode(db *gorm.DB, techCode uint, documentosData []rD_Types.DocumentsCreationData) error {
	return ReemplazarDocumentos[rModels.MrAsset](db, "tech_code", techCode, documentosData)
}

// ==================== FUNCIÓN AUXILIAR ====================

// crearDocumentosEnTx crea documentos dentro de una transacción
func crearDocumentosEnTx(tx *gorm.DB, documentosData []rD_Types.DocumentsCreationData) ([]rModels.MrDocuments, error) {
	if len(documentosData) == 0 {
		return []rModels.MrDocuments{}, nil
	}

	documentos := make([]rModels.MrDocuments, len(documentosData))
	now := time.Now()

	for i, data := range documentosData {
		documentos[i] = rModels.MrDocuments{
			Nombre:      data.Nombre,
			Descripcion: data.Descripcion,
			URL:         data.URL,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	if err := tx.Create(&documentos).Error; err != nil {
		return nil, err
	}

	return documentos, nil
}

// ==================== FUNCIONES DE CONSULTA GENÉRICAS ====================

// ObtenerDocumentosDeEntidad obtiene todos los documentos de una entidad
func ObtenerDocumentosDeEntidad[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
) ([]rModels.MrDocuments, error) {
	var entity T
	var documentos []rModels.MrDocuments

	// Buscar la entidad con sus documentos
	if err := db.Where(fieldName+" = ?", fieldValue).Preload("Documents").First(&entity).Error; err != nil {
		return nil, fmt.Errorf("error buscando entidad: %w", err)
	}

	// Obtener los documentos a través de la asociación
	if err := db.Model(&entity).Association("Documents").Find(&documentos); err != nil {
		return nil, fmt.Errorf("error obteniendo documentos: %w", err)
	}

	return documentos, nil
}

// ContarDocumentosDeEntidad cuenta los documentos de una entidad
func ContarDocumentosDeEntidad[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
) (int64, error) {
	var entity T

	// Buscar la entidad
	if err := db.Where(fieldName+" = ?", fieldValue).First(&entity).Error; err != nil {
		return 0, fmt.Errorf("error buscando entidad: %w", err)
	}

	// Contar documentos - Count() no retorna error, retorna el count directamente
	count := db.Model(&entity).Association("Documents").Count()

	return count, nil
}

// RemoverDocumentosDeEntidad remueve documentos específicos de una entidad
func RemoverDocumentosDeEntidad[T EntityConstraint](
	db *gorm.DB,
	fieldName string,
	fieldValue interface{},
	nombresDocumentos []string,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var entity T

		// Buscar la entidad
		if err := tx.Where(fieldName+" = ?", fieldValue).First(&entity).Error; err != nil {
			return fmt.Errorf("error buscando entidad: %w", err)
		}

		// Buscar los documentos a remover
		var documentos []rModels.MrDocuments
		if err := tx.Where("nombre IN ?", nombresDocumentos).Find(&documentos).Error; err != nil {
			return fmt.Errorf("error buscando documentos: %w", err)
		}

		// Remover documentos
		if err := tx.Model(&entity).Association("Documents").Delete(&documentos); err != nil {
			return fmt.Errorf("error removiendo documentos: %w", err)
		}

		return nil
	})
}
