CREATE KEYSPACE platform_image_db  with replication = {'class':'SimpleStrategy'};
use platform_image_db;

CREATE TYPE platform_image_db.image_store_data (
  url_link text,
  store_url text,
);

CREATE TABLE platform_image_db.image_metadata
(tenant_id text,
 search_term text,
 search_term_alias text,
 store_type text,
 image_store_by_title map<text,text>,
 image_count int,
 PRIMARY KEY (tenant_id, search_term, search_term_alias));

CREATE TABLE platform_image_db.api_metadata
(tenant_id text,
 api_name text,
 parameters map<text,text>,
 PRIMARY KEY (tenant_id));

