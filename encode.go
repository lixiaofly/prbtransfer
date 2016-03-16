package prbtransfer

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
)

type ProbDataForm interface {
	ProbFormatComb() []byte
	ProbDataType() uint32
	DesKey() []byte
	DesIv() []byte
}

type probProcPack struct {
	DataTypeBig []byte
	BodySizeBig []byte
	Context     []byte
}

func ProbDataDeal(pdf ProbDataForm) ([]byte, error) {
	datatype := make([]byte, 4)
	datasize := make([]byte, 4)

	binary.BigEndian.PutUint32(datatype, pdf.ProbDataType())
	binary.BigEndian.PutUint32(datasize, uint32(len(pdf.ProbFormatComb())))
	fmt.Println("len(buf3)=", uint32(len(pdf.ProbFormatComb())), "datasize=", datasize)
	procPack := &probProcPack{
		DataTypeBig: datatype,
		BodySizeBig: datasize,
		Context:     make([]byte, 80),
	}
	fmt.Println("context before zip:", pdf.ProbFormatComb())
	procPack.Context = probeDateZip(pdf.ProbFormatComb())
	fmt.Println("context after zip:", procPack.Context)
	return probeDesEncry(structToBuff(procPack), pdf.DesKey(), pdf.DesIv())
}

func structToBuff(strt *probProcPack) []byte {
	buf := strt.DataTypeBig	
	fmt.Println("buf1=", strt.DataTypeBig)
	fmt.Println("buf2=", strt.BodySizeBig)
	fmt.Println("buf3=", strt.Context, "len(buf3)=", len(strt.Context))
	buf = append(buf, strt.BodySizeBig...)
	buf = append(buf, strt.Context...)
	fmt.Printf("\n")
	fmt.Println("context before des:", buf, "len=",len(buf))
	return buf
}
func probeDateZip(context []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	//w.Write([]byte("hello, world\n"))
	w.Write(context)
	w.Close()
	return b.Bytes()
	//fmt.Println(b.Bytes())
}
func probeDesEncry(buf, key, iv []byte) ([]byte, error) {
	datalen := make([]byte, 4)
	binary.BigEndian.PutUint32(datalen, uint32(len(buf)))
	buf = append(buf, datalen...)
	
	pading := []byte{'A', 'B', 'C', 'D'}
	buf = append(buf, pading...)

	return DesEncrypt(buf, key, iv)
}
