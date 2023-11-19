CREATE TABLE clients (
  "id" uuid primary key default uuid_generate_v4(),
  "first_name" varchar ,
  "last_name" varchar , 
  "phone" varchar ,
  "photo" varchar , 
  "date_of_birth" varchar ,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP
)
