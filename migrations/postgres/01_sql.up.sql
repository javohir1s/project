CREATE TABLE "branch" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR(64), 
    "address" VARCHAR(64),
    "phone" VARCHAR(20),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "client" (
    "id" UUID NOT NULL PRIMARY KEY,
    "first_name" VARCHAR(48), 
    "last_name" VARCHAR(48),
    "father_name" VARCHAR(48),
    "phone" VARCHAR(24),
    "birthday" DATE,
    "gender" VARCHAR(10) CHECK ("gender" IN ('male', 'female')),
    "branch_id" UUID REFERENCES "branch"("id"), 
    "active" BOOLEAN DEFAULT true,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "product" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR(48),
    "price" NUMERIC,
    "branch_id" UUID REFERENCES "branch"("id"), 
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "coming" (
    "id" UUID NOT NULL PRIMARY KEY,
    "increment_id" VARCHAR(12),
    "branch_id" UUID REFERENCES "branch"("id"), 
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "picking_list" (
    "id" UUID NOT NULL PRIMARY KEY,
    "product_id" UUID REFERENCES "product"("id"),
    "price" NUMERIC,
    "quantity" INT,
    "total_price" NUMERIC,
    "coming_id" UUID REFERENCES "coming"("id"),
    "coming_increment_id" VARCHAR(12),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "remainder" (
    "id" UUID NOT NULL PRIMARY KEY,
    "product_id" UUID REFERENCES "product"("id"),
    "name" VARCHAR(24),
    "quantity" INT,
    "coming_price" NUMERIC,
    "sale_price" NUMERIC,
    "branch_id" UUID REFERENCES "branch"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "sale" (
    "id" UUID NOT NULL PRIMARY KEY,
    "client_id" UUID REFERENCES "client"("id"),
    "branch_id" UUID REFERENCES "branch"("id"),
    "increment_id" VARCHAR(12),
    "total_price" NUMERIC,
    "paid" NUMERIC,
    "debt" NUMERIC, 
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "sale_product" (
    "id" UUID NOT NULL PRIMARY KEY,
    "product_id" UUID REFERENCES "product"("id"),
    "sale_id" UUID REFERENCES "sale"("id"),
    "sale_increment_id" VARCHAR(12),
    "quantity" INT,
    "price" NUMERIC,
    "total_price" NUMERIC,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);