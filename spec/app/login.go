package app

// Authentication is a marker interface to declare which kind of authentication is used.
// If any endpoint returns http.StatusUnauthenticated (401) the frontend will present the according
// login possibility. However, http.StatusForbidden (403) will only display a notification.
type Authentication interface {
	IsAuthentication() bool
}

// None means that the application does not have any authentication mechanism and that there is no
// login page nor any redirect.
type None struct {
}

func (None) IsAuthentication() bool {
	return true
}

// Bearer is more like a debug or developer mechanism and not suited for normal end users.
// This means, that the application presents a single input field to enter the bearer token.
// The user must get his token over a secure and independent channel.
type Bearer struct {
}

func (Bearer) IsAuthentication() bool {
	return true
}
