package web

import (
	"crypto/sha256"
	"crypto/subtle"
	"html/template"
	"net/http"
	"os"
   "os/user"
	"strconv"
	"time"

	"web-service/pkg/io"
	"web-service/pkg/utils"

	"github.com/julienschmidt/httprouter"
   "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)


type IndexData struct {
   Version string
   HostName string
   UserName string
   Uptime string
   NCores string
   NCoresGT4 bool
   CPUInfo string
   Date string
   CPULoad []string
   CPUno []int
}

func BasicAuth(h httprouter.Handle) httprouter.Handle {
   
   return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

      username, password, hasAuth := request.BasicAuth()

      io.LogDebug("WEB - BasicAuth", "username: " + username)
      io.LogDebug("WEB - BasicAuth", "password: " + password)

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

   stats, err := cpu.Info()
    if err != nil {
      io.LogError("WEB - Index", "error getting CPU info")
      os.Exit(1)
    }

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
   
   userinfo, err := user.Current()
   if err != nil {
      io.LogError("WEB - Index", "error getting username")
      os.Exit(1)
   }
   data.UserName = userinfo.Username

   utime, err := host.Uptime()
   if err != nil {
      io.LogError("WEB - Index", "error getting uptime")
      os.Exit(1)
   }
   ftime := float64(utime / 3600)
   data.Uptime = strconv.FormatFloat(ftime, 'f', 2, 64)
    
   data.NCores = strconv.Itoa(len(stats))
   data.CPUInfo = stats[0].ModelName
   
   percent := utils.GetCPUsLoad()
   if len(percent) < 0 {
      io.LogError("WEB - InitServer", "cannot have CPU count < 0")
   }

   data.NCoresGT4 = false
   if len(percent) > 4 {
      data.NCoresGT4 = true
   }
   
   for i := 0; i < len(percent); i++ {
      str := strconv.FormatFloat(percent[i], 'f', 2, 64)
      data.CPULoad = append(data.CPULoad, str)
      data.CPUno = append(data.CPUno, i)
   }

   tmpl := template.Must(template.ParseFiles("web/html/index.html"))
   _ = tmpl.Execute(writer, data)
   io.LogInfo("INDEX", "Page sent in "+time.Since(timer).String())
}
