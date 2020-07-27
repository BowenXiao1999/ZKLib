package zetcd




func main()  {
	zk := NewZKClient([]string{"127.0.0.1:2380"})
	err := zk.Create("/hbase/meta-region-server", "9999")
	if err != nil {

	}

	resp, err := zk.Get("/hbase/meta-region-server")
	
	// expect to be 9999
	
	
	err := zk.Create("/hbase/meta-region-server", "8888")
	if err != nil {

	}

	resp, err := zk.Get("/hbase/meta-region-server")
	// expect to be 8888

}