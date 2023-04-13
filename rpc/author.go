package rpc

// AuthorSubmitAndWatchExtrinsic  submits an extrinsic and watches it until block finalized
func AuthorSubmitAndWatchExtrinsic(id int, signedExtrinsic string) []byte {
	rpc := Param{Id: id, Method: "author_submitAndWatchExtrinsic", Params: []string{signedExtrinsic}}
	return rpc.structureQuery()
}

// AuthorSubmitExtrinsic submits an extrinsic and returns the hash
func AuthorSubmitExtrinsic(id int, signedExtrinsic string) []byte {
	rpc := Param{Id: id, Method: "author_submitExtrinsic", Params: []string{signedExtrinsic}}
	return rpc.structureQuery()
}
