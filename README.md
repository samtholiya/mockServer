# Mock Server [![Build Status](https://travis-ci.com/samtholiya/mockServer.svg?branch=master)](https://travis-ci.com/samtholiya/mockServer)

Mock Server can be used by people who want to test their client by mocking the backend responses instead of changing the backend
The only thing you need to do is configure the config.yaml and you are ready to run the server which gives responses as per requirement.

Mock Server has two modes to help you actually get started with mocking your desired server.

### Proxy Mode
This mode can be used when you want to record the responses in the config that you will be using while running the mock server without actual server where you are expecting your responses from.

You can start server in proxy mode using
```
./mockServer -proxy -host="https://urloftheactualserver"
```
This will create a file proxy_generated.yml which can be used a starting point for configuring your mock server.
Once you have started the mock server in proxy mode you can see all the logs related to endpoints in the log. Once you have recorded your required request and responses. You can shut down the mock server and edit the proxy_generated.yml for custom responses.

### Mock Mode
This mode can be used to replace your orignal server with the mock one.
Please look at the the wiki for detail on how to configure mock server with config.yml
Your can start server in mock mode
```
./mockServer
```

To build the repo you can use considering you have **go version 1.13.x**

```
go build .
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
