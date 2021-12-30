// Package io provides ...
package mesa

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"

   "web-service/pkg/io"
)

var starHistoryName = "history.data"
var binaryHistoryName = "binary_history.data"

var binaryLogDirectory = "LOGS_binary"
var starLogDirectory = "LOGS"
var star1LogDirectory = "LOGS1"
var star2LogDirectory = "LOGS2"


// struct holding info on MESAstar
type MESAstarInfo struct {
   Version int
   Date string
   HistoryName string
   ModelNumber int
   NumZones int
   Mass float64
   LogMdot float64
   Age float64
   CenterH1, CenterHe4 float64
   LogTcntr float64
   NumRetries, NumIters int
   ElapsedTime float64
   EvolState string
}


// struct holding info on MESAbinary
type MESAbinaryInfo struct {
   ModelNumber int
   InitialDonorMass, InitialAccretorMass float64
   InitialPeriod float64
   Age float64
   Star1Mass, Star2Mass float64
   Period float64
   MTCase string
   HistoryName string
   DonorIndex, PointMassIndex int
   RelRLOF1, RelRLOF2 float64
}


type MESAInfo struct {
   ProcId int
   RootDir string
   BinaryFilename string
   Star1Filename string
   Star2Filename string
   BinaryInfo *MESAbinaryInfo
   Star1Info *MESAstarInfo
   Star2Info *MESAstarInfo
   Have2Stars bool
   IsBinaryEvolution bool
}


// get useful information of a MESA run
func (m *MESAInfo) LoadMESAData () error {
   
   // find out if it is a binary or isolated evolution
   m.IsBinaryEvolution = IsBinary(m.RootDir)

   // get LOG names for either single or binary evolutions.
   // in the case of a single evolution, only star1LogName should not be empty
   // for a binary, binaryLogName and star1LogName will not be empty; star2LogName might, if its a
   // star + point-mass simulations
   err := m.getLogNames()
   if err != nil {
      io.LogError("MESA - LoadMESAData", "problem getting LOGS names")
   }

   return nil

}


// return logs names from MESA folder
func (m *MESAInfo) getLogNames () error {

   io.LogInfo("MESA - getLogNames", "searching for MESA LOGS filename(s)")

   binaryLogName := ""
   star1LogName := ""
   star2LogName := ""

   if (m.IsBinaryEvolution) {

      // search for binary output
      // use defaults values defined at beginning of module
      binaryLogName = fmt.Sprintf("%s%s", m.RootDir, binaryHistoryName)
      _, err := os.Stat(binaryLogName)
      if err != nil {
         binaryLogName = fmt.Sprintf("%s%s/%s", m.RootDir, binaryLogDirectory, binaryHistoryName)
         _, err = os.Stat(binaryLogName)
         if err != nil {
            io.LogError("MESA - getLogNames", "cannot find binary LOG output file")
         }
      }
      io.LogInfo("MESA - getLogNames", "found binary output: " + binaryLogName)

      // now look for star 1 data
      star1LogName = fmt.Sprintf("%s%s/%s", m.RootDir, starLogDirectory, starHistoryName)
      _, err = os.Stat(star1LogName)
      if err != nil {
         star1LogName = fmt.Sprintf("%s%s/%s", m.RootDir, star1LogDirectory, starHistoryName)
         _, err = os.Stat(star1LogName)
         if err != nil {
            star1LogName = fmt.Sprintf("%s%s/%s", m.RootDir, star1LogDirectory, "primary_history.data")
            _, err = os.Stat(star1LogName)
            if err != nil {
               star1LogName = fmt.Sprintf("%s%s/%s", m.RootDir, "LOGS_companion", starHistoryName)
               _, err = os.Stat(star1LogName)
               if err != nil {
                  io.LogError("MESA - getLogNames", "cannot find star 1 LOG output file")
               }
            }
         }
      }
      io.LogInfo("MESA - getLogNames", "found star 1 output: " + star1LogName)

      // now look for star 2 data (though not always found if doing star + point-mass)
      star2LogName = fmt.Sprintf("%s%s/%s", m.RootDir, star2LogDirectory, starHistoryName)
      _, err = os.Stat(star2LogName)
      if err != nil {
         star2LogName = fmt.Sprintf("%s%s/%s", m.RootDir, star2LogDirectory, "secondary_history.data")
         _, err = os.Stat(star2LogName)
         if err != nil {
            io.LogInfo("MESA - getLogNames", "cannot find star 2 LOG output file. maybe doing star + point-mass evolution")
            star2LogName = ""
         } else {
            io.LogInfo("MESA - getLogNames", "found star 2 output: " + star2LogName)
         }
      } else {
         io.LogInfo("MESA - getLogNames", "found star 2 output: " + star2LogName)
      }

   } else {

      // only need to search for star1LogName
      star1LogName = fmt.Sprintf("%s%s/%s", m.RootDir, starLogDirectory, starHistoryName)

      _, err := os.Stat(star1LogName)
      if err != nil {
         io.LogError("MESA - getLogNames", "cannot find star LOG output file of single evolution")
      }

      io.LogInfo("MESA - getLogNames", "found single evolution output: " + star1LogName)
   }

   // update struct with LOG filenames
   m.BinaryFilename = binaryLogName
   m.Star1Filename = star1LogName
   m.Star2Filename = star2LogName

   return nil

}


// get useful information for the summary of a MESAstar run
func (s *MESAstarInfo) LoadMESAstarData () error {
   
   // MESA specific row numbers for header names & values in history output
   nr_header_names := 2
   nr_header_values := 3
   nr_column_names := 6

    // open star file
   fstar, err := os.Open(s.HistoryName)
   if err != nil {
      io.LogError("MESA - loadMESAstarData", "problem opening star data file")
      return nil
   }
   defer fstar.Close()

   // scan star file
   scanner := bufio.NewScanner(fstar)

   // arrays holding star header names & values
   var header_names []string
   var header_values []string
   var header_names_found, header_values_found bool
   var column_names []string
   var column_values []string
   var column_names_found bool
   lineCount := 0

   for scanner.Scan() {

      lineCount++

      // get header names
      if lineCount == nr_header_names {
         header_names = strings.Fields(scanner.Text())
         header_names_found = true
      }

      // find header values
      if lineCount == nr_header_values {
         header_values = strings.Fields(scanner.Text())
         header_values_found = true
      }

      if (header_names_found && header_values_found) {
         for k, name := range header_names {
            if name == "version_number" {
               i, err := strconv.Atoi(strings.Split(header_values[k], "\"")[1])
               // handle error
               if err != nil {
                  fmt.Println(err)
                  return nil
               }
               s.Version = i
            }
            if name == "date" {s.Date = header_values[k]}
         }
      }

      // get header names
      if lineCount == nr_column_names {
         column_names = strings.Fields(scanner.Text())
         column_names_found = true
      }

      if column_names_found {
         break
      }

   }

   if (column_names_found) {
      column_values = strings.Fields(GetLastLineWithSeek(s.HistoryName))

      for k, name := range column_names {
         val := column_values[k]
         if name == "model_number" {
            i, err := strconv.Atoi(val)
            // handle error
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.ModelNumber = i
         }
         if name == "num_zones" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.NumZones = i
         }
         if name == "star_mass" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.Mass = i
         }
         if name == "log_abs_mdot" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.LogMdot = i
         }
         if name == "star_age" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.Age = i
         }
         if name == "center_h1" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.CenterH1 = i
         }
         if name == "center_he4" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.CenterHe4 = i
         }
         if name == "log_center_T" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.LogTcntr = i
         }
         if name == "log_cntr_T" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.LogTcntr = i
         }
         if name == "num_retries" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.NumRetries = i
         }
         if name == "num_iters" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.NumIters = i
         }
         if name == "elapsed_time" {
            i, err := strconv.ParseFloat(val, 64)  // i is in sec
            if err != nil {
               fmt.Println(err)
               return nil
            }
            s.ElapsedTime = i / 60 // from sec to min
         }

      }

   }

   s.EvolState = SetEvolutionaryStage(s.Mass, s.CenterH1, s.CenterHe4, s.LogTcntr)

   if err := scanner.Err(); err != nil {
      io.LogError("MESA - loadMESAstarData", "problem while scanning binary data file")
      return nil
   }

   return nil

}


// get useful information for the summary of a MESAbinary run
func (b *MESAbinaryInfo) LoadMESAbinaryData () error {

   // MESA specific row numbers for header names & values in history output
   nr_header_names := 2
   nr_header_values := 3
   nr_column_names := 6

   // open binary file
   fbinary, err := os.Open(b.HistoryName)
   if err != nil {
      io.LogError("MESA - loadMESAbinaryData", "problem opening binary data file")
      return nil
   }
   defer fbinary.Close()

   // scan star file
   scanner := bufio.NewScanner(fbinary)

   // arrays holding header, column names and values
   var header_names []string
   var header_values []string
   var header_names_found, header_values_found bool
   var column_names []string
   var column_values []string
   var column_names_found bool
   lineCount := 0

   for scanner.Scan() {

      lineCount++

      // get header names
      if lineCount == nr_header_names {
         header_names = strings.Fields(scanner.Text())
         header_names_found = true
      }

      // find header values
      if lineCount == nr_header_values {
         header_values = strings.Fields(scanner.Text())
         header_values_found = true
      }

      if (header_names_found && header_values_found) {
         for k, name := range header_names {
            val := header_values[k]
            if name == "initial_don_mass" {
               i, err := strconv.ParseFloat(val, 64)  // i is in Msun
               if err != nil {
                  fmt.Println(err)
                  return nil
               }
               b.InitialDonorMass = i
            }
            if name == "initial_acc_mass" {
               i, err := strconv.ParseFloat(val, 64)  // i is in Msun
               if err != nil {
                  fmt.Println(err)
                  return nil
               }
               b.InitialAccretorMass = i
            }
            if name == "initial_period_days" {
               i, err := strconv.ParseFloat(val, 64)  // i is in days
               if err != nil {
                  fmt.Println(err)
                  return nil
               }
               b.InitialPeriod = i
            }
         }
      }

      // get column names
      if lineCount == nr_column_names {
         column_names = strings.Fields(scanner.Text())
         column_names_found = true
      }

      // once row with column names is found, just exit loop
      if column_names_found {
         break
      }

   }

   // load last row of file and loop through it to match each column value with name
   if (column_names_found) {
      column_values = strings.Fields(GetLastLineWithSeek(b.HistoryName))

      for k, name := range column_names {
         val := column_values[k]
         if name == "model_number" {
            i, err := strconv.Atoi(val)
            // handle error
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.ModelNumber = i
         }
         if name == "age" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.Age = i
         }
         if name == "period_days" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.Period = i
         }
         if name == "star_1_mass" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.Star1Mass = i
         }
         if name == "star_2_mass" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.Star2Mass = i
         }
         if name == "donor_index" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.DonorIndex = i
         }
         if name == "point_mass_index" {
            i, err := strconv.Atoi(val)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.PointMassIndex = i
         }
         if name == "rl_relative_overflow_1" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.RelRLOF1 = i
         }
         if name == "rl_relative_overflow_2" {
            i, err := strconv.ParseFloat(val, 64)
            if err != nil {
               fmt.Println(err)
               return nil
            }
            b.RelRLOF2 = i
         }
      }
   }

   if err := scanner.Err(); err != nil {
      io.LogError("MESA - loadMESAbinaryData", "problem while scanning binary data file")
      return nil
   }

   return nil

}
