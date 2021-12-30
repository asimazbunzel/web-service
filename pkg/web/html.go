package web

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	// "fmt"

	"html/template"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"time"

	"web-service/pkg/io"
	"web-service/pkg/mesa"
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

func (Data *IndexData) GetIndexData () error {

   Data.Version = "0.0.1"

   // get CPU info
   stats, err := cpu.Info()
    if err != nil {
      io.LogError("HTML - GetIndexData", "error getting CPU info")
      os.Exit(1)
    }
   
   // get time from server
   dt := time.Now()
   Data.Date = dt.Format("01-02-2006 15:04:05")

   // hostname
   hostname, err := os.Hostname()
   if err != nil {
      io.LogError("HTML - GetIndexData", "error getting hostname")
      os.Exit(1)
   }
   Data.HostName = hostname

   // info on the user running the server
   userinfo, err := user.Current()
   if err != nil {
      io.LogError("HTML - GetIndexData", "error getting username")
      os.Exit(1)
   }
   Data.UserName = userinfo.Username

   // time server has been up
   utime, err := host.Uptime()
   if err != nil {
      io.LogError("HTML - GetIndexData", "error getting uptime")
      os.Exit(1)
   }
   ftime := float64(utime / 3600)
   Data.Uptime = strconv.FormatFloat(ftime, 'f', 2, 64)

   // number of cores and type
   Data.NCores = strconv.Itoa(len(stats))
   Data.CPUInfo = stats[0].ModelName

   // CPU load
   percent := utils.GetCPUsLoad()
   if len(percent) < 0 {
      io.LogError("HTML - GetIndexData", "cannot have CPU count < 0")
   }

   Data.NCoresGT4 = false
   if len(percent) > 4 {
      Data.NCoresGT4 = true
   }
   
   for i := 0; i < len(percent); i++ {
      str := strconv.FormatFloat(percent[i], 'f', 2, 64)
      Data.CPULoad = append(Data.CPULoad, str)
      Data.CPUno = append(Data.CPUno, i)
   }

   return nil
}

func BasicAuth(h httprouter.Handle) httprouter.Handle {
   
   return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

      username, password, hasAuth := request.BasicAuth()

      io.LogDebug("HTML - BasicAuth", "username: " + username)
      io.LogDebug("HTML - BasicAuth", "password: " + password)

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


// index.html serving func
func Index (writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

   // start counting time until serve files
   timer := time.Now()

   // load info for the index page
   data := new(IndexData)
   data.GetIndexData()

   // serve index.html
   tmpl := template.Must(template.ParseFiles("web/html/index.html"))
   _ = tmpl.Execute(writer, data)
   io.LogInfo("HTML - Index", "Page sent in "+time.Since(timer).String())

}


// dashboard.html serving func
func Dashboard (writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

   // start counting time until serve files
   timer := time.Now()

   // load info for the index page
   data := new(IndexData)
   data.GetIndexData()

   tmpl := template.Must(template.ParseFiles("web/html/dashboard.html"))
   _ = tmpl.Execute(writer, data)
   io.LogInfo("HTML - Dashboard", "Page sent in "+time.Since(timer).String())

}


// mesa.html serving func
func MESAhtml (writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {

   // start counting time until serve files
   timer := time.Now()

   // set struct which has tha ability to find the process running a MESA executable
   mesaProc := new(utils.MESAprocess)
   mesaProc.WalkProc()
   
   // get abs path where run is being done, in case there is no simulation, this is simply and
   // string
   if mesaProc.Id > 0 {
      mesaProc.GetAbsPath()
   }

   // set struct with MESA info, which will later be connected to html file via Templates
   mesaInfo := new(mesa.MESAInfo)
   bInfo := new(mesa.MESAbinaryInfo)
   star1Info := new(mesa.MESAstarInfo)
   star2Info := new(mesa.MESAstarInfo)

   // set some defaults
   mesaInfo.ProcId = mesaProc.Id
   mesaInfo.RootDir = mesaProc.Loc

   // load all the info on the binary run
   if mesaProc.Id > 0 {

      err := mesaInfo.LoadMESAData()

      // if problems while loading stuff, just set the ProcId to a reserve value so that the html
      // will warn about it
      if err != nil {
         io.LogError("WEB - MESAhtml", "problem loading MESA data")
         mesaInfo.ProcId = -98
      }

      // some defaults for MESAbinary search
      bInfo.HistoryName = mesaInfo.BinaryFilename
      bInfo.MTCase = "none"

      // load MESAbinary info
      err = bInfo.LoadMESAbinaryData()

      // again, if problems were found, give some warning in the html
      if err != nil {
         io.LogError("WEB - MESAhtml", "problem loading MESAbinary data")
         mesaInfo.ProcId = -97
      }

      // star1 defaults
      star1Info.HistoryName = mesaInfo.Star1Filename

      // load MESAstar data for star1
      err = star1Info.LoadMESAstarData()
      if err != nil {
         io.LogError("WEB - MESAhtml", "problem loading MESAstar data for star 1")
         mesaInfo.ProcId = -96
      }

      // by default, Have2Stars is false. change accordingly to the value of point_mass_index column
      if bInfo.PointMassIndex == 0 {
         mesaInfo.Have2Stars = true
      }

      // star2 defaults
      star2Info.HistoryName = mesaInfo.Star2Filename

      // load MESAstar data for star2
      err = star2Info.LoadMESAstarData()
      if err != nil {
         io.LogError("WEB - MESAhtml", "problem loading MESAstar data for star 2")
         mesaInfo.ProcId = -95
      }

      // also. after having info on (possibly) both stars, check for a MT phase
      if bInfo.DonorIndex == 1 {
         bInfo.MTCase = mesa.SetMTCase(bInfo.RelRLOF1, star1Info.EvolState)
      } else {
         bInfo.MTCase = mesa.SetMTCase(bInfo.RelRLOF2, star2Info.EvolState)
      }

      // finally, store useful information inside the MESAInfo struct
      mesaInfo.BinaryInfo = bInfo
      mesaInfo.Star1Info  = star1Info
      mesaInfo.Star2Info  = star2Info

   } else {

      // this is to get the correct message in the html page
      mesaInfo.ProcId = -99

   }

   // server html
   tmpl := template.Must(template.ParseFiles("web/html/mesa.html"))
   _ = tmpl.Execute(writer, mesaInfo)
   io.LogInfo("HTML - MESAhtml", "Page sent in "+time.Since(timer).String())

}
