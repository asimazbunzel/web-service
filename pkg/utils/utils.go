// Package utils provides ...
package utils

import (
   "time"
   "strconv"
   
   "web-service/pkg/io"

   "github.com/shirou/gopsutil/cpu"
)

func GetCPUsLoad () []float64 {
   
   io.LogInfo("UTILS - utils.go - GetCPUsLoad", "loading CPU percentage")
   
   percents, _ := cpu.Percent(time.Second, true)


   for i := 0; i < len(percents); i++ {

      cpu_val := strconv.Itoa(i)
      percent_val := strconv.FormatFloat(percents[i], 'f', -1, 64)
      dbg_string := "CPU " + cpu_val + ": " + percent_val
   
      io.LogDebug("UTILS - utils.go - GetCPUsLoad", dbg_string)

   }

   return percents

}
