package astits

import (
	"testing"

	"github.com/asticode/go-astitools/binary"
	"github.com/stretchr/testify/assert"
)

func TestParseData(t *testing.T) {
	// Init
	pm := newProgramMap()
	ps := []*Packet{}

	// Custom parser
	cds := []*Data{{PID: 1}}
	var c = func(ps []*Packet) (o []*Data, skip bool, err error) {
		o = cds
		skip = true
		return
	}
	ds, err := parseData(ps, c, pm)
	assert.NoError(t, err)
	assert.Equal(t, cds, ds)

	// Do nothing for CAT
	ps = []*Packet{{Header: &PacketHeader{PID: PIDCAT}}}
	ds, err = parseData(ps, nil, pm)
	assert.NoError(t, err)
	assert.Empty(t, ds)

	// PES
	p := pesWithHeaderBytes()
	ps = []*Packet{
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[:33],
		},
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[33:],
		},
	}
	ds, err = parseData(ps, nil, pm)
	assert.NoError(t, err)
	assert.Equal(t, []*Data{{PES: pesWithHeader, PID: uint16(256)}}, ds)

	// PSI
	pm.set(uint16(256), uint16(1))
	p = psiBytes()
	ps = []*Packet{
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[:33],
		},
		{
			Header:  &PacketHeader{PID: uint16(256)},
			Payload: p[33:],
		},
	}
	ds, err = parseData(ps, nil, pm)
	assert.NoError(t, err)
	assert.Equal(t, psi.toData(uint16(256)), ds)
}

func TestIsPSIPayload(t *testing.T) {
	pm := newProgramMap()
	var pids []int
	for i := 0; i <= 255; i++ {
		if isPSIPayload(uint16(i), pm) {
			pids = append(pids, i)
		}
	}
	assert.Equal(t, []int{0, 16, 17, 18, 19, 20, 30, 31}, pids)
	pm.set(uint16(1), uint16(0))
	assert.True(t, isPSIPayload(uint16(1), pm))
}

func TestIsPESPayload(t *testing.T) {
	w := astibinary.New()
	w.Write("0000000000000001")
	assert.False(t, isPESPayload(w.Bytes()))
	w.Reset()
	w.Write("000000000000000000000001")
	assert.True(t, isPESPayload(w.Bytes()))
}
