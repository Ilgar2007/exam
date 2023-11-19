CREATE TABLE products (
  "id" varchar primary key ,
  "category_id" uuid references "category"("id") , 
  "title"  varchar ,
  "description" varchar,
  "photo" varchar ,
  "price" int ,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP
)

