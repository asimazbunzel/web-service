package mesa

import (
   "fmt"
   "os"

   "web-service/pkg/io"
)


// function to find out if path contains a single or binary evolution
// to check if this is a binary simulation, look for the MESAbinary output
func IsBinary (path string) bool {

   io.LogInfo("IO - IsBinary", "searching for binary evolution")

   // name of the MESAbinary output
   binaryFile := fmt.Sprintf("%s%s", path, binaryHistoryName)

   io.LogInfo("IO - IsBinary", "looking for file " + binaryFile)
   _, err := os.Stat(binaryFile)
   if err != nil {

      io.LogInfo("IO - IsBinary", binaryFile + " file not found. now searching inside " + binaryLogDirectory + " folder")

      binaryFile := fmt.Sprintf("%s%s/%s", path, binaryLogDirectory, binaryHistoryName)

      _, err := os.Stat(binaryFile)
      if err != nil {

         io.LogInfo("IO - IsBinary", "binary logs not found. single evolution assumed")
         return false

      } else {

         io.LogInfo("IO - IsBinary", "found binary log: " + binaryFile + ". binary evolution assumed")
         return true

      }
   } else {

      io.LogInfo("IO - IsBinary", "found binary logs. binary evolution assumed")
      return true

   }
}

