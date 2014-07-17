### Overview
Here are some resources that provide a technical overview of ahimsa.
- http://ahimsa.io/about
- wiki

### Protocol Buffers
We use google protocol buffers to encode data in output scripts of bitcoin transactions.

This gives us two things:
- a language agnostic specification
- a specification

To contribute to this project or build tools for it. You need a protocol buffer extension
for the langauge you are operating in.

For golang the compiler was retreived from [here](http://code.google.com/p/goprotobuf/). 
The file bulletin.pb.go was built using this command:
```bash
protoc --go_out=./ bulletin.proto
```
