package mp4

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/container/mp4/box/free"
	"github.com/wangyoucao577/medialib/container/mp4/box/ftyp"
	"github.com/wangyoucao577/medialib/container/mp4/box/mdat"
	"github.com/wangyoucao577/medialib/container/mp4/box/moof"
	"github.com/wangyoucao577/medialib/container/mp4/box/moov"
	"github.com/wangyoucao577/medialib/container/mp4/box/sidx"
	"github.com/wangyoucao577/medialib/container/mp4/box/wide"
	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/video/avc/annexbes"
	"github.com/wangyoucao577/medialib/video/avc/es"
)

// MoofMdat represents composition of one moof and one mdat, since they're stored interleavely like this.
type MoofMdat struct {
	Moof moof.Box `json:"moof"`
	Mdat mdat.Box `json:"mdat"`
}

// Boxes represents mp4 boxes.
type Boxes struct {
	Ftyp     *ftyp.Box  `json:"ftyp,omitempty"`
	Free     []free.Box `json:"free,omitempty"`
	Wide     *wide.Box  `json:"wide,omitempty"`
	Moov     *moov.Box  `json:"moov,omitempty"`
	MoofMdat []MoofMdat `json:"moof_mdat,omitempty"` // for fmp4, make sure moof,mdat can be pared and stored interleavely
	Mdat     []mdat.Box `json:"mdat,omitempty"`      // for  mp4 that doesn't have moof
	Sidx     []sidx.Box `json:"sidx,omitempty"`

	//TODO: other boxes

	// internal vars for parsing or other handling
	boxesCreator map[string]box.NewFunc `json:"-"`
}

func newBoxes() Boxes {
	return Boxes{
		boxesCreator: map[string]box.NewFunc{
			box.TypeFtyp: ftyp.New,
			box.TypeFree: free.New,
			box.TypeSkip: free.New,
			box.TypeWide: wide.New,
			box.TypeMdat: mdat.New,
			box.TypeMoov: moov.New,
			box.TypeMoof: moof.New,
			box.TypeSidx: sidx.New,
		},
	}
}

// JSON marshals boxes to JSON representation
func (b Boxes) JSON() ([]byte, error) {
	return json.Marshal(b)
}

// JSONIndent marshals boxes to JSON representation with customized indent.
func (b Boxes) JSONIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(b, prefix, indent)
}

// YAML formats boxes to YAML representation.
func (b Boxes) YAML() ([]byte, error) {
	j, err := json.Marshal(b)
	if err != nil {
		return j, err
	}
	return yaml.JSONToYAML(j)
}

// CSV formats boxes to CSV representation, which isn't supported at the moment.
func (b Boxes) CSV() ([]byte, error) {
	return nil, fmt.Errorf("csv representation does not support yet")
}

