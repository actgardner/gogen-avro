package avro

//go:generate $GOPATH/bin/gogen-avro -namespaced-names=full -containers . schema.avsc
//go:generate mkdir -p evolution
//go:generate $GOPATH/bin/gogen-avro -namespaced-names=full -containers evolution evolution.avsc
