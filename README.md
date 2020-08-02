# ZKLib

This project is based on [zetcd](https://github.com/etcd-io/zetcd). 

New changes are you can use it as a library instead of client request. 

Now you are able to use ZK API with only one etcd instance w/o ZK server. 

## Why build it?
We try to pack Zetcd into a library instead of a server to make it a binary release. It will be better for operations management and the complexity of whole software. 


## How to use it?
See [API](./zjetcd_lib.go) and [DEMO](./zjetcd_lib_test.go) for more info.  