// CreateSubBox creates directly included box, such as create `mvhd` in `moov`, or create `moov` on top level.
//
//	return ErrNotImplemented is the box doesn't have any sub box.
func (b *Boxes) CreateSubBox(h box.Header) (box.Box, error) {
	creator, ok := b.boxesCreator[h.Type.String()]
	if !ok {
		glog.V(2).Infof("unknown box type %s, size %d payload %d", h.Type.String(), h.Size, h.PayloadSize())
		return nil, box.ErrUnknownBoxType
	}

	createdBox := creator(h)
	if createdBox == nil {
		glog.Fatalf("create box type %s failed", h.Type.String())
	}

	switch h.Type.String() {
	case box.TypeFtyp:
		b.Ftyp = createdBox.(*ftyp.Box)
	case box.TypeFree, box.TypeSkip:
		b.Free = append(b.Free, *createdBox.(*free.Box))
		createdBox = &b.Free[len(b.Free)-1] // reference to the last empty free box
	case box.TypeWide:
		b.Wide = createdBox.(*wide.Box)
	case box.TypeMdat:
		if len(b.MoofMdat) > 0 {
			if err := b.MoofMdat[len(b.MoofMdat)-1].Mdat.Validate(); err == nil { // expect error
				glog.Warningf("expect empty mdat but got a valid one %v", b.MoofMdat[len(b.MoofMdat)-1].Mdat)
				b.MoofMdat = append(b.MoofMdat, MoofMdat{}) // append new one to avoid lost mdat
			}

			b.MoofMdat[len(b.MoofMdat)-1].Mdat = *createdBox.(*mdat.Box)
			createdBox = &b.MoofMdat[len(b.MoofMdat)-1].Mdat
		} else if len(b.MoofMdat) == 0 {
			b.Mdat = append(b.Mdat, *createdBox.(*mdat.Box))
			createdBox = &b.Mdat[len(b.Mdat)-1]
		}
	case box.TypeMoov:
		b.Moov = createdBox.(*moov.Box)
	case box.TypeMoof:
		// Moof is required present before Mdat, so always create a new one if moof encountered.
		b.MoofMdat = append(b.MoofMdat, MoofMdat{Moof: *createdBox.(*moof.Box)})
		createdBox = &b.MoofMdat[len(b.MoofMdat)-1].Moof // reference to the last empty moof box
	case box.TypeSidx:
		b.Sidx = append(b.Sidx, *createdBox.(*sidx.Box))
		createdBox = &b.Sidx[len(b.Sidx)-1] // reference to the last empty sidx box
	}

	return createdBox, nil
}

// ParsePayload acts as an root box to parse all sub boxes.
func (b *Boxes) ParsePayload(r io.Reader) error {

	for {
		if _, err := box.ParseBox(r, b, math.MaxUint64); err != nil {
			if err == io.EOF || err == box.ErrInsufficientSize {
				break
			} else if err == box.ErrUnknownBoxType {
				continue
			}
			return err
		}
	}

	return nil
}

// ExtractES extracts AVC or HEVC Elementary Stream.
// Use trackID to select the specified one, trackID <= 0 means use the first found one.
func (b *Boxes) ExtractES(trackID int) (*es.ElementaryStream, error) {

	if b.Moov == nil || (b.MoofMdat == nil && b.Mdat == nil) {
		return nil, fmt.Errorf("moov, moof or mdat not found")
	}

	trackFound := false
	e := es.ElementaryStream{}
	for _, track := range b.Moov.Trak {
		if track.Mdia.Hdlr.HandlerType.String() == box.TypeVide {
			if trackID > 0 && uint32(trackID) != track.Tkhd.TrackID {
				continue
			}
			if trackID <= 0 {
				trackID = int(track.Tkhd.TrackID)
			}
			trackFound = true
			e.SetLengthSize(uint32(track.Mdia.Minf.Stbl.Stsd.AVC1SampleEntries[0].AVCConfig.AVCConfig.LengthSize()))
			break
		}
	}

	if !trackFound {
		return nil, fmt.Errorf("trackID %d not found", trackID)
	}

	// fragment-mp4 if exist
	for i := 0; i < len(b.MoofMdat); i++ {
		for _, tf := range b.MoofMdat[i].Moof.Traf {
			if int(tf.Tfhd.TrackID) != trackID {
				continue
			}

			for _, tr := range tf.Trun {
				var startPos uint32
				for _, sampleSize := range tr.SampleSize {
					data := b.MoofMdat[i].Mdat.Data[startPos : startPos+sampleSize]
					if _, err := e.Parse(bytes.NewReader(data), len(data)); err != nil {
						return &e, err
					}
					startPos += sampleSize
				}
			}
			break
		}
	}

	// mp4 if exist
	for i := 0; i < len(b.Mdat); i++ {

	}

	return &e, nil
}

// ExtractAnnexBES extracts AVC or HEVC Elementary Stream with AnnexB byte format.
// Use trackID to select the specified one, trackID <= 0 means use the first found one.
func (b *Boxes) ExtractAnnexBES(trackID int) (*annexbes.ElementaryStream, error) {
	mp4ES, err := b.ExtractES(trackID)
	if err != nil {
		return nil, err
	}

	annexbES := annexbes.ElementaryStream{}
	for i := range mp4ES.LengthNALU {
		annexbES.NALU = append(annexbES.NALU, mp4ES.LengthNALU[i].NALU)
	}

	return &annexbES, nil
}

