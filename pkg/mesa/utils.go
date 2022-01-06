package mesa

import (
   "fmt"
   "os"

   "web-service/pkg/io"
)


// function to find out if path contains a single or binary evolution
// to check if this is a binary simulation, look for the MESAbinary output
func IsBinary (path string) bool {

   io.LogInfo("MESA - utils.go - IsBinary", "searching for binary evolution")

   // name of the MESAbinary output
   binaryFile := fmt.Sprintf("%s%s", path, binaryHistoryName)

   io.LogInfo("MESA - utils.go - IsBinary", "looking for file " + binaryFile)
   _, err := os.Stat(binaryFile)
   if err != nil {

      io.LogInfo("MESA - utils.go - IsBinary", binaryFile + " file not found. now searching inside " + binaryLogDirectory + " folder")

      binaryFile := fmt.Sprintf("%s%s/%s", path, binaryLogDirectory, binaryHistoryName)

      _, err := os.Stat(binaryFile)
      if err != nil {

         io.LogInfo("MESA - utils.go - IsBinary", "binary logs not found. single evolution assumed")
         return false

      } else {

         io.LogInfo("MESA - utils.go - IsBinary", "found binary log: " + binaryFile + ". binary evolution assumed")
         return true

      }
   } else {

      io.LogInfo("MESA - utils.go - IsBinary", "found binary logs. binary evolution assumed")
      return true

   }
}


// define phase of evolution for a star based on abundances and central temperature
func SetEvolutionaryStage(mass float64, center_h1 float64, center_he4 float64, log_T_cntr float64) string {

   // H threshold
   h1_threshold := 1e-5

   // He$ threshold
   he4_threshold := 1e-4

   // log T for He ignition
   log_he_temp := 7.8

   // chandra mass
   chandra_mass := 1.4

   // Main sequence star
   if (center_h1 > h1_threshold) {
      return "MS star"
   }

   // Hertzsprung gap star or WD
   if (log_T_cntr < log_he_temp) {

      // depending on the center He4 value, either HG star of becoming WD or ECSN
      if (center_he4 < he4_threshold) {

         // chandra mass separates between WD & ECSN
         if (mass < chandra_mass) {
            return "WD"
         } else {
            return "He depleted star, possible EC SN"
         }

      } else {

         // HG star
         return "HG star"
      }
   }

   // core He burning star
   if (center_he4 > he4_threshold) {

      return "CHeB star"

   } else {

      // already past core He depletion
      // chandra mass separates between WD & ECSN
      if (mass < chandra_mass) {
         return "WD"
      } else {
         return "He depleted star"
      }
   }

}


// check for MT case in a binary
func SetMTCase(rel_rlof float64, stage string) string {

   // if no RLOF, no MT phase
   if (rel_rlof < 0) {

      return "No MT (R < RL)"

   } else {

      if (stage == "MS star") {return "Case A"}

      if (stage == "HG star") {return "early Case B"}

      if (stage == "CheB star") {return "Case B"}

      return "Case C"

   }

}
