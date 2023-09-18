[![CI](https://github.com/wangyoucao577/medialib/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyoucao577/medialib/actions/workflows/ci.yml)    
# medialib

## Tools
Below tools are available in [cmd](cmd).     

```
cmd
├── flv2avc
├── mediadump
└── mp42avc
```


| Name | Description | 
| - | - |
| `mediadump` | displays the container or elementary stream structure of an input media file, as `json` or `yaml` |
| `flv2avc` | extract a raw AVC/H.264 elementary stream from an flv file |
| `mp42avc` | extract a raw AVC/H.264 elementary stream from an mp4 file, only support fragemented mp4 at the moment |

### Examples     

- dump `tags` of an `flv` file    

```
./mediadump -logtostderr -i in.flv -o dump.json
./mediadump -logtostderr -i in.flv -of yaml -o dump.yaml 
```

- dump `boxes` of an `mp4` file    

```
./mediadump -logtostderr -i in.mp4 -o dump.json
./mediadump -logtostderr -i in.mp4 -of yaml -o dump.yaml 
```

- dump `nalus` of an `.h264` file

```
./mediadump -logtostderr -i in.h264 -o dump.json
./mediadump -logtostderr -i in.h264 -of yaml -o dump.yaml 
```

- extract `.h264` of an `flv` file 

```
./flv2avc -logtostderr -i in.flv -o out.h264 
```

- extract `.h264` of an `mp4` file 

```
./mp42avc -logtostderr -i in.mp4 -o out.h264 
```

