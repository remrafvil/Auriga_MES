package services

import (
	"github.com/remrafvil/Auriga_API/internal/services/sAssets"
	"github.com/remrafvil/Auriga_API/internal/services/sAuth"
	"github.com/remrafvil/Auriga_API/internal/services/sEvents"
	"github.com/remrafvil/Auriga_API/internal/services/sInfluxQuery"
	"github.com/remrafvil/Auriga_API/internal/services/sLabor"
	"github.com/remrafvil/Auriga_API/internal/services/sLabor1"
	"github.com/remrafvil/Auriga_API/internal/services/sProducts"
	"github.com/remrafvil/Auriga_API/internal/services/sSap"
	"github.com/remrafvil/Auriga_API/internal/services/sUsers"
	"go.uber.org/fx"
)

var Module = fx.Module("services", fx.Provide(
	sAssets.New,
	sAuth.NewJWKSValidator,
	sAuth.New,
	sProducts.New,
	sSap.New,
	sEvents.New,
	//sLandProperties.New,
	sUsers.New,
	sInfluxQuery.New,
	sLabor1.New,
	sLabor.New,
))
