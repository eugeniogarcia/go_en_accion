module gz.com/ch10

go 1.25.1

require (
	github.com/shopspring/decimal v1.4.0
	gz.com/ch10/paquetes/convert v0.0.0-00010101000000-000000000000
	gz.com/ch10/paquetes/person v0.0.0-00010101000000-000000000000
)

replace gz.com/ch10/paquetes/convert => ./paquetes/convert

replace gz.com/ch10/paquetes/person => ./paquetes/person
