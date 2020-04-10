package avro

//go:generate $GOPATH/bin/gogen-avro . union.avsc
//go:generate mkdir -p evolution
//go:generate $GOPATH/bin/gogen-avro evolution evolution.avsc
