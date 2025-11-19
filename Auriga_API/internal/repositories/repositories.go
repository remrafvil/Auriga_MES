package repositories

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rAssets"
	"github.com/remrafvil/Auriga_API/internal/repositories/rAuth"
	"github.com/remrafvil/Auriga_API/internal/repositories/rDocuments"
	"github.com/remrafvil/Auriga_API/internal/repositories/rEvents"
	"github.com/remrafvil/Auriga_API/internal/repositories/rInfluxQuery"
	"github.com/remrafvil/Auriga_API/internal/repositories/rLabor"
	"github.com/remrafvil/Auriga_API/internal/repositories/rLabor_KKK"
	"github.com/remrafvil/Auriga_API/internal/repositories/rLineOrders"
	"github.com/remrafvil/Auriga_API/internal/repositories/rProducts"
	"github.com/remrafvil/Auriga_API/internal/repositories/rUsers"
	"github.com/remrafvil/Auriga_API/internal/repositories/riInfluxdb"
	"github.com/remrafvil/Auriga_API/internal/repositories/rsSap"
	"github.com/remrafvil/Auriga_API/internal/repositories/rwWorkera"
	"go.uber.org/fx"
)

var Module = fx.Module("repositories", fx.Provide(
	rsSap.New,
	riInfluxdb.New,
	rUsers.New,
	rLineOrders.New,
	rEvents.New,
	rProducts.New,
	rAssets.New,
	rAuth.New,
	rDocuments.New,
	rInfluxQuery.New,
	rLabor.New,
	rLabor_KKK.New,
	rwWorkera.New,
))
