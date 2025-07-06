package types

func NewMsgStoreCode(creator string, wasmByteCode []byte) *MsgStoreCode {
	return &MsgStoreCode{
		Creator:      creator,
		WasmByteCode: wasmByteCode,
	}
}
