package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// SubChunkRequest requests specific sub-chunks from the server using a center point.
type SubChunkRequest struct {
	// Dimension is the dimension of the sub-chunk.
	Dimension int32
	// Offsets contains all requested offsets around the center point.
	Offsets []protocol.SubChunkOffset
	// Position is an absolute sub-chunk center point used as a base point for all sub-chunks requested. The X and Z
	// coordinates represent the chunk coordinates, while the Y coordinate is the absolute sub-chunk index.
	Position protocol.SubChunkPos
}

// ID ...
func (*SubChunkRequest) ID() uint32 {
	return IDSubChunkRequest
}

// Marshal follows SubChunkRequestPacketPayload in Mojang's 1.26.40 schema, which orders the fields Dimension
// Type (0), SubChunk Position Offset List (1), Center Pos (2), with SubChunkPos being three plain int32s
// rather than varints. An earlier local change here had the position before the offsets and encoded it as
// varints, which is the pre-1.26.30 layout.
func (pk *SubChunkRequest) Marshal(io protocol.IO) {
	io.Varint32(&pk.Dimension)
	protocol.Slice(io, &pk.Offsets)
	io.SubChunkPos(&pk.Position)
}
