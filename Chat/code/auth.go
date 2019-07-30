package code

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
	"strings"
)

type authHandler struct{
	next http.Handler
}
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authentication started")
	//discard the cookie using _ character and only keep error
	_,err :=r.Cookie("auth")
	if err ==http.ErrNoCookie{
		//not authenticated
		w.Header().Set("Location","/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err !=nil{
		//some other error
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	//success -call the nexxt handler
	h.next.ServeHTTP(w,r)
}
func MustAuth(handler http.Handler) http.Handler{
	return &authHandler{next:handler}
}


// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}



func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler")
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		//get the provider object that matches the object specified in URL
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//get the location where we must send the users in order to start the author process
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location",loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",
				provider, err), http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s",
				provider, err), http.StatusInternalServerError)
			return
		}
		//Base64 to ensure there is no special or unpredicted character
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/"})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

