# ARImport

### Building
```console
>> go build .\src\server\main.go
>> .\main.exe
```

### Packages
*   output
*   exe
*   scheduler

#### Generate docs for package (examples)
#### CMD :
```console
>> go doc output
Package output contains interface for logging errors to text file.

var Log *log.Logger
func Close()
func Pf(fS string, err error, fatal bool)
```

```console
>> go doc output Log
var Log *log.Logger
    Log is a pointer to the log.Logger struct.
```

#### HTTP :
```console
>> godoc -http=:8002
```
*   Navigate to Packages -> Third party -> ajagnic -> ARImport
