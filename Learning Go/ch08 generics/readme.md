- `comparable_stack`. Implementa usando generics un tipo `Stack`. Al declarar el parametro especificamos un interface. Si elegimos `any` - esto es `Stack[T any]` - `T` seria cualquier tipo, y por lo tanto podría ser algo no comparable. Si queremos asegurar que el tipo que usemos como parametro tenemos que especificar un interface más restrictivo. En este ejemplo indicamos `comparable` para restringir `T` a un tipo que se pueda comparar.

- Con `map_filter_reduce` mostramos que también se pueden usar generics con funciones (las funciones son un tipo y generics se puede usar con todos los tipos). Implementamos un Map, Reduce & un Filter.

- En `impossible` mostramos con restringir los tipos que se pueden usar en un generic. Definimos un interface en el que además de indicar los métodos que hay que implementar, se indica una relación de tipos "admisibles". Este interface no se puede usar para definir una variable, solo sirve para definir los parametros de un generic
