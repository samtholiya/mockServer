## Mock Server [![Build Status](https://travis-ci.com/samtholiya/mockServer.svg?branch=master)](https://travis-ci.com/samtholiya/mockServer)

Mock Server can be used by people who want to test their client by mocking the backend responses instead of changing the backend
The only thing you need to do is configure the config.yaml and you are ready to run the server which gives responses as per requirement.

To build the repo you can use considering you have **go version 1.13.x**

```
go build .
```

To run the mock server in proxy mode you can use
```
./mockServer -proxy -host="https://urloftheactualserver"
```

Following are the environment variables in Mock Server <br/>

|Environment Variable|Default Value|Description|
| --- | --- | --- |
|MOCK_CONFIG|./config.yaml   |  This variable is used to configure the path of the yaml which contains the routes to be mocked with details |
|PROXY_GENERATED_CONFIG  | ./proxy_generated.yaml | This variable is used to configure the path of the file where the proxy mode will write the config it reads from the yaml. |

Following are the parameters in Mock Server <br/>

|Parameters |Default Value |Description |
| --- | --- | --- |
| proxy | false | Used to run mock server as a proxy server and to generate config file |
| host  | https://httpbin.org | Used in proxy mode to set the host to redirect the endpoints |
| port  | 3000 | Used to change the port number of Mock Server|
| debug | false | Used to debug mock server |
