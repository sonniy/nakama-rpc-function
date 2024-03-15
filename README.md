Nakama rpc function test project
===
This project represent posibility of Nakama game server run custom logic using rpc call
## Goals
* Сreate a rpc function
* Сover it by tests
* Should be launched simply using a Docker

## Story
- rpc function read a file from the disk by template __path=%type/%version.json__ (e.g. __"core/1.0.0.json"__)
- save information to database using template __%type/%version__ as key and store __content__ of the file like value
- If hashes are not equal, then content will be null.
- If file doesn't exist, then return error.
- defaults parameter: type=core, version=1.0.0, hash=null
  


## How to run
- download a project using GitHub
- go to project directory
- run command `docker-compose up --build nakama`

## How to check results 
### easy way
- make POST request to http://127.0.0.1:7351/v2/console/api/endpoints/rpc/versionchecker via terminal
`curl "http://127.0.0.1:7350/v2/rpc/versionchecker?http_key=defaulthttpkey" \
	-d '"{\"type\": \"core\",\"version\": \"1.0.0\",\"hash\": \"c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2\"}"' \
	-H 'Content-Type: application/json' \
	-H 'Accept: application/json'`
  
### hard way
- go to http://127.0.0.1:7351/
- enter username: `admin` and password: `password`
<img width="563" alt="image" src="https://github.com/sonniy/nakama-rpc-function/assets/564889/db312ace-e011-49e4-a769-c28035e7f18e">

- go to API explorer and select `versionchecker`
<img width="570" alt="image" src="https://github.com/sonniy/nakama-rpc-function/assets/564889/5994eb96-1886-43f5-b34b-61b43cb4ba9d">

- put `{}` in to `Request Body` fild and press `SEND REQUEST`
<img width="1106" alt="image" src="https://github.com/sonniy/nakama-rpc-function/assets/564889/b03fd332-f1ac-4e0e-aa2d-ebf252c9edd8">

- congratulation everything is work :)
- you also can check more dificult request:
`{
  "type": "core",
  "version": "1.0.0",
  "hash": "c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2"
}`
<img width="1108" alt="image" src="https://github.com/sonniy/nakama-rpc-function/assets/564889/a64f7f16-3aca-4ef2-9b6d-ad17c5c90206">

## Reqest and responce examples
### success, all parameters are presented, file exist
#### request
```
curl "http://127.0.0.1:7350/v2/rpc/versionchecker?http_key=defaulthttpkey" \
	-d '"{\"type\": \"core\",\"version\": \"1.0.0\",\"hash\": \"c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2\"}"' \
	-H 'Content-Type: application/json' \
	-H 'Accept: application/json'
```
#### response 
```json
{ "payload":"{\"type\":\"core\",\"version\":\"1.0.0\",\"hash\":\"c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2\",\"content\":\"nakama should read this file\"}"}
```
### half success, type/version - presented, hash - no
#### request
```
curl "http://127.0.0.1:7350/v2/rpc/versionchecker?http_key=defaulthttpkey" \
	-d '"{\"type\": \"core\",\"version\": \"1.0.0\"}"' \
	-H 'Content-Type: application/json' \
	-H 'Accept: application/json'
```
#### response, content is empty
```json
{"payload":"{\"type\":\"core\",\"version\":\"1.0.0\",\"hash\":\"c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2\",\"content\":\"\"}"
```
### half success, no param, params will be defaults
#### request
```
curl "http://127.0.0.1:7350/v2/rpc/versionchecker?http_key=defaulthttpkey" \
	-d '"{}"' \
	-H 'Content-Type: application/json' \
	-H 'Accept: application/json'
```
#### response, content is empty
```json
{"payload":"{\"type\":\"core\",\"version\":\"1.0.0\",\"hash\":\"c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2\",\"content\":\"\"}"
```
### unsuccess, wrong type or version
#### request
```
curl "http://127.0.0.1:7350/v2/rpc/versionchecker?http_key=defaulthttpkey" \
	-d '"{\"type\":\"rpc\",\"version\":\"1.0.99\"}"' \
	-H 'Content-Type: application/json' \
	-H 'Accept: application/json'
```
#### response
```json
{"code":13,"error":{},"message":"file not found: stat rpc/1.0.99.json: no such file or directory"}
```
# What can be done better
- I would find out the business purpose of this logic, becaose return hash when incoming hash is different looks like security issue, probably here we can find more issue
- I used default user for storing data, it looks no so good, but it general way to save something, so I'd dig deeper and would create my own table(nakama not recomended use other DB instead of their own, built-in one, like coocroach or postgress)
