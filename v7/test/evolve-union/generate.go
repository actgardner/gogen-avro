package avro

//go:generate $GOPATH/bin/gogen-avro . defaults.avsc
//go:generate mkdir -p evolution
//go:generate $GOPATH/bin/gogen-avro evolution evolution.avsc
