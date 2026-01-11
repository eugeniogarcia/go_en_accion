En esta carpeta hay tres variantes en las que en lugar de usar Postgres usamos MySql, MongoDB o DynamoDB. He dejado únicamente los archivos en los que hay alguna diferencia.

- Archivos de configuración `toml`, porque los atributos para conectarnos con la base de datos difieren
- `Main.go`, porque se importa el driver que usamos, y por lo tanto el paquete que importamos no es el mismo, aunque el resto del código si lo es
- `dbscripts`. Los scripts para crear los artefactos en la base de datos (no es código _go_ de la aplicación propiamente dicho) es diferente, incluso en el caso de Postgres y MySql aunque ambas usan SQL, hay algunas diferencias en los dialectos
- `repositories`. Aqui era de esperar que hubiera diferencias porque es la capa de acceso a datos. En las capas "superiores" - controller y servicios -, no hay diferencias. Comentar los siguientes cambios:
    - En MySql con respecto a Postgres
        - Se usa diferente placeholder ($1, $2, frente a ?, ?,...)
        - El campo ID se formatea de forma diferente
        - Hay diferencias pequeñas en las sentencias SQL (por ejemplo, no se admite la sentencia `RETURNING`, o los booleanos pasan de ser `false` a ser `FALSE`)
    - En MongoDB manejamos colecciones en lugar de tablas, y desnormalizamos el modelo de modo que results y runners se fusionan en la misma colección. Las primitivas de acceso a los datos ya no son SQL, y lo que manejamos son Binary Json (`BSON`)
    - En DynamoDB las tablas son key/value stores, las primitivas de acceso no son tampoco SQL, y solo tenemos tipados los atributos que usaremos como primary key o indices secundarios