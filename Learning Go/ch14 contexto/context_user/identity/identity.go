package identity

import (
	"context"
	"net/http"
)

// Declaramos un tipo con base int para usarlo como key en el contexto
type userKey int

// Definimos la lista de keys usando iota. En este caso solo tenemos una key. Con esto garantizamos la unicidad de la key. key no se exporta es solo para uso interno.
const (
	_ userKey = iota
	key
)

func ContextWithUser(ctx context.Context, user string) context.Context {
	// Usamos la helper function context.WithValue para crear un nuevo contexto que wrappea ctx e incluye una nueva key/valor
	return context.WithValue(ctx, key, user)
}

func UserFromContext(ctx context.Context) (string, bool) {
	// Recuperamos un valor. Como es el valor es un interface - any - podemos hacer type assertion. En este caso esperamos un string. ok sera true si la aserción ha ido bien. En user tendremos el valor asociado a la key
	user, ok := ctx.Value(key).(string)
	return user, ok
}

// a real implementation would be signed to make sure
// the identity didn't spoof their identity
func extractUser(req *http.Request) (string, error) {
	userCookie, err := req.Cookie("identity")
	if err != nil {
		return "", err
	}
	return userCookie.Value, nil
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		user, err := extractUser(req)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("unauthorized"))
			return
		}
		ctx := req.Context()
		ctx = ContextWithUser(ctx, user)
		req = req.WithContext(ctx)
		h.ServeHTTP(rw, req)
	})
}

func SetUser(user string, rw http.ResponseWriter) {
	http.SetCookie(rw, &http.Cookie{
		Name:  "identity",
		Value: user,
	})
}

func DeleteUser(rw http.ResponseWriter) {
	http.SetCookie(rw, &http.Cookie{
		Name:   "identity",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