// DumpDurations dumps duration information.
func (b *Boxes) DumpDurations() {

	type trackInfo struct {
		// from moov
		trackID   uint32
		trackType string
		timescale uint64

		// from fragments
		duration    uint64
		sampleCount uint64
	}
	tracksInfo := map[uint32]*trackInfo{}

	if b.Moov != nil && b.Moov.Mvhd != nil { // moov
		mvhdTimescale := b.Moov.Mvhd.Timescale
		mvhdDuration := b.Moov.Mvhd.Duration
		glog.Infof("movie (mvhd) timescale %d duration %d(%fs)", mvhdTimescale, mvhdDuration, util.DurationInSeconds(mvhdDuration, uint64(mvhdTimescale)))

		for _, track := range b.Moov.Trak {
			if track.Tkhd == nil {
				continue
			}
			trackID := track.Tkhd.TrackID

			trackType := "unknown"
			if track.Mdia != nil && track.Mdia.Hdlr != nil {
				switch track.Mdia.Hdlr.HandlerType.String() {
				case "vide":
					trackType = "video"
				case "soun":
					trackType = "audio"
				}
			}

			tkhdDuration := track.Tkhd.Duration
			glog.Infof("  track %d %s (tkhd) duration %d (%fs)", trackID, trackType, tkhdDuration, util.DurationInSeconds(tkhdDuration, uint64(mvhdTimescale)))

			var mdhdTimescale uint32
			if track.Mdia != nil && track.Mdia.Mdhd != nil {
				mdhdTimescale = track.Mdia.Mdhd.Timescale
				mdhdDuration := track.Mdia.Mdhd.Duration
				glog.Infof("  track %d %s (mdhd) timescale %d duration %d (%fs)", trackID, trackType, mdhdTimescale, mdhdDuration, util.DurationInSeconds(mdhdDuration, uint64(mdhdTimescale)))
			}

			tracksInfo[trackID] = &trackInfo{trackID: trackID, trackType: trackType, timescale: uint64(mdhdTimescale)}
		}
	}

	if len(b.MoofMdat) == 0 {
		return
	}

	for _, moofMdat := range b.MoofMdat { // fragments
		for _, moofTraf := range moofMdat.Moof.Traf {
			if moofTraf.Tfhd == nil {
				continue
			}
			trackID := moofTraf.Tfhd.TrackID

			if _, ok := tracksInfo[trackID]; !ok { // set default track info if no header
				tracksInfo[trackID] = &trackInfo{
					trackID, "unknown", 0, 0, 0,
				}
			}
			currInfo := tracksInfo[trackID]

			var defaultSampleDuration uint32
			if moofTraf.Tfhd.FullHeader.Flags&0x000008 > 0 { // default sample duration found
				defaultSampleDuration = moofTraf.Tfhd.DefaultSampleDuration
			}

			for _, trun := range moofTraf.Trun {
				currInfo.sampleCount += uint64(trun.SampleCount)

				if trun.Flags&0x100 > 0 { // sample‐duration‐present, indicates that each sample has its own duration
					for _, d := range trun.SampleDuration {
						currInfo.duration += uint64(d)
					}
					glog.V(2).Infof("    track %d %s use sample durations for sample count %d", trackID, tracksInfo[trackID].trackType, trun.SampleCount)
				} else { // otherwise the default is used (default‐sample‐duration‐present)
					glog.V(2).Infof("    track %d %s use default sample duration %d for sample count %d", trackID, tracksInfo[trackID].trackType, defaultSampleDuration, trun.SampleCount)
					currInfo.duration += uint64(defaultSampleDuration * trun.SampleCount)
				}
			}
		}
	}

	for _, v := range tracksInfo {
		glog.Infof("  track %d %s (tfhd/trun) sample_count %d duration %d (%fs)", v.trackID, v.trackType, v.sampleCount, v.duration, util.DurationInSeconds(v.duration, v.timescale))
	}
}
