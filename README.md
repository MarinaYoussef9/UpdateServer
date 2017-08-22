![alt text](https://travis-ci.org/xmarcoied/UpdateServer.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmarcoied/updateserver)](https://goreportcard.com/report/github.com/xmarcoied/updateserver)
# UpdateServer

An UpdateServer written in golang to ship releases and manage update requests


## Installing
Nothing fancy, Just use the golang ```go get```

```
go get code.videolan.org/GSoC2017/Marco/UpdateServer
```

## Configuration
  The configuration comes in JSON format, you either edit the default config.json attached or give the path of the configuration file through flags 
  Example : 
  ```
  {
    "psqlinfo": {
        "psqlhost": "host" , 
        "psqlname" : "marcoied" ,
        "psqluser" : "postgres" , 
        "psqlpassword" : "postgres" 
    }
  }
  
  ```

## Flags
  The UpdateServer uses ```-port``` flag to customize port web application will be running on (8080 default)
  and ```-config``` to add a path for the configuration file 
  
## Usage
  - ```cd $GOPATH/src/code.videolan.org/GSoC2017/Marco/UpdateServer```
  
  - ```go build``` to build a binary UpdateServer

  - ```./UpdateServer``` to run the Server
  
      Optional flags: 
      ```
      ./UpdateServer -port 80
      ./UpdateServer -config $HOME/config.json
      ```
      Default admin authentication : username:admin , password:admin
      
  -  ```<host>/admin/dashboard/channels``` to add new channel with public and private keys

  - ```<host>/admin/dashboard/releaese```  to add new release

       insert the needed fields and choose a channel and you can add rules against it later .
       
       Sign the status given via this commond :
       
       ```$ echo "<release metadata/json to sign here>" | gpg --default-key <fingerprint of the public key> --clearsign```
     
  - add ```<host>/u/update``` to you client with a querystring for the update_request paramaters.
    
    Example : ```<host>/u/update?os=Win&os_ver=7&os_arch=64&product_ver=2.1.4```
  
  - Monitor the incoming update_requests through ```<host>/u/releases```
  
  
