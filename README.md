# memdb
a minimalistic key/value store implemented in Go.

**NOTE:** I wrote these libraries for learning purposes. It may not be completely thought out and error free. Use at Your Own Risk.

---

inspired by redis, this project provides a key/value data store intended to hold a small amount of data within memory.<br>
!!! No intended to be used for large data storage !!!

### How it works
A tcp server listening and accepting connection on port 8000 per default. The client can send a variety of command the follows this pattern
```
  COMMAND KEY VALUE TTL VALUE
```
Available commands:
* SET key value [TTL value]\n
* GET key\n
* REMOVE key\n
* INC key value\n
* DEC key value\n
* NEGATE key\n
* LAPPEND key value\n
* LREMOVE key value\n

Responses follows this pattern:
```
Type Length\r\n
value
```
Response types are special characters identifiers:<br>
    + string<br>
    : integer<br>
    # boolean<br>
    * list<br>
    - error

value can be a string, integer, boolean or a list

#### examples:
###### SET command
```bash
SET key value
Ok
```
###### GET command
```bash
GET example
+5
value
```
###### SET command for list type
```bash
SET example [2,3,4,[test,5,7]]
Ok
```
###### GET command for list type
```bash
*15
[2,3,4,[test,5,7]]
```
### Install
```bash
go install hgithub.com/Wa4h1h/memdb/tree/main/cmd/server@latest
```
---

### Todos
- [X] TCP server
- [X] Command evaluator
- [X] Support for:
  - [X] SET,GET
  - [X] INC,DEC
  - [X] NEGATE
  - [X] LAPPEND,LREMOVE
- [ ] Tests
- [ ] Docker deployment
- [ ] Golang client to interact with db
  - [ ] Implementation
  - [ ] Tests