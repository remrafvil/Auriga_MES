package rAssets

import (
	"context"
	"errors"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	AssetList(ctx context.Context) []rModels.MrAsset
	AssetListByParentCodes(parentCodes []string, ctx context.Context) ([]rModels.MrAsset, error)
	AssetInfo(id uint) (rModels.MrAsset, error)
	AssetAdd(productID uint, parentId *uint, location string, techCode uint, code string, sn string, SapCode string, hierarchicalLevel []string) (rModels.MrAsset, error)
	AssetUpdate(id uint, productID uint, parentId *uint, location string, techCode uint, code string, sn string, SapCode string, hierarchicalLevel []string) (rModels.MrAsset, error)
	AssetDel(id uint) error
	AssetByFactLine(hierarchicalLevel_1 string, hierarchicalLevel_2 string) (rModels.MrAsset, error)

	GetAssetWithProductDetails(assetID uint) (rModels.MrAsset, error)
	GetDosingSystemByLine(ctx context.Context, lineID uint) ([]rModels.MrAsset, error)
	GetDoserComponents(ctx context.Context, doserID uint) ([]rModels.MrAsset, error)
}

type repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

/* func (m *repository) AssetList(ctx context.Context) []rModels.MrAsset { // Lista todos
	var data = []rModels.MrAsset{}
	log.Println("*****************      Llego por aquí ASSET LIST  *****************")
	err := m.db.Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}
	return data
} */

func (m *repository) AssetList(ctx context.Context) []rModels.MrAsset {
	// Obtener información del contexto para logging
	userID, _ := ctx.Value("user_id").(string)
	userName, _ := ctx.Value("user_name").(string)
	factoryNames, _ := ctx.Value("factory_names").([]string)

	m.logger.Debug("Executing database query: AssetList",
		zap.String("repository", "AssetList"),
		zap.String("user_id", userID),
		zap.String("user_name", userName),
		zap.Strings("factory_names", factoryNames))

	var data []rModels.MrAsset
	start := time.Now()

	err := m.db.Find(&data).Error
	if err != nil {
		m.logger.Error("Database query failed",
			zap.Error(err),
			zap.String("repository", "AssetList"),
			zap.String("query", "SELECT * FROM mr_assets"),
			zap.String("user_id", userID),
			zap.String("user_name", userName))
		return []rModels.MrAsset{}
	}

	duration := time.Since(start)
	m.logger.Debug("Database query completed",
		zap.String("repository", "AssetList"),
		zap.Int("recordsFound", len(data)),
		zap.Duration("queryDuration", duration),
		zap.String("user_id", userID),
		zap.String("user_name", userName))

	return data
}

// ✅ NUEVO MÉTODO: AssetListByParentCodes para múltiples códigos padre
func (m *repository) AssetListByParentCodes(parentCodes []string, ctx context.Context) ([]rModels.MrAsset, error) {
	if len(parentCodes) == 0 {
		return []rModels.MrAsset{}, nil
	}

	// Obtener información del contexto para logging
	userID, _ := ctx.Value("user_id").(string)
	userName, _ := ctx.Value("user_name").(string)

	m.logger.Debug("Executing database query: AssetListByParentCodes",
		zap.String("repository", "AssetListByParentCodes"),
		zap.Strings("parentCodes", parentCodes),
		zap.String("user_id", userID),
		zap.String("user_name", userName))

	// 1. Primero obtenemos los IDs de todos los padres
	var parentAssets []rModels.MrAsset
	err := m.db.Where("code IN ?", parentCodes).Find(&parentAssets).Error
	if err != nil {
		m.logger.Error("Failed to find parent assets",
			zap.Error(err),
			zap.Strings("parentCodes", parentCodes),
			zap.String("user_id", userID))
		return nil, err
	}

	if len(parentAssets) == 0 {
		m.logger.Warn("No parent assets found for codes",
			zap.Strings("parentCodes", parentCodes),
			zap.String("user_id", userID))
		return []rModels.MrAsset{}, nil
	}

	// 2. Extraer los IDs de los padres encontrados
	parentIDs := make([]uint, len(parentAssets))
	for i, asset := range parentAssets {
		parentIDs[i] = asset.ID
	}

	// 3. Consulta recursiva para obtener todos los hijos de múltiples padres
	query := `
        WITH RECURSIVE asset_tree AS (
            -- Caso base: los activos padres
            SELECT * 
            FROM mr_assets 
            WHERE id IN (?)
            
            UNION ALL
            
            -- Caso recursivo: todos los hijos
            SELECT a.* 
            FROM mr_assets a
            INNER JOIN asset_tree at ON a.parent_id = at.id
        )
        SELECT * FROM asset_tree
        ORDER BY hierarchical_level, code
    `

	var data []rModels.MrAsset
	start := time.Now()

	// Ejecutar la consulta con los IDs de los padres
	err = m.db.Raw(query, parentIDs).Scan(&data).Error
	if err != nil {
		m.logger.Error("Database query failed",
			zap.Error(err),
			zap.String("repository", "AssetListByParentCodes"),
			zap.Strings("parentCodes", parentCodes),
			zap.Any("parentIDs", parentIDs),
			zap.String("user_id", userID))
		return nil, err
	}

	duration := time.Since(start)
	m.logger.Debug("Database query completed",
		zap.String("repository", "AssetListByParentCodes"),
		zap.Strings("parentCodes", parentCodes),
		zap.Any("parentIDs", parentIDs),
		zap.Int("recordsFound", len(data)),
		zap.Duration("queryDuration", duration),
		zap.String("user_id", userID))

	return data, nil
}

