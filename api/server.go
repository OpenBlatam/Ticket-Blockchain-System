package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const baseUrl = "/ext"

type Server struct {
	log     logging.Logger
	router  *router
	portUrl string
}

// Initialize creates the API server at the provided port
func (s *Server) Initialize(log logging.Logger, factory logging.Factory, port uint16) {
	s.log = log
	s.factory = factory
	s.portURL = fmt.Sprintf(":%d", port)
	s.router = newRouter()
}

// Dispatch starts the API server
func (s *Server) Dispatch() error {
	handler := cors.Default().Handler(s.router)
	return http.ListenAndServe(s.portURL, handler)
}

// DispatchTLS starts the API server with the provided TLS certificate
func (s *Server) DispatchTLS(certFile, keyFile string) error {
	handler := cors.Default().Handler(s.router)
	return http.ListenAndServeTLS(s.portURL, certFile, keyFile, handler)
}

// RegisterChain registers the API endpoints associated with this chain That
// is, add <route, handler> pairs to server so that http calls can be made to
// the vm
func (s *Server) RegisterChain(ctx *snow.Context, vmIntf interface{}) {
	vm, ok := vmIntf.(common.VM)
	if !ok {
		return
	}

	// all subroutes to a chain begin with "bc/<the chain's ID>"
	defaultEndpoint := "bc/" + ctx.ChainID.String()
	httpLogger, err := s.factory.MakeChain(ctx.ChainID, "http")
	if err != nil {
		s.log.Error("Failed to create new http logger: %s", err)
		return
	}
	s.log.Verbo("About to add API endpoints for chain with ID %s", ctx.ChainID)

	// Register each endpoint
	for extension, service := range vm.CreateHandlers() {
		// Validate that the route being added is valid
		// e.g. "/foo" and "" are ok but "\n" is not
		_, err := url.ParseRequestURI(extension)
		if extension != "" && err != nil {
			s.log.Warn("could not add route to chain's API handler because route is malformed: %s", extension)
			continue
		}
		s.log.Verbo("adding API endpoint: %s", defaultEndpoint+extension)
		if err := s.AddRoute(service, &ctx.Lock, defaultEndpoint, extension, httpLogger); err != nil {
			s.log.Error("error adding route: %s", err)
		}
	}
}

// AddRoute registers the appropriate endpoint for the vm given an endpoint
func (s *Server) AddRoute(handler *common.HTTPHandler, lock *sync.RWMutex, base, endpoint string, log logging.Logger) error {
	url := fmt.Sprintf("%s/%s", baseURL, base)
	s.log.Info("adding route %s%s", url, endpoint)
	h := handlers.CombinedLoggingHandler(log, handler.Handler)
	switch handler.LockOptions {
	case common.WriteLock:
		return s.router.AddRouter(url, endpoint, middlewareHandler{
			before:  lock.Lock,
			after:   lock.Unlock,
			handler: h,
		})
	case common.ReadLock:
		return s.router.AddRouter(url, endpoint, middlewareHandler{
			before:  lock.RLock,
			after:   lock.RUnlock,
			handler: h,
		})
	case common.NoLock:
		return s.router.AddRouter(url, endpoint, h)
	default:
		return errUnknownLockOption
	}
}

// AddAliases registers aliases to the server
func (s *Server) AddAliases(endpoint string, aliases ...string) error {
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)
	endpoints := make([]string, len(aliases))
	for i, alias := range aliases {
		endpoints[i] = fmt.Sprintf("%s/%s", baseURL, alias)
	}
	return s.router.AddAlias(url, endpoints...)
}

// AddAliasesWithReadLock registers aliases to the server assuming the http read
// lock is currently held.
func (s *Server) AddAliasesWithReadLock(endpoint string, aliases ...string) error {
	// This is safe, as the read lock doesn't actually need to be held once the
	// http handler is called. However, it is unlocked later, so this function
	// must end with the lock held.
	s.router.lock.RUnlock()
	defer s.router.lock.RLock()

	return s.AddAliases(endpoint, aliases...)
}

// Call ...
func (s *Server) Call(
	writer http.ResponseWriter,
	method,
	base,
	endpoint string,
	body io.Reader,
	headers map[string]string,
) error {
	url := fmt.Sprintf("%s/vm/%s", baseURL, base)

	handler, err := s.router.GetHandler(url, endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "*", body)
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	handler.ServeHTTP(writer, req)

	return nil
}
