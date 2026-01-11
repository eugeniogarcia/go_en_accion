-- incluimos la extesión pgcrypto para hashear contraseñas
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

CREATE TABLE users (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    username text NOT NULL UNIQUE,
    user_password text NOT NULL,
    user_role text NOT NULL,
    access_token text,
    CONSTRAINT users_pk PRIMARY KEY (id) -- definimos la clave principal de la tabla
);

CREATE INDEX user_access_token
ON users (access_token); -- creamos un índice en la columna access_token para optimizar las consultas que filtren por token de acceso

-- La función gen_salt('bf') genera una sal aleatoria para el algoritmo Blowfish. Produce una salida diferente cada vez que se llama. La salida tiene la siguiente forma: $2a$<cost>$<22 character salt>, esto es, un identificador del alfgortimo que se ha usado ($2a$ corresponde al algoritmo blowfish), el coste de computación (cost) y el salt propiamente dicho (que tendra 22 caracteres de largo).
-- La función crypt() toma dos argumentos: el valor a hashear  y la sal (salt). Del salt toma el algoritmo y la salt propiamente dicha para hashear el valor. El valor y la salt se combinan (es más complejo que una concatenación de ambos) y se aplica el algoritmo de hash especificado en la salt para producir el hash resultante. El resultado es una cadena que incluye el identificador del algoritmo, el coste, la salt y el hash resultante
-- guardamos la contraseñas de los dos usuarios que hemos creado hasheadas con un salt y utilizando el algoritmo Blowfish
INSERT INTO users(username, user_password, user_role)
VALUES
    ('admin', crypt('admin', gen_salt('bf')), 'admin'),
    ('runner', crypt('runner', gen_salt('bf')), 'runner');
