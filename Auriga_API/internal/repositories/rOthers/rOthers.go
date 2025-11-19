package rOthers

import (
	"github.com/lib/pq"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

type MrDocumentsBase64 struct {
	gorm.Model
	DocName  string              `gorm:"type:text;not null"`
	DocFile  string              `gorm:"type:text;not null"`
	DocType  string              `gorm:"type:text;not null"`
	Products []rModels.MrProduct `gorm:"many2many:mr_produc_Documentsbase64;"`
}

type MrGrafanaDashboards struct {
	gorm.Model
	AssetType     string         `gorm:"type:text;not null"`
	DashboardID   pq.StringArray `gorm:"type:text[];not null"`
	DashboardName string         `gorm:"type:text;not null"`
}

func DB_InitDocumentsBase64(c *gorm.DB) { //
	c.Create(&MrDocumentsBase64{DocName: "", DocFile: "", DocType: ""})
}

func DB_InitGrafanaDashboards(c *gorm.DB) { //										General			Process	        Compsumptions	   	Events			Maintenance
	c.Create(&MrGrafanaDashboards{AssetType: "General", DashboardID: /*       */ []string{"ae9r8d95ad9mof", "fe9r8erhrc7i8d", "de9r8hhjnumf4c", "be9r8ib5tifi8f", "fe9r8j90fz4e8b"}, DashboardName: "0100company"})
	c.Create(&MrGrafanaDashboards{AssetType: "Factory", DashboardID: /*       */ []string{"fe9r8psoelvcwf", "fe9r8t3f3cdtse", "ae9r8u8xdlog0f", "ee9r8z78g658ga", "ce9r936udil1cf"}, DashboardName: "0200factory"})
	c.Create(&MrGrafanaDashboards{AssetType: "ExtrusionLine", DashboardID: /* */ []string{"ee9r9qc4slj40a", "de9r9tyyon6dcf", "be9r9v2j6ql8gd", "ce9r9zidntm2of", "de9r9yp35or28c"}, DashboardName: "0310lineextr"})
	c.Create(&MrGrafanaDashboards{AssetType: "LineCommon", DashboardID: /*    */ []string{"              ", "eeeixyl986s5ce", "              ", "              ", "              "}, DashboardName: "0310lineextrcomm"})
	c.Create(&MrGrafanaDashboards{AssetType: "LineTransport", DashboardID: /* */ []string{"              ", "deeixb4uady4gf", "              ", "              ", "              "}, DashboardName: "0311lineextrtrans"})
	c.Create(&MrGrafanaDashboards{AssetType: "LineDosing", DashboardID: /*    */ []string{"              ", "aeeixmspfyqyof", "              ", "              ", "              "}, DashboardName: "0312lineextrtdosing"})
	c.Create(&MrGrafanaDashboards{AssetType: "LineExtrusion", DashboardID: /* */ []string{"              ", "aeeixocelyfi8e", "              ", "              ", "              "}, DashboardName: "0313lineextrtextr"})
	c.Create(&MrGrafanaDashboards{AssetType: "LineStack", DashboardID: /*     */ []string{"              ", "deeixpk035ou8d", "              ", "              ", "              "}, DashboardName: "0315lineextrtstack"})
	c.Create(&MrGrafanaDashboards{AssetType: "LineWinders", DashboardID: /*   */ []string{"              ", "eeeiyjsryv9j4c", "              ", "              ", "              "}, DashboardName: "0318lineextrwinders"})
	// 	c.Create(&MrGrafanaDashboards{AssetType: "LineCommon", DashboardID: "de9o4edxu8buoe", DashboardName: "0000-base1"})
	// 	c.Create(&MrGrafanaDashboards{AssetType: "LineTransport", DashboardID: "fe9o5j4fhp7nkc", DashboardName: "0000base"})
	// 	c.Create(&MrGrafanaDashboards{AssetType: "LineDosing", DashboardID: "de9o4edxu8buoe", DashboardName: "0000-base1"})
	// 	c.Create(&MrGrafanaDashboards{AssetType: "LineExtrusion", DashboardID: "fe9o5j4fhp7nkc", DashboardName: "0000base"})
	// 	c.Create(&MrGrafanaDashboards{AssetType: "LineStack", DashboardID: "de9o4edxu8buoe", DashboardName: "0000-base1"})
	// 	c.Create(&MrGrafanaDashboards{AssetType: "LineWinders", DashboardID: "fe9o5j4fhp7nkc", DashboardName: "0000base"})
}