func (m *repository) AssetInfo(id uint) (rModels.MrAsset, error) { // Encuentra uno
	var data rModels.MrAsset

	if err := m.db.Where("id  = ? ", id).Find(&data).Error; err != nil {
		return rModels.MrAsset{}, errors.New("no existe registro con es id")
	}

	return data, nil
}

func (m *repository) AssetAdd(productID uint, parentId *uint, location string, techCode uint, code string, sn string, SapCode string, hierarchicalLevel []string) (rModels.MrAsset, error) {
	var data rModels.MrAsset

	data.ProductID = productID
	data.ParentID = parentId
	data.Location = location
	data.TechCode = techCode
	data.Code = code
	data.Sn = sn
	data.SapCode = SapCode
	data.HierarchicalLevel = hierarchicalLevel

	if err := m.db.Create(&data).Error; err != nil {
		return rModels.MrAsset{}, err
	} else {

		return data, nil
	}
}

func (m *repository) AssetUpdate(id uint, productID uint, parentId *uint, location string, techCode uint, code string, sn string, SapCode string, hierarchicalLevel []string) (rModels.MrAsset, error) {
	var data rModels.MrAsset

	if err := m.db.Where("id = ? ", id).First(&data).Error; err != nil {
		return rModels.MrAsset{}, errors.New("no hay registro con ese id")
	}
	data.ProductID = productID
	data.ParentID = parentId
	data.Location = location
	data.TechCode = techCode
	data.Code = code
	data.Sn = sn
	data.SapCode = SapCode
	data.HierarchicalLevel = hierarchicalLevel

	if err := m.db.Save(&data).Error; err != nil {
		return rModels.MrAsset{}, errors.New("no se pudo actualizar")
	}
	return data, nil

}

func (m *repository) AssetDel(id uint) error { // Borra uno
	var data rModels.MrAsset
	if err := m.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

func (m *repository) AssetByFactLine(hierarchicalLevel_1 string, hierarchicalLevel_2 string) (rModels.MrAsset, error) { // Encuentra uno
	var data rModels.MrAsset
	// if err := m.db.Where("hierarchical_Level[1] = ? AND hierarchical_Level[2] = ? AND hierarchical_Level[3] IS NULL", hierarchicalLevel_1, hierarchicalLevel_2).Find(&data).Error; err != nil {
	// 	return rModels.MrAsset{}, errors.New("no existe registro con ese id")
	// }
	m.db.Where("hierarchical_level[1] = ? AND hierarchical_level[2] = ? AND hierarchical_level[3] IS NULL", hierarchicalLevel_1, hierarchicalLevel_2).Find(&data)
	return data, nil
}
