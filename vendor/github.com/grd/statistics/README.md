# Statistics

Pure [Go](http://www.golang.org) implementation of the [GSL Statistics library](http://www.gnu.org/software/gsl/manual/html_node/Statistics.html).

For the API overview see [Godoc](http://godoc.org/github.com/grd/statistics).

### Note:
An updated version of this package is https://github.com/grd/stat.

### API Interfaces:
- It uses interfaces for the data access. So different data types can be used.
- Two datatypes are pre-defined: Float64 and Int64.
- For sampling purposes there is a Strider type provided.

Testing 100% pass. Testing covers the complete functionality.  
Tested on Debian6 32-bit and Windows 7 64-bit.

**Stable API**. I have absolutely no plans in changing the API.

### Offtopic:
- GSL 1.9 Statistics C library: 181 kb, 54 files
- Go Statistics library       :  83 kb, 26 files  :-)

