# servemuxpathtraversal
Script written in Go that checks targets for path traversal vulnerabilities induced by the ServeMux HTTP Multiplexer as described by Ilya Glotov's blog ->  https://ilyaglotov.com/blog/servemux-and-path-traversal

Usage:

 ```$ go run servemuxpathtraversal.go -t [target]```
 
 Or
 
 ```$ go run servemuxpathtraversal.go -i [targets_file]```
 
 
adding ports is optional and target should be in the form **http[s]://host[:port]**

example target = *https://localhost:50000*
