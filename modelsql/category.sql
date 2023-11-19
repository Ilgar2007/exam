CREATE TABLE category (
  "id" uuid primary key default  uuid_generate_v4() , 
  "title" varchar ,
  "image" varchar , 
  "parent_id" uuid references "category"("id"),
  "created_at" timestamp default current_timestamp,   
  "updated_at" timestamp
);