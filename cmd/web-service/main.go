
package main

import (
   
   "web-service/pkg/io"
   "web-service/pkg/web"

)


const serviceName = "Computer status web service"
const serviceDescription = "Service that monitors the status of a computer doing scientific computations"
const version = "2021.12.2.1"


func main() {

   io.LogInfo("MAIN - main.go - main", serviceName)


   // Struct with information of the service
   Info := new(web.ServiceInfo)
   Info.ServiceName = serviceName
   Info.ServiceDescription = serviceDescription
   Info.Version = version


   // Initialize service
   io.LogInfo("MAIN - main.go - main", "initializing service")
   web.InitServer(Info)

}
