CREATE TABLE allergens
(
    id character(36) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    name varchar(40) NOT NULL,
    description TEXT,
	CONSTRAINT allergen_id PRIMARY KEY (id)
);

CREATE TABLE dishes
(
    id character(36) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    name varchargit (40) NOT NULL,
	description TEXT,
    CONSTRAINT dish_id PRIMARY KEY (id)
);

CREATE TABLE dish_allergens
(
    id character(36) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    dish_id character(36) COLLATE pg_catalog."default" NOT NULL,
    allergen_id character(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT dish_allergen_id PRIMARY KEY (id),
    CONSTRAINT dish_id_fkey FOREIGN KEY (dish_id)
        REFERENCES dishes (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT allergen_id_fkey FOREIGN KEY (allergen_id)
        REFERENCES allergens (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE menu_entries
(
    id character(36) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    day DATE NOT NULL,
    dish_id character(36) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT dish_id_fkey FOREIGN KEY (dish_id)
        REFERENCES dishes (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

