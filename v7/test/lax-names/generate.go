package avro

//go:generate $GOPATH/bin/gogen-avro -containers . schema.avsc
//go:generate mkdir -p evolution
//go:generate $GOPATH/bin/gogen-avro -containers evolution evolution.avsc
