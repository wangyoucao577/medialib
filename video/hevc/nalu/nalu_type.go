// Package nalu represents HEVC NAL Units.
package nalu

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/ghodss/yaml"
)

// TypeInfo contains basic information of nalu type, such as name, description, etc.
type TypeInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NALU Type Codes, defined in Rec. ITU-T H.265 Table 7-1 – NAL unit type codes and NAL unit type classes.
const (
	TypeTRAIL_N = iota
	TypeTRAIL_R
	TypeTSA_N
	TypeTSA_R
	TypeSTSA_N
	TypeSTSA_R
	TypeRADL_N
	TypeRADL_R
	TypeRASL_N
	TypeRASL_R
	TypeRSV_VCL_N10
	TypeRSV_VCL_R11
	TypeRSV_VCL_N12
	TypeRSV_VCL_R13
	TypeRSV_VCL_N14
	TypeRSV_VCL_R15
	TypeBLA_W_LP
	TypeBLA_W_RADL
	TypeBLA_N_LP
	TypeIDR_W_RADL
	TypeIDR_N_LP
	TypeCRA_NUT
	TypeRSV_IRAP_VCL22
	TypeRSV_IRAP_VCL23
	TypeRSV_VCL24
	TypeRSV_VCL25
	TypeRSV_VCL26
	TypeRSV_VCL27
	TypeRSV_VCL28
	TypeRSV_VCL29
	TypeRSV_VCL30
	TypeRSV_VCL31
	TypeVPS_NUT
	TypeSPS_NUT
	TypePPS_NUT
	TypeAUD_NUT
	TypeEOS_NUT
	TypeEOB_NUT
	TypeFD_NUT
	TypePREFIX_SEI_NUT
	TypeSUFFIX_SEI_NUT
	TypeRSV_NVCL41
	TypeRSV_NVCL42
	TypeRSV_NVCL43
	TypeRSV_NVCL44
	TypeRSV_NVCL45
	TypeRSV_NVCL46
	TypeRSV_NVCL47
	TypeUNSPEC48
	TypeUNSPEC49
	TypeUNSPEC50
	TypeUNSPEC51
	TypeUNSPEC52
	TypeUNSPEC53
	TypeUNSPEC54
	TypeUNSPEC55
	TypeUNSPEC56
	TypeUNSPEC57
	TypeUNSPEC58
	TypeUNSPEC59
	TypeUNSPEC60
	TypeUNSPEC61
	TypeUNSPEC62
	TypeUNSPEC63
)

