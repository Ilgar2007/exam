CREATE TABLE branches (
  "id" uuid default uuid_generate_v4()  ,
  "name" varchar , 
  "phone" varchar , 
  "photo" varchar , 
  "work_start_hour" varchar ,
  "work_end_hour" varchar ,
  "address" varchar ,
  "delivery_price" int , 
  "active"  boolean ,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP
)