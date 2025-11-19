-- SERVER SITE 1 (CXC)
-- Crear usuario para replicación
CREATE USER replicator_cxc WITH REPLICATION PASSWORD 'TYzPrJX5qlJ2o47y6ZqT25';
GRANT CONNECT ON DATABASE "CX_OT_CXC" TO replicator_cxc;

-- Crear publicación para la tabla mr_raw_events
CREATE PUBLICATION pub_cxc_mr_raw_events FOR TABLE mr_raw_events;

-- Dar permisos al usuario replicator
GRANT USAGE ON SCHEMA public TO replicator_cxc;
GRANT SELECT ON mr_raw_events TO replicator_cxc;





-- SERVER SITE 2 (FSP)
-- Crear usuario para replicación
CREATE USER replicator_fsp WITH REPLICATION PASSWORD 'TYzPrJX5qlJ2o47y6ZqT21';
GRANT CONNECT ON DATABASE "CX_OT_FSP" TO replicator_fsp;

-- Crear publicación para la tabla mr_raw_events
CREATE PUBLICATION pub_fsp_mr_raw_events FOR TABLE mr_raw_events;

-- Dar permisos al usuario replicator
GRANT USAGE ON SCHEMA public TO replicator_fsp;
GRANT SELECT ON mr_raw_events TO replicator_fsp;






-- SERVER CLOUD 
-- Crear suscripción para Site 1 (CXC)
CREATE SUBSCRIPTION sub_cxc_mr_raw_events
CONNECTION 'host=10.80.50.46 port=33307 dbname=CX_OT_CXC user=replicator_cxc password=TYzPrJX5qlJ2o47y6ZqT25'
PUBLICATION pub_cxc_mr_raw_events
WITH (
    copy_data = true,
    create_slot = true,
    enabled = true
);

-- Crear suscripción para Site 2 (FSP)
CREATE SUBSCRIPTION sub_fsp_mr_raw_events
CONNECTION 'host=172.21.99.1 port=33307 dbname=CX_OT_FSP user=replicator_fsp password=TYzPrJX5qlJ2o47y6ZqT21'
PUBLICATION pub_fsp_mr_raw_events
WITH (
    copy_data = true,
    create_slot = true,
    enabled = true
);