CREATE KEYSPACE platform_image_db with replication = {'class':'SimpleStrategy', 'replication_factor':1};
CREATE TABLE platform_image_db.api_metadata (tenant_id text, api_name text, parameters map<text,text>, PRIMARY KEY ((tenant_id), api_name) );
Insert into platform_image_db.api_metadata (tenant_id, api_name, parameters) values ('TEST01','google_image_search_api',{'key':'AIzaSyCkV8-7L8vJ7kUOuzkty9qBHYrFKUTvgJk','cx':'015644724633820582793%3Aynfqetiyhmi','url':'https://www.googleapis.com/customsearch/v1'})
CREATE TABLE platform_image_db.image_store(tenant_id text,term text,)

