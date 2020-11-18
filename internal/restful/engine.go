package restful

import (
	"context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/zliang90/kingRest/internal/app/conf"
	"github.com/zliang90/kingRest/internal/restful/api"
	"github.com/zliang90/kingRest/internal/restful/errors"
	"github.com/zliang90/kingRest/internal/restful/router"
	"github.com/zliang90/kingRest/internal/restful/validator"
	"github.com/zliang90/kingRest/pkg/log"
	"github.com/zliang90/kingRest/pkg/util/table"
	"net/http"
	"time"
)

type Server struct {
	// env
	env string

	// web server config
	cf conf.WebServer

	// serve handler
	r http.Handler

	// http server
	srv *http.Server
}

func New(cf *conf.Config) *Server {
	return &Server{
		env: cf.Env,
		cf:  cf.WebServer,
	}
}

func (s *Server) Run(ctx context.Context) {
	defer func() {
		if e := recover(); e != nil {
			log.Errorf("%v", e)
		}
	}()
	// init engine, and run it
	log.Info("start restful api")
	go s.serve()

	select {
	case <-ctx.Done():
		log.Info("stopping restful api")
		if err := s.srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (s *Server) serve() {
	log.Infof("register routers")
	s.r = router.InitRouter(s.env)

	// upgrade to v10
	log.Infof("upgrade gin validator to v10")
	validator.ToV10()

	// no route
	s.r.(*gin.Engine).NoRoute(func(c *gin.Context) {
		api.Failure(c, errors.NotFound(
			fmt.Sprintf("%s '%s'", c.Request.Method, c.Request.URL)))
	})
	// middleware
	middleware := []gin.HandlerFunc{
		handlerRequestId(),
		handlerRecovery(),
	}
	if s.env != "prod" {
		middleware = append(middleware, handlerLogger())
		log.Debugf("register pprof routers")
		pprof.Register(s.r.(*gin.Engine))

		s.printRouters()
	}
	s.r.(*gin.Engine).Use(middleware...)

	log.Infof("serving restful api on: %s", s.cf.Addr)

	// web http server
	srv := &http.Server{
		Addr:              s.cf.Addr,
		Handler:           s.r,
		MaxHeaderBytes:    s.cf.MaxHeaderBytes,
		ReadTimeout:       time.Duration(s.cf.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.cf.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.cf.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.cf.IdleTimeout) * time.Second,
	}
	s.srv = srv
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
		log.Info(err)
	}
	//log.Fatal(s.r.(*gin.Engine).Run(s.l))
}

func (s *Server) printRouters() {
	t := &table.Table{}

	tHeader := []string{"METHOD", "PATH"}
	tData := make([][]string, 0)

	routers := s.r.(*gin.Engine).Routes()
	for _, r := range routers {
		d := []string{r.Method, r.Path}
		tData = append(tData, d)
	}
	t.Header = tHeader
	t.Data = tData

	t.PrintTable(true)
}
