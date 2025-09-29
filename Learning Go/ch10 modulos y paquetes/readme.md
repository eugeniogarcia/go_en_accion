## setup

En el directorio raiz inicializamos un modulo

```ps
go mod init gz.com/ch10
```

una vez iniciado el modulo, en el directorio donde tenemos el paquete, hacemos lo siguiente para inicializar el paquete con nombre `gz.com/ch10/paquetes/person`. Se creará el archivo `go.mod` en el directorio.

```ps
go mod init gz.com/ch10/paquetes/person

go mod tidy
```

hacemos también `go mod tidy` en el directorio raíz

con esto ya podemos importar el paquete:

```go
import (
	"fmt"

	"gz.com/ch10/paquetes/person"
)
```

Se usa la directiva replace para indicar que el modulo no esta disponible en internet sino en local:

```go
replace gz.com/ch10/paquetes/convert => ./paquetes/convert

replace gz.com/ch10/paquetes/person => ./paquetes/person
```

## exporting

Todas las constantes, variables o funciones del paquete que queramos exportar tendremos que capitalizarlas. Los items que no esten capitalizados no estarán visibles fuera del paquete.

## alias

Podemos poner un alias al importar un paquete. Esto nos vendrá muy bien cuando haya dos paquetes diferentes con el mismo nombre.

```go
import (
	"fmt"

	persona "gz.com/ch10/paquetes/person"
)
```

## documentacion

Si documentamos el código podemos usar `go doc [paquete]` o `go doc [paquete.elemento]` para acceder a la documentación.

En `convert` tenemos un ejemplo de como documentar el código:

- Place the comment directly before the item being documented, with no blank lines between the comment and the declaration of the item.

- Start each line of the comment with double slashes (//), followed by a space.  While it's legal to use /* and */ to mark your comment block, it is idiomatic to use double slashes.

- The first word in the comment for a symbol (a function, type, constant, variable, or method) should be the name of the symbol. You can also use "A" or "An" before the symbol name to help make the comment text grammatically correct.

- Use a blank comment line (double slashes and a newline) to break your comment into multiple paragraphs.

- If you want your comment to contain some preformatted content (such as a table or source code), put an additional space after the double slashes to indent the lines with the content.

- If you want a header in your comment, put a # and a space after the double slashes. Unlike with Markdown, you cannot use multiple # characters to make different levels of headers.

- To make a link to another package (whether or not it is in the current module), put the package path within brackets ([ and ]).

- To link to an exported symbol, place its name in brackets. If the symbol is in another package, use [pkgName.SymbolName].

- If you include a raw URL in your comment, it will be converted into a link.

- If you want to include text that links to a web page, put the text within brackets ([ and ]). At the end of the comment block, declare the mappings between your text and their URLs with the format // [TEXT]: URL. You'll see a sample of this in a moment.

### herramienta

Podemos usar una herramienta `pkgsite` para acceder a la documentación. La herramienta se instala así:
```ps
go install golang.org/x/pkgsite/cmd/pkgsite@latest
```

para arrancar la herramienta:

```ps
pkgsite
```

podemos acceder al site en `http://localhost:8080/`

## internal

Hay una convención de go que consiste en considerar como interno un modulo que este definido en un subdirectorio llamado `internal`. Los elementos exportados por un paquete `internal` solo estarán accesibles desde su paquete padre (no el abuelo, solo desde su padre), y desde sus paquetes hermanos.

## init

Cuando se importa por primera vez un paquete se ejecutaran todas las funciones `func init() {}` que esten definidas en él. He creado un par de funciones en el paquete `gz.com/ch10/paquetes/person`.

## get modules

Podemos hacer lo siguiente para descargar todos los modulos _requeridos_ en todos los paquetes de la estructura de directorios. El comando get recupera los datos y crea o actualizar el `go.mod`:

```ps
go get ./...
```

la otra opción es descargar un módulo concreto:

```ps
go get [modulo]
```

Los modulos que se descargan se guardan en la cache. Este comando recupera todas las variables de entorno de GO:

```ps
go env
```

con `GOMODCACHE` podemos ver la ruta en la que se están cacheando los módulos:

```ps
go env GOMODCACHE
```

podemos limpiar la cache con:

```ps
go clean -modcache
```

se marcan como `indirect` en `go.mod` aquellas dependencias que son dependencias de los módulos que importamos.

Podemos ver las versiones instaladas con:

```ps
go list
```