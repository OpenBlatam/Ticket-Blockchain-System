package api

const baseUrl = "/ext"

type Server struct {
	log     logging.Logger
	router  *router
	portUrl string
}
