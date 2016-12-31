package generator

const AVRO_BLOCK_SCHEMA = `
{"type": "record", "name": "AvroContainerBlock",
 "fields" : [
   {"name": "numRecords", "type": "long"},
   {"name": "recordBytes", "type": "bytes"},
   {"name": "sync", "type": {"type": "fixed", "name": "sync", "size": 16}}
  ]
}
`

const AVRO_HEADER_SCHEMA = `
{"type": "record", "name": "AvroContainerHeader",
 "fields" : [
   {"name": "magic", "type": {"type": "fixed", "name": "Magic", "size": 4}},
   {"name": "meta", "type": {"type": "map", "values": "bytes"}},
   {"name": "sync", "type": {"type": "fixed", "name": "Sync", "size": 16}}
  ]
}
`
