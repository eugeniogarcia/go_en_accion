package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {
	// Métodos simples para crear un log
	slog.Debug("debug log message")
	slog.Info("info log message")
	slog.Warn("warning log message")
	slog.Error("error log message")

	// Podemos también enviar pares key/value
	userID := "fred"
	loginCount := 20
	slog.Info("user login",
		"id", userID, //primer par key/value
		"login_count", loginCount) //segundo key/value

	// Si necesitamos enviar la información estructurada, por ejemplo con un json, creamos un handler
	//1. definimos las opciones del handler
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	//2. creamos el handler, en este caso un JSONHandler. Para crearlos indicamos un io.Writer, y las opciones
	handler := slog.NewJSONHandler(os.Stderr, options)
	//3. creamos el logger a partir del handler
	mySlog := slog.New(handler) //a partir del handler creamos un logger, que pasamos a utilizar con los métodos estadard

	lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)
	mySlog.Debug("debug message", "id", userID, "last_login", lastLogin)

	// Para optimizar el rendimiento del logger podemos usar LogAttrs en lugar de Info, Debug, etc
	// Toma un contexto, y la información a registrar. Los pares key-value se crear con helpers. slog.Any nos sirve para cualquier tipo
	ctx := context.Background()
	mySlog.LogAttrs(ctx, slog.LevelInfo, "faster logging", slog.String("id", userID), slog.Time("last_login", lastLogin))

	myLog := slog.NewLogLogger(mySlog.Handler(), slog.LevelDebug)
	myLog.Println("using the mySlog Handler")
}
