package web

import (
   "crypto/sha256"
   "crypto/subtle"
   "html/template"
	"net/http"
   "os"
   "time"

   "web-service/pkg/io"
   "web-service/pkg/utils"

   "github.com/shirou/gopsutil/host"
	"github.com/julienschmidt/httprouter"
)

/*
var (
        cachedUsersByEmail = map[string]database.User{}
        usersByEmailSync   sync.RWMutex
)
*/

type IndexData struct {
   Version string
   HostName string
   Uptime float64
   Date string
   CPULoad []float64
}

func BasicAuth(h httprouter.Handle) httprouter.Handle {
   
   return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

      username, password, hasAuth := request.BasicAuth()

      io.LogDebug("WEB - BasicAuth", "username: " + username)
      io.LogDebug("WEB - BasicAuth", "password: " + password)

      /*
      usersByEmailSync.RLock()
      user, userFound := cachedUsersByEmail[email]
      usersByEmailSync.RUnlock()
      userMatchesPassword := comparePasswords(user.Password, []byte(password))
      */

      if hasAuth {

         usernameHash := sha256.Sum256([]byte(username))
         passwordHash := sha256.Sum256([]byte(password))
         expectedUsernameHash := sha256.Sum256([]byte("test"))
         expectedPasswordHash := sha256.Sum256([]byte("1010"))

         usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
         passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

         if usernameMatch && passwordMatch {
            h(writer, request, params)
            return
         } else {
            writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
            http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
         }

      } else {
         writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
         http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
      }
   }
}


func Index (writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	timer := time.Now()

   var data IndexData

   data.Version = "0.0.1"
   dt := time.Now()
   data.Date = dt.Format("01-02-2006 15:04:05")
   hostname, err := os.Hostname()
   if err != nil {
      io.LogError("WEB - Index", "error getting hostname")
      os.Exit(1)
   }
   data.HostName = hostname

   utime, err := host.Uptime()
   if err != nil {
      io.LogError("WEB - Index", "error getting uptime")
      os.Exit(1)
   }

   data.Uptime = float64(utime) / 3600
   
   percent := utils.GetCPUsLoad()
   if len(percent) < 0 {
      io.LogError("WEB - InitServer", "cannot have CPU count < 0")
   }
   
   for i := 0; i < len(percent); i++ {
      data.CPULoad = append(data.CPULoad, percent[i])
   }

   tmpl := template.Must(template.ParseFiles("web/html/index.html"))
   _ = tmpl.Execute(writer, data)
   io.LogInfo("INDEX", "Page sent in "+time.Since(timer).String())
}
