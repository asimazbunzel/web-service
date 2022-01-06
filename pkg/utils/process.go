// Package utils provides ...
package utils

import (
   "bytes"
   "errors"
   "io/ioutil"
   "path/filepath"
   "strconv"
   "strings"
   "os"
   "os/exec"
   
   "web-service/pkg/io"
)


// possible names of MESA executable binaries
var execNames = []string{"star", "binary", "bin2dco"}


// Mp structure holds the info 
type MESAprocess struct {
   ExecName string
   Id int
   Loc string
}


// wrapper function for searching MESA simulation process
func (M *MESAprocess) WalkProc () {

   for k, name := range execNames {

      M.ExecName = name

      io.LogInfo("UTILS - process.go - Walkproc", "doing filepath.Walk for name: " + name)
      err := filepath.Walk("/proc", M.FindMESAProcess)

      if err != nil {

         io.LogDebug("UTILS - process.go - WalkProc", M.ExecName)
         io.LogDebug("UTILS - process.go - WalkProc", strconv.Itoa(M.Id))

         return

      } else {

         io.LogError("UTILS - process.go - WalkProc", "could not find MESA process for name: " + name)

         if (k == len(execNames)) {
            return
         }

      }

   }

}


// search process with MESA simulation
// idea from this post:
// https://stackoverflow.com/questions/41060457/golang-kill-process-by-name
func (M *MESAprocess) FindMESAProcess (path string, info os.FileInfo, err error) error {

   if err != nil {
      return nil
   }

   if strings.Count(path, "/") == 3 {

      if strings.Contains(path, "/status") {

         pid, err := strconv.Atoi(path[6:strings.LastIndex(path, "/")])
         if err != nil {
            io.LogInfo("UTILS - process.go - FindMESAProcess", "problem converting Atoi")
            return nil
         }

         f, err := ioutil.ReadFile(path)
         if err != nil {
            io.LogInfo("UTILS - process.go - FindMESAProcess", "problem reading status file")
            return nil
         }

         // Extract the process name from within the first line in the buffer
         name := string(f[6:bytes.IndexByte(f, '\n')])
         if name == M.ExecName {

            io.LogInfo("UTILS - process.go - FindMESAProcess", "found MESA process")
            io.LogInfo("UTILS - process.go - FindMESAProcess", strconv.Itoa(pid))

            // asign PID and abs path of process
            M.Id = pid
            M.Loc = "/proc/" + strconv.Itoa(pid)

            // Let's return a fake error to abort the walk through the
            // rest of the /proc directory tree
            return errors.New("found MESA process")
         }

      }
   }

   return nil

}


func (M *MESAprocess) GetAbsPath () {

   cmd := exec.Command("ls", "-l", "/proc/" + strconv.Itoa(M.Id) + "/exe")
   out, err := cmd.CombinedOutput()
   if err != nil {
      io.LogError("PROCESS - GetAbsPath", "problem running cmd")
   }

   // get only last element of slice
   res1 := strings.Split(string(out), " ")
   last := res1[len(res1)-1]

   // strip from newline char
   last2 := strings.Replace(last, "\n", "", -1)
   last2 = last2[:len(last2) - len(M.ExecName)]
   M.Loc = last2

   io.LogDebug("PROCESS - GetAbsPath", "found AbsPath on " + M.Loc)

}
