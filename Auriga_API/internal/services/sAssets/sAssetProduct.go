package sAssets

import (
	"log"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

type msValueType string

const (
	ValueTypeString msValueType = "string"
	ValueTypeNumber msValueType = "number"
)

type msProductFeature struct {
	Name        string
	Symbol      string
	Unit        string
	ValueType   msValueType
	VisuOrder   uint
	StringValue string
	NumberValue float64
}

type msProduct struct {
	ID              uint
	Name            string
	ProductTypeName string
	Manufacturer    string
	Family          string
	Description     string
	FeatureValues   []msProductFeature
	Documents       []sDocuments
}

type msAssetLong struct {
	ID        uint
	Product   msProduct
	Location  string
	TechCode  uint
	Code      string
	Sn        string
	SapCode   string
	Documents []sDocuments
}

type sDocuments struct {
	ID          uint
	Nombre      string
	Descripcion string
	URL         string
}

type msAssetHierarchy struct {
	ID                uint
	ProductID         uint
	ParentID          *uint
	Location          string
	TechCode          uint
	Code              string
	Sn                string
	SapCode           string
	HierarchicalLevel []string
}

func (s *service) AssetInfoDetail(id uint) (*msAssetLong, error) {
	rData, err := s.repository.AssetInfo(id)
	if err != nil {
		log.Println("Error service aqui:", err)
		return nil, err
	}
	data := &msAssetLong{
		ID:   rData.ID,
		Code: rData.Code,
	}
	return data, nil
}

// GetAssetWithSimplifiedProduct obtiene un activo con los detalles del producto en formato simplificado
func (s *service) GetAssetWithSimplifiedProduct(assetID uint) (msAssetLong, error) {
	// Obtener los datos completos del repositorio
	log.Println("*****************      Llego por aquí SERVICE *****************")
	dbAsset, err := s.repository.GetAssetWithProductDetails(assetID)
	if err != nil {
		return msAssetLong{}, err
	}

	// Transformar los datos al formato simplificado
	return transformToSimplifiedAsset(dbAsset), nil
}

// transformToSimplifiedAsset convierte los modelos de DB a los modelos simplificados
func transformToSimplifiedAsset(dbAsset rModels.MrAsset) msAssetLong {
	// Transformar el producto
	simplifiedProduct := msProduct{
		ID:              dbAsset.Product.ID,
		Name:            dbAsset.Product.Name,
		ProductTypeName: dbAsset.Product.ProductType.Name,
		Manufacturer:    dbAsset.Product.Manufacturer,
		Family:          dbAsset.Product.Family,
		Description:     dbAsset.Product.Description,
		FeatureValues:   make([]msProductFeature, 0, len(dbAsset.Product.FeatureValues)),
		Documents:       transformToSimplifiedDocuments(dbAsset.Product.Documents),
	}

	// Transformar los valores de las características
	for _, fv := range dbAsset.Product.FeatureValues {
		// Buscar el orden de visualización (VisuOrder)
		visuOrder := getVisuOrder(dbAsset.Product.ProductType.Relations, fv.MrProductFeatureTypeID)

		simplifiedProduct.FeatureValues = append(simplifiedProduct.FeatureValues, msProductFeature{
			Name:        fv.ProductFeature.Name,
			Symbol:      fv.ProductFeature.Symbol,
			Unit:        fv.ProductFeature.Unit,
			ValueType:   msValueType(fv.ProductFeature.ValueType),
			VisuOrder:   visuOrder,
			StringValue: fv.StringValue,
			NumberValue: fv.NumberValue,
		})
	}

	// Transformar el activo
	return msAssetLong{
		ID:        dbAsset.ID,
		Product:   simplifiedProduct,
		Location:  dbAsset.Location,
		TechCode:  dbAsset.TechCode,
		Code:      dbAsset.Code,
		Sn:        dbAsset.Sn,
		SapCode:   dbAsset.SapCode,
		Documents: transformToSimplifiedDocuments(dbAsset.Documents),
	}
}

// getVisuOrder busca el orden de visualización para una característica
func getVisuOrder(relations []rModels.MrProductFeatureTypeRelation, featureTypeID uint) uint {
	for _, rel := range relations {
		if rel.MrProductFeatureTypeID == featureTypeID {
			return rel.VisuOrder
		}
	}
	return 0 // Valor por defecto si no se encuentra
}

// transformToSimplifiedDocuments convierte los documentos a formato simplificado
func transformToSimplifiedDocuments(dbDocs []rModels.MrDocuments) []sDocuments {
	docs := make([]sDocuments, 0, len(dbDocs))
	for _, d := range dbDocs {
		docs = append(docs, sDocuments{
			ID:          d.ID,
			Nombre:      d.Nombre,
			Descripcion: d.Descripcion,
			URL:         d.URL,
		})
	}
	return docs
}
