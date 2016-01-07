package main

var conf = struct {
	realPanic      bool
	goFilesPause   bool
	jadeFilesPause bool
	sassFilesPause bool
	sqlFilesPause  bool
	jsFilesPause   bool
}{
	false,
	true,
	true,
	false,
	false,
	false,
}
