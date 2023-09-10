[![CI](https://github.com/wangyoucao577/medialib/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyoucao577/medialib/actions/workflows/ci.yml)    
# medialib

## Tools
Below tools are available in [cmd](cmd).     

```
cmd
|-- flv2avc
|-- flvdump
|-- h26xdump
|-- mp42avc
`-- mp4dump
```


| Name | Description | 
| - | - |
| `flvdump` | displays the tags structure of an flv file, as `json` or `yaml` |
| `flv2avc` | extract a raw AVC/H.264 elementary stream from an flv file |
| `mp4dump` | displays the entire atom/box structure of an mp4 file, as `json` or `yaml` |
| `mp42avc` | extract a raw AVC/H.264 elementary stream from an mp4 file, only support fragemented mp4 at the moment |
| `h26xdump`| displays the structure of an AVC/H.264 elementary stream file, as `json` or `yaml`. Be aware that the Elementary Stream file is mandatory stored by AnnexB byte stream format. |

### Examples     

- dump `tags` of an `flv` file    

```
./flvdump -logtostderr -i in.flv -o dump.json
./flvdump -logtostderr -i in.flv -format yaml -o dump.yaml 
```

- extract `.h264` of an `flv` file 

```
./flv2avc -logtostderr -i in.flv -o out.h264 
```

- dump `boxes` of an `mp4` file    

```
./mp4dump -logtostderr -i in.mp4 -o dump.json
./mp4dump -logtostderr -i in.mp4 -format yaml -o dump.yaml 
```

- extract `.h264` of an `mp4` file 

```
./mp42avc -logtostderr -i in.mp4 -o out.h264 
```

- dump `nalus` of an `.h264` file

```
./h26xdump -logtostderr -i in.h264 -o dump.json
./h26xdump -logtostderr -i in.h264 -format yaml -o dump.yaml 
```

