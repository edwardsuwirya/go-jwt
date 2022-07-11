package request

type AuthHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}
