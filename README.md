# ARImport

### Building
```golang
>> go build .\server\main.go
>> ./main.exe
```

### Packages
*   output
*   exe
#### Generate docs for package (examples)
#### CMD :
```golang
>> go doc output
Package output contains interface for logging errors to text file.

var Log *log.Logger
func Close()
func Pf(fS string, err error, fatal bool)
```
```golang
>> go doc output Log
var Log *log.Logger
    Log is a pointer to the log.Logger struct.
```
#### HTTP :
```golang
>> godoc -http=:8002
```
*   Navigate to Packages -> Third party -> ajagnic -> ARImport
