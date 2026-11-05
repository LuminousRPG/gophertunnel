package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ClientCacheBlobStatus is part of the blob cache protocol. It is sent by the client to let the server know
// what blobs it needs and which blobs it already has, in an ACK type system.
type ClientCacheBlobStatus struct {
	// MissHashes is a list of blob hashes that the client does not have a blob available for. The server
	// should send the blobs matching these hashes as soon as possible.
	MissHashes []uint64
	// HitHashes is a list of blob hashes that the client has a blob available for. The blobs hashes here mean
	// that the client already has them: The server does not need to send the blobs anymore.
	HitHashes []uint64
}

// ID ...
func (pk *ClientCacheBlobStatus) ID() uint32 {
	return IDClientCacheBlobStatus
}

// maxBlobStatusHashes bounds each hash list. The client answers every sub chunk blob it was sent at
// once, so a large view loading after a join or teleport passes the usual slice limit of 1024.
const maxBlobStatusHashes = 1 << 16

func (pk *ClientCacheBlobStatus) Marshal(io protocol.IO) {
	blobHashes(io, &pk.MissHashes)
	blobHashes(io, &pk.HitHashes)
}

func blobHashes(io protocol.IO, x *[]uint64) {
	count := uint32(len(*x))
	io.Varuint32(&count)
	if r, ok := io.(interface{ SliceLength(value, max uint32) }); ok {
		r.SliceLength(count, maxBlobStatusHashes)
		*x = make([]uint64, count)
	}
	for i := range *x {
		io.Uint64(&(*x)[i])
	}
}
