package prbtransfer

import (
	//"bytes"
	//"compress/zlib"
	"archive/zip"
	"encoding/binary"
	"fmt"
	//"io"
	"io/ioutil"
	"os"
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

	fmt.Println("ProbFormatComb=", string(pdf.ProbFormatComb()))
	binary.BigEndian.PutUint32(datatype, pdf.ProbDataType())
	binary.BigEndian.PutUint32(datasize, uint32(len(pdf.ProbFormatComb())))
	fmt.Println("len(buf3)=", uint32(len(pdf.ProbFormatComb())), "datasize=", datasize)
	procPack := &probProcPack{
		DataTypeBig: datatype,
		BodySizeBig: datasize,
		Context:     make([]byte, 80),
	}
	fmt.Println("context before zip:", pdf.ProbFormatComb(), "len=", len(pdf.ProbFormatComb()))
	var err error
	procPack.Context, err = probeDesEncry(probeDateZip(pdf.ProbFormatComb()), pdf.DesKey(), pdf.DesIv())
	if err != nil {
		fmt.Println("Error: des encrypt failed!")
		return []byte(""), err
	}
	fmt.Println("context after des:", procPack.Context, "len=", len(procPack.Context))
	//return probeDesEncry(structToBuff(procPack), pdf.DesKey(), pdf.DesIv())
	return structToBuff(procPack), nil
}

func structToBuff(strt *probProcPack) []byte {
	buf := strt.DataTypeBig
	//fmt.Println("buf1=", strt.DataTypeBig)
	//fmt.Println("buf2=", strt.BodySizeBig)
	//fmt.Println("buf3=", strt.Context, "len(buf3)=", len(strt.Context))
	buf = append(buf, strt.BodySizeBig...)
	buf = append(buf, strt.Context...)
	return buf
}
func probeDateZip(context []byte) []byte {
	//写文件
	midName := GetGuid()
	dirData := "./data/data/"
	dirZip := "./data/zip/"
	originF := midName + ".data"
	zipF := midName + ".zip"

	err := ioutil.WriteFile(dirData+originF, context, 0644)
	if err != nil {
		fmt.Println("WriteFile failed!")
		return []byte("")
	}
	//压缩
	fzip, _ := os.Create(dirZip + zipF)
	w := zip.NewWriter(fzip)
	//fw, _ := w.Create("zip.data")
	fw, _ := w.Create(originF)
	filecontent, err := ioutil.ReadFile(dirData + originF)
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	n, err := fw.Write(filecontent)
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	fmt.Println(n)
	w.Close()
	if err := os.Remove(dirData + originF); err != nil {
		fmt.Println("remove file(", dirData+originF, ")failed!")
	}
	fzip.Close()
	//读文件

	dat, err := ioutil.ReadFile(dirZip + zipF)
	if err != nil {
		fmt.Println("err:", err)
		return []byte("")
	}
	if err := os.Remove(dirZip + zipF); err != nil {
		fmt.Println("remove file(", dirZip+zipF, ")failed!")
	}
	fmt.Println("dat:", dat, "len=", len(dat))
	return dat
}
func probeDesEncry(buf, key, iv []byte) ([]byte, error) {
	datalen := make([]byte, 4)
	binary.BigEndian.PutUint32(datalen, uint32(len(buf)))
	buf = append(buf, datalen...)

	pading := []byte{'A', 'B', 'C', 'D'}
	buf = append(buf, pading...)
	fmt.Println("context before des:", buf, "len=", len(buf))
	return DesEncrypt(buf, key, iv)
}
