package box

import (
	"bytes"
	"encoding/csv"
	"encoding/json"

	"github.com/ghodss/yaml"
)

// BasicInfo contains basic information of box, such as name, etc.
type BasicInfo struct {
	Name string `json:"name"`
}

// Box types
const (
	TypeUUID = "uuid"

	TypeFtyp   = "ftyp"
	TypeFree   = "free"
	TypeSkip   = "skip"
	TypeMdat   = "mdat"
	TypeMoov   = "moov"
	TypeMvhd   = "mvhd"
	TypeUdta   = "udta"
	TypeCprt   = "cprt"
	TypeMeta   = "meta"
	TypeHdlr   = "hdlr"
	TypeIlst   = "ilst"
	TypeTrak   = "trak"
	TypeTkhd   = "tkhd"
	TypeMdia   = "mdia"
	TypeMdhd   = "mdhd"
	TypeMinf   = "minf"
	TypeStbl   = "stbl"
	TypeDinf   = "dinf"
	TypeSmhd   = "smhd"
	TypeVmhd   = "vmhd"
	TypeStsd   = "stsd"
	TypeStts   = "stts"
	TypeStss   = "stss"
	TypeStsc   = "stsc"
	TypeStsz   = "stsz"
	TypeStco   = "stco"
	TypeCtts   = "ctts"
	TypeSdtp   = "sdtp"
	TypeDref   = "dref"
	TypeUrl    = "url "
	TypeUrn    = "urn"
	TypeMoof   = "moof"
	TypeMfhd   = "mfhd"
	TypeTraf   = "traf"
	TypeTfhd   = "tfhd"
	TypeTrun   = "trun"
	TypeTfdt   = "tfdt"
	TypeMvex   = "mvex"
	TypeMehd   = "mehd"
	TypeTrex   = "trex"
	TypeEdts   = "edts"
	TypeElst   = "elst"
	TypeSidx   = "sidx"
	TypeData   = "data"
	TypeDottoo = "\251too"
	TypeDesc   = "desc"

	// sample entry
	TypeVide = "vide"
	TypeSoun = "soun"
	TypeAvc1 = "avc1"
	TypeAvcC = "avcC"
	TypeHev1 = "hev1"
	TypeHvc1 = "hvc1"
	TypeHvcC = "hvcC"
	TypeLhvC = "lhvC"
	TypeAv01 = "av01"
	TypeAv1C = "av1C"
	TypeBtrt = "btrt"
	TypeMp4a = "mp4a"
	TypeEsds = "esds"
	TypePasp = "pasp"
	TypeVexu = "vexu"
	TypeColr = "colr"
	TypeHfov = "hfov"
	TypeSgpd = "sgpd"
)

var boxTypes = map[string]BasicInfo{
	TypeUUID: {Name: "UUID"},

	TypeFtyp:   {Name: "File Type Box"},
	TypeFree:   {Name: "Free Space Box"},
	TypeSkip:   {Name: "Free Space Box"},
	TypeMdat:   {Name: "Media Data Box"},
	TypeMoov:   {Name: "Movie Box"},
	TypeMvhd:   {Name: "Movie Header Box"},
	TypeUdta:   {Name: "User Data Box"},
	TypeCprt:   {Name: "Copyright Box"},
	TypeMeta:   {Name: "Meta Box"},
	TypeHdlr:   {Name: "Handler Reference Box"},
	TypeIlst:   {Name: "unknown"},
	TypeTrak:   {Name: "Track Reference Box"},
	TypeTkhd:   {Name: "Track Header Box"},
	TypeMdia:   {Name: "Media Box"},
	TypeMdhd:   {Name: "Media Header Box"},
	TypeMinf:   {Name: "Media Information Box"},
	TypeStbl:   {Name: "Sample Table Box"},
	TypeDinf:   {Name: "Data Information Box"},
	TypeSmhd:   {Name: "Sound Media Header"},
	TypeVmhd:   {Name: "Video Media Header"},
	TypeStsd:   {Name: "Sample Description Box"},
	TypeStts:   {Name: "Decoding Time to Sample Box"},
	TypeStss:   {Name: "Sync Sample Box"},
	TypeStsc:   {Name: "Sample To Chunk Box"},
	TypeStsz:   {Name: "Sample Size Box"},
	TypeStco:   {Name: "Chunk Offset Box"},
	TypeCtts:   {Name: "Composition Time to Sample Box"},
	TypeSdtp:   {Name: "Independent and Disposable Samples Box"},
	TypeDref:   {Name: "Data Reference Box"},
	TypeUrl:    {Name: "Data Entry Url Box"},
	TypeUrn:    {Name: "Data Entry Urn Box"},
	TypeMoof:   {Name: "Movie Fragment Box"},
	TypeMfhd:   {Name: "Movie Fragment Header Box"},
	TypeTraf:   {Name: "Track Fragment Box"},
	TypeTfhd:   {Name: "Track Fragment Header Box"},
	TypeTrun:   {Name: "Track Fragment Run Box"},
	TypeTfdt:   {Name: "Track Fragment Base Media Decode Time Box"},
	TypeMvex:   {Name: "Movie Extends Box"},
	TypeMehd:   {Name: "Movie Extends Header Box"},
	TypeTrex:   {Name: "Track Extends Box"},
	TypeEdts:   {Name: "Edit Box"},
	TypeElst:   {Name: "Edit List Box"},
	TypeSidx:   {Name: "Segment Index Box"},
	TypeData:   {Name: "Data Box (Quicktime file format defined)"},
	TypeDottoo: {Name: "Encoding tool information box"},
	TypeDesc:   {Name: "Description Box"},

	TypeVide: {Name: "Visual Sample Entry"},
	TypeSoun: {Name: "Audio Sample Entry"},
	TypeAvc1: {Name: "AVC Sample Entry"},
	TypeAvcC: {Name: "AVC Configuration Box"},
	TypeHev1: {Name: "HEVC Sample Entry"},
	TypeHvc1: {Name: "HEVC Sample Entry"},
	TypeHvcC: {Name: "HEVC Decoder Configuration Record"},
	TypeLhvC: {Name: "L-HEVC Decoder Configuration Record"},
	TypeAv01: {Name: "AV1 Sample Entry"},
	TypeAv1C: {Name: "AV1 Configuration Box"},
	TypeBtrt: {Name: "MPEG4 Bit Rate Box"},
	TypeMp4a: {Name: "MP4 Visual Sample Entry"},
	TypeEsds: {Name: "ES Descriptor Box"},
	TypePasp: {Name: "Pixel Aspect Ratio Box"},
	TypeVexu: {Name: "Video Extended Usage Box"},
	TypeColr: {Name: "Colour Information Box"},
	TypeHfov: {Name: "hfov Box"},
	TypeSgpd: {Name: "Sample Group Description Box"},
}

// BoxTypes returns box types map.
func BoxTypes() map[string]BasicInfo {
	return boxTypes
}

// IsValidBoxType checks whether input box type is valid or not.
func IsValidBoxType(t string) bool {
	_, ok := boxTypes[t]
	return ok
}

// TypesMarshaler implements util.Marshaler.
type TypesMarshaler struct{}

// JSON marshalls supported boxes and relerrant information to JSON.
func (t TypesMarshaler) JSON() ([]byte, error) {
	return json.Marshal(boxTypes)
}

// JSONIndent marshalls boxes to JSON representation with customized indent.
func (t TypesMarshaler) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(boxTypes, prefix, indent)
}

// YAML formats boxes to YAML representation.
func (t TypesMarshaler) YAML() ([]byte, error) {
	j, err := json.Marshal(boxTypes)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CSV marshalls all supported boxes to csv.
func (t TypesMarshaler) CSV() ([]byte, error) {
	records := [][]string{
		{"Type", "Name"}, // csv header
	}

	for k, v := range boxTypes {
		records = append(records, []string{k, v.Name})
	}

	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	err := w.WriteAll(records)

	return buf.Bytes(), err
}
