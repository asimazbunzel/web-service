
package web

import (
   "net/http"
   "os"
   "sync"
   "time"

   "web-service/pkg/io"
   
   "github.com/julienschmidt/httprouter"
   "github.com/kardianos/service"
)


// Struct with information on the service 
type ServiceInfo struct {
   ServiceName string
   ServiceDescription string
   Version string
}


// some useful variables for the server
var (
   serviceIsRunning bool
   programIsRunning bool
   writingSync    sync.Mutex
)


// wrapper structure for start & stop service
type Program struct{}

// Program method to start service
func (p Program) Start(s service.Service) error {
   io.LogDebug("WEB - server.go - Start", s.String() + " started")
   writingSync.Lock()
   serviceIsRunning = true
   writingSync.Unlock()
   go p.run()
   return nil
}

// Program method to stop service
func (p Program) Stop(s service.Service) error {
   writingSync.Lock()
   serviceIsRunning = false
   writingSync.Unlock()
   for programIsRunning {
      io.LogDebug("WEB - server.go - Stop", s.String() + " stopping...")
      time.Sleep(5 * time.Second)
   }
   io.LogDebug("WEB - server.go - Stop", s.String() + " stopped")
   return nil
}

// Program method that runs service
func (p Program) run() {

   io.LogDebug("WEB - server.go - run", "serving web files")

   router := httprouter.New()
   router.ServeFiles("/html/*filepath", http.Dir("web/html"))
   router.ServeFiles("/css/*filepath", http.Dir("web/css"))
   router.ServeFiles("/js/*filepath", http.Dir("web/js"))
   router.ServeFiles("/vendors/*filepath", http.Dir("web/vendors"))
 
   router.GET("/", BasicAuth(Index))
   router.GET("/index", BasicAuth(Index))
   router.GET("/dashboard", BasicAuth(Dashboard))
   router.GET("/mesa", BasicAuth(MESAhtml))

   // get port number from env variables
   port := os.Getenv("PORT")
   if port == "" {
      port = "8080"
   }
   io.LogInfo("WEB - server.go - run", "setting PORT to: " + port)

   err := http.ListenAndServe(":"+port, router)
   if err != nil {
      io.LogError("WEB - server.go - run", "problem starting web server: " + err.Error())
      os.Exit(-1)
   }

   io.LogInfo("WEB - server.go - run", "service running")

}


// Initialization of the server
func InitServer(s *ServiceInfo) {

   // Start service Config
   serviceConfig := &service.Config{
      Name:        s.ServiceName,
      DisplayName: s.ServiceName,
      Description: s.ServiceDescription,
   }

   prg := &Program{}

   serv, err := service.New(prg, serviceConfig)
   if err != nil {
      io.LogDebug("WEB - server.go - InitServer", "cannot create the service: " + err.Error())
   }


   err = serv.Run()
   if err != nil {
      io.LogDebug("WEB - server.go - InitServer", "cannot start the service: " + err.Error())
   }
}
