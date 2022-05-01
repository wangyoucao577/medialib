package nalu

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/ghodss/yaml"
)

// TypeInfo contains basic information of nalu type, such as description, etc.
type TypeInfo struct {
	Description string `json:"description"`
}

// NALU Type Codes, defined in ISO/IEC-14496-10 7.4.1
const (
	TypeUnspecified = iota
	TypeNonIDR
	TypeSliceDataPartitionA
	TypeSliceDataPartitionB
	TypeSliceDataPartitionC
	TypeIDR
	TypeSEI
	TypeSPS
	TypePPS
	TypeAccessUnitDelimiter
	TypeEndOfSequence
	TypeEndOfStream
	TypeFillerData
	TypeSPSExt
	TypePrefix
	TypeSubsetSPS
	TypeReserved16
	TypeReserved17
	TypeReserved18
	TypeAuxiliaryPicture
	TypeSliceExtersion
	TypeReserved21
	TypeReserved22
	TypeReserved23
	TypeTypeUnspecified24
	TypeTypeUnspecified25
	TypeTypeUnspecified26
	TypeTypeUnspecified27
	TypeTypeUnspecified28
	TypeTypeUnspecified29
	TypeTypeUnspecified30
	TypeTypeUnspecified31
)

var naluTypes = map[int]TypeInfo{
	TypeUnspecified:         {Description: "Unspecified"},
	TypeNonIDR:              {Description: "Coded slice of a non-IDR picture, slice_layer_without_partitioning_rbsp()"},
	TypeSliceDataPartitionA: {Description: "Coded slice data partition A, slice_data_partition_a_layer_rbsp()"},
	TypeSliceDataPartitionB: {Description: "Coded slice data partition B, slice_data_partition_a_layer_rbsp()"},
	TypeSliceDataPartitionC: {Description: "Coded slice data partition C, slice_data_partition_a_layer_rbsp()"},
	TypeIDR:                 {Description: "Coded slice of an IDR picture, slice_layer_without_partitioning_rbsp()"},
	TypeSEI:                 {Description: "Supplemental enhancement information (SEI), sei_rbsp( )"},
	TypeSPS:                 {Description: "Sequence parameter set, seq_parameter_set_rbsp()"},
	TypePPS:                 {Description: "Picture parameter set, pic_parameter_set_rbsp()"},
	TypeAccessUnitDelimiter: {Description: "Access unit delimiter, access_unit_delimiter_rbsp()"},
	TypeEndOfSequence:       {Description: "End of sequence, end_of_seq_rbsp()"},
	TypeEndOfStream:         {Description: "End of stream, end_of_stream_rbsp()"},
	TypeFillerData:          {Description: "Filler data, filler_data_rbsp( )"},
	TypeSPSExt:              {Description: "Sequence parameter set extension, seq_parameter_set_extension_rbsp()"},
	TypePrefix:              {Description: "Prefix NAL unit, prefix_nal_unit_rbsp(}"},
	TypeSubsetSPS:           {Description: "Subset sequence parameter set, subset_seq_parameter_set_rbsp()"},
	TypeReserved16:          {Description: "Reserved"},
	TypeReserved17:          {Description: "Reserved"},
	TypeReserved18:          {Description: "Reserved"},
	TypeAuxiliaryPicture:    {Description: "Coded slice of an auxiliary coded picture without partitioning,slice_layer_without_partitioning_rbsp()"},
	TypeSliceExtersion:      {Description: "Coded slice extension, slice_layer_extension_rbsp()"},
	TypeReserved21:          {Description: "Reserved"},
	TypeReserved22:          {Description: "Reserved"},
	TypeReserved23:          {Description: "Reserved"},
	TypeTypeUnspecified24:   {Description: "Unspecified"},
	TypeTypeUnspecified25:   {Description: "Unspecified"},
	TypeTypeUnspecified26:   {Description: "Unspecified"},
	TypeTypeUnspecified27:   {Description: "Unspecified"},
	TypeTypeUnspecified28:   {Description: "Unspecified"},
	TypeTypeUnspecified29:   {Description: "Unspecified"},
	TypeTypeUnspecified30:   {Description: "Unspecified"},
	TypeTypeUnspecified31:   {Description: "Unspecified"},
}

// TypeDescription represents nalu type description.
func TypeDescription(t int) string {
	n, ok := naluTypes[t]
	if !ok {
		return ""
	}
	return n.Description
}

// IsValidNALUType checks whether input NAL Unit Type is valid or not.
func IsValidNALUType(t int) bool {
	_, ok := naluTypes[t]
	return ok
}

// TypesMarshaler implements util.Marshaler
type TypesMarshaler struct{}

// JSON marshalls nalu types and relerrant information to JSON.
func (t TypesMarshaler) JSON() ([]byte, error) {
	return json.Marshal(naluTypes)
}

// JSONIndent marshalls nalu types to JSON representation with customized indent.
func (t TypesMarshaler) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(naluTypes, prefix, indent)
}

// YAML formats nalu types to YAML representation.
func (t TypesMarshaler) YAML() ([]byte, error) {
	j, err := json.Marshal(naluTypes)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CSV marshalls all supported nalu types to csv.
func (t TypesMarshaler) CSV() ([]byte, error) {
	records := [][]string{
		{"Type", "Description"}, // csv header
	}

	keys := make([]int, 0, len(naluTypes))
	for k := range naluTypes {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k := range keys {
		records = append(records, []string{strconv.Itoa(k), naluTypes[k].Description})
	}

	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	err := w.WriteAll(records)

	return buf.Bytes(), err
}
