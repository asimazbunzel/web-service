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


// Mp structure holds the info 
type MESAprocess struct {
   ExecName string
   Id int
   Loc string
}

// wrapper function for searching MESA simulation process
func (M *MESAprocess) WalkProc () {

   err := filepath.Walk("/proc", M.FindMESAProcess)
   io.LogInfo("PROCESS - Walkproc", "exiting filepath.Walk")

   if err != nil {
      io.LogDebug("PROCESS - WalkProc", M.ExecName)
      io.LogDebug("PROCESS - WalkProc", strconv.Itoa(M.Id))
      return
   } else {
      io.LogError("PROCESS - WalkProc", "could not find MESA process")
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
            io.LogInfo("PROCESS - FindMESAProcess", "problem converting Atoi")
            return nil
         }

         f, err := ioutil.ReadFile(path)
         if err != nil {
            io.LogInfo("PROCESS - FindMESAProcess", "problem reading status file")
            return nil
         }

         // Extract the process name from within the first line in the buffer
         name := string(f[6:bytes.IndexByte(f, '\n')])
         if name == M.ExecName {

            io.LogInfo("PROCESS - FindMESAProcess", "found MESA process")
            io.LogInfo("PROCESS - FindMESAProcess", strconv.Itoa(pid))

            M.ExecName = name
            M.Id = pid
            M.Loc = "/proc/" + strconv.Itoa(pid)

            // proc, err := os.FindProcess(pid)
            // if err != nil {
                // log.Println(err)
            // }

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
   M.Loc = strings.Replace(last, "\n", "", -1)

   io.LogDebug("PROCESS - GetAbsPath", "found AbsPath on " + M.Loc)

}