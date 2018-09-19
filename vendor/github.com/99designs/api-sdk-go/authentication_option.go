package smartling

// AuthenticationOption specifies should request to API use authentication or
// not. See Post and Get methods.
type AuthenticationOption bool

const (
	// WithAuthentication equal to use of authentication in request.
	WithAuthentication = AuthenticationOption(true)

	// WithoutAuthentication equal to not use of authentication in request.
	WithoutAuthentication = AuthenticationOption(false)
)
