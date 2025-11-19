# Auriga_API

Auriga_API/
├── config/
│   └── config.go
├── databases/
│   ├── influxdbConnection.go
│   └── postgresqlConnection.go
├── internal/
│   ├── grafana/
│   │   └── client.go
│   ├── httpapi/
│   │   ├── handlers/
│   │   │   ├──  hAssets/
│   │   │   │   ├── asset.go
│   │   │   │   ├── assets.go
│   │   │   │   └── hDosingSystemByLine.go
│   │   │   ├──  hEvents/
│   │   │   │   ├── hEvents.go
│   │   │   │   ├── hEventsCategory.go
│   │   │   │   ├── hEventsCommit.go
│   │   │   │   ├── hEventsRaw.go
│   │   │   │   └── hEventsSap.go
│   │   │   ├──  hProducts/
│   │   │   │   ├── product.go
│   │   │   │   └── hEventsSap.go
│   │   │   ├──  hSap/
│   │   │   │   ├── hConsumptionByOrder.go
│   │   │   │   ├── hOrdersByLine.go
│   │   │   │   └── hEventsShRecipeByOrderap.go
│   │   ├── middlewares/
│   │   │   ├── middlewares.go
│   │   │   └── middlewaresMain.go
│   │   └── handlers.go
│   ├── repositories/
│   │   ├──  rAssets/
│   │   │   ├── rAssetHelper.go
│   │   │   ├── rAssetProduct.go
│   │   │   ├── rAssets.go
│   │   │   ├── rAssetsInit.go
│   │   │   └── rDosingSystemByLine.go
│   │   ├──  rEvents/
│   │   │   ├── rEvents.go
│   │   │   ├── rEventsCategory.go
│   │   │   ├── rEventsCommit.go
│   │   │   └── rEventsRaw.go
│   │   ├──  riInfluxdb/
│   │   │   └── riInfluxdb.go
│   │   ├──  rLineOrders/
│   │   │   ├── rConsumptionByOrder.go
│   │   │   ├── rLineOrders.go
│   │   │   ├── rOrdersByLine.go
│   │   │   └── rRecipeByOrder.go
│   │   ├──  rModels/
│   │   │   ├── rmAssets.go
│   │   │   ├── rmEvents.go
│   │   │   └── rmOrders.go
│   │   ├──  rOthers/
│   │   │   ├── InfluxGrafanaQuery.go
│   │   │   └── rOthers.go
│   │   ├──  rProducts/
│   │   │   ├── rProductDataTypeInit.go
│   │   │   ├── rProductHelper.go.go
│   │   │   ├── rProducts.go
│   │   │   ├── rProductsInitDrives.go
│   │   │   ├── rProductsInitExtrLinePartss.go
│   │   │   ├── rProductsInitExtrusionLines.go
│   │   │   ├── rProductsInitFactories.go
│   │   │   ├── rProductsInitMotors.go
│   │   │   ├── rProductsInitOtherProducts.go
│   │   │   ├── rProductTypes.go
│   │   │   └── rProdrProductTypesHelperucts.go
│   │   └──  rsSap/
│   │   │   ├── rsFactoryLineList.go
│   │   │   ├── rsLineOrderList.go
│   │   │   ├── rsLineRecipe.go
│   │   │   ├── rsLineStopEvent.go
│   │   │   └── rsSap.go
│   │   └── repositories.go
│   └── services/
│       ├──  sAssets/
│       │   ├── sAssetProduct.go
│       │   ├── sAssets.go
│       │   └── sDosingSystemByLine.go
│       ├──  sEvents/
│       │   ├── sEvents.go
│       │   ├── sEventsCategType.go
│       │   ├── sEventsCommit.go
│       │   ├── sEventsRaw.go
│       │   └── sEventsSap.go
│       ├──  sProducts/
│       │   └── sProducts.go
│       ├──  ssSap/
│       │   ├── sConsumptionByOrder.go
│       │   ├── sOrdersByLine.go
│       │   ├── sRecipeByOrder.go
│       │   └── sSap.go
│       └── services.go
├── .env
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