var naluTypes = map[int]TypeInfo{
	TypeTRAIL_N:        {Name: "TRAIL_N", Description: "Coded slice segment of a non-TSA, non-STSA trailing picture, slice_segment_layer_rbsp()"},
	TypeTRAIL_R:        {Name: "TRAIL_R", Description: "Coded slice segment of a non-TSA, non-STSA trailing picture, slice_segment_layer_rbsp()"},
	TypeTSA_N:          {Name: "TSA_N", Description: "Coded slice segment of a TSA picture, slice_segment_layer_rbsp()"},
	TypeTSA_R:          {Name: "TSA_R", Description: "Coded slice segment of a TSA picture, slice_segment_layer_rbsp()"},
	TypeSTSA_N:         {Name: "STSA_N", Description: "Coded slice segment of an STSA picture, slice_segment_layer_rbsp()"},
	TypeSTSA_R:         {Name: "STSA_R", Description: "Coded slice segment of an STSA picture, slice_segment_layer_rbsp()"},
	TypeRADL_N:         {Name: "RADL_N", Description: "Coded slice segment of a RADL picture, slice_segment_layer_rbsp()"},
	TypeRADL_R:         {Name: "RADL_R", Description: "Coded slice segment of a RADL picture, slice_segment_layer_rbsp()"},
	TypeRASL_N:         {Name: "RASL_N", Description: "Coded slice segment of a RASL picture, slice_segment_layer_rbsp()"},
	TypeRASL_R:         {Name: "RASL_R", Description: "Coded slice segment of a RASL picture, slice_segment_layer_rbsp()"},
	TypeRSV_VCL_N10:    {Name: "RSV_VCL_N10", Description: "Reserved non-IRAP SLNR VCL NAL unit types"},
	TypeRSV_VCL_R11:    {Name: "RSV_VCL_R11", Description: "Reserved non-IRAP sub-layer reference VCL NAL unit types"},
	TypeRSV_VCL_N12:    {Name: "RSV_VCL_N12", Description: "Reserved non-IRAP SLNR VCL NAL unit types"},
	TypeRSV_VCL_R13:    {Name: "RSV_VCL_R13", Description: "Reserved non-IRAP sub-layer reference VCL NAL unit types"},
	TypeRSV_VCL_N14:    {Name: "RSV_VCL_N14", Description: "Reserved non-IRAP SLNR VCL NAL unit types"},
	TypeRSV_VCL_R15:    {Name: "RSV_VCL_R15", Description: "Reserved non-IRAP sub-layer reference VCL NAL unit types"},
	TypeBLA_W_LP:       {Name: "BLA_W_LP", Description: "Coded slice segment of a BLA picture, slice_segment_layer_rbsp()"},
	TypeBLA_W_RADL:     {Name: "BLA_W_RADL", Description: "Coded slice segment of a BLA picture, slice_segment_layer_rbsp()"},
	TypeBLA_N_LP:       {Name: "BLA_N_LP", Description: "Coded slice segment of a BLA picture, slice_segment_layer_rbsp()"},
	TypeIDR_W_RADL:     {Name: "IDR_W_RADL", Description: "Coded slice segment of an IDR picture, slice_segment_layer_rbsp()"},
	TypeIDR_N_LP:       {Name: "IDR_N_LP", Description: "Coded slice segment of an IDR picture, slice_segment_layer_rbsp()"},
	TypeCRA_NUT:        {Name: "CRA_NUT", Description: "Coded slice segment of a CRA picture, slice_segment_layer_rbsp()"},
	TypeRSV_IRAP_VCL22: {Name: "RSV_IRAP_VCL22", Description: "Reserved IRAP VCL NAL unit types"},
	TypeRSV_IRAP_VCL23: {Name: "RSV_IRAP_VCL23", Description: "Reserved IRAP VCL NAL unit types"},
	TypeRSV_VCL24:      {Name: "RSV_VCL24", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL25:      {Name: "RSV_VCL25", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL26:      {Name: "RSV_VCL26", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL27:      {Name: "RSV_VCL27", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL28:      {Name: "RSV_VCL28", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL29:      {Name: "RSV_VCL29", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL30:      {Name: "RSV_VCL30", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeRSV_VCL31:      {Name: "RSV_VCL31", Description: "Reserved non-IRAP VCL NAL unit types"},
	TypeVPS_NUT:        {Name: "VPS_NUT", Description: "Video parameter set, video_parameter_set_rbsp()"},
	TypeSPS_NUT:        {Name: "SPS_NUT", Description: "Sequence parameter set, seq_parameter_set_rbsp()"},
	TypePPS_NUT:        {Name: "PPS_NUT", Description: "Picture parameter set, pic_parameter_set_rbsp()"},
	TypeAUD_NUT:        {Name: "AUD_NUT", Description: "Access unit delimiter, access_unit_delimiter_rbsp()"},
	TypeEOS_NUT:        {Name: "EOS_NUT", Description: "End of sequence, end_of_seq_rbsp()"},
	TypeEOB_NUT:        {Name: "EOB_NUT", Description: "End of bitstream, end_of_bitstream_rbsp()"},
	TypeFD_NUT:         {Name: "FD_NUT", Description: "Filler data, filler_data_rbsp()"},
	TypePREFIX_SEI_NUT: {Name: "PREFIX_SEI_NUT", Description: "Supplemental enhancement information, sei_rbsp()"},
	TypeSUFFIX_SEI_NUT: {Name: "SUFFIX_SEI_NUT", Description: "Supplemental enhancement information, sei_rbsp()"},
	TypeRSV_NVCL41:     {Name: "RSV_NVCL41", Description: "Reserved"},
	TypeRSV_NVCL42:     {Name: "RSV_NVCL42", Description: "Reserved"},
	TypeRSV_NVCL43:     {Name: "RSV_NVCL43", Description: "Reserved"},
	TypeRSV_NVCL44:     {Name: "RSV_NVCL44", Description: "Reserved"},
	TypeRSV_NVCL45:     {Name: "RSV_NVCL45", Description: "Reserved"},
	TypeRSV_NVCL46:     {Name: "RSV_NVCL46", Description: "Reserved"},
	TypeRSV_NVCL47:     {Name: "RSV_NVCL47", Description: "Reserved"},
	TypeUNSPEC48:       {Name: "UNSPEC48", Description: "Unspecified"},
	TypeUNSPEC49:       {Name: "UNSPEC49", Description: "Unspecified"},
	TypeUNSPEC50:       {Name: "UNSPEC50", Description: "Unspecified"},
	TypeUNSPEC51:       {Name: "UNSPEC51", Description: "Unspecified"},
	TypeUNSPEC52:       {Name: "UNSPEC52", Description: "Unspecified"},
	TypeUNSPEC53:       {Name: "UNSPEC53", Description: "Unspecified"},
	TypeUNSPEC54:       {Name: "UNSPEC54", Description: "Unspecified"},
	TypeUNSPEC55:       {Name: "UNSPEC55", Description: "Unspecified"},
	TypeUNSPEC56:       {Name: "UNSPEC56", Description: "Unspecified"},
	TypeUNSPEC57:       {Name: "UNSPEC57", Description: "Unspecified"},
	TypeUNSPEC58:       {Name: "UNSPEC58", Description: "Unspecified"},
	TypeUNSPEC59:       {Name: "UNSPEC59", Description: "Unspecified"},
	TypeUNSPEC60:       {Name: "UNSPEC60", Description: "Unspecified"},
	TypeUNSPEC61:       {Name: "UNSPEC61", Description: "Unspecified"},
	TypeUNSPEC62:       {Name: "UNSPEC62", Description: "Unspecified"},
	TypeUNSPEC63:       {Name: "UNSPEC63", Description: "Unspecified"},
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
		{"Type", "Name", "Description"}, // csv header
	}

	keys := make([]int, 0, len(naluTypes))
	for k := range naluTypes {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k := range keys {
		records = append(records, []string{strconv.Itoa(k), naluTypes[k].Name, naluTypes[k].Description})
	}

	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	err := w.WriteAll(records)

	return buf.Bytes(), err
}
