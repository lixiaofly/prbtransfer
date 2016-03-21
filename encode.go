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
	fmt.Println("buf1=", strt.DataTypeBig)
	fmt.Println("buf2=", strt.BodySizeBig)
	fmt.Println("buf3=", strt.Context, "len(buf3)=", len(strt.Context))
	buf = append(buf, strt.BodySizeBig...)
	buf = append(buf, strt.Context...)
	fmt.Printf("\n")
	//fmt.Println("context before des:", buf, "len=", len(buf))
	return buf
}
func probeDateZip(context []byte) []byte {
	/*
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	//w.Write([]byte("hello, world\n"))
	w.Write(context)
	w.Close()
	return b.Bytes()
	*/
	//fmt.Println(b.Bytes())
	/*
			var buf bytes.Buffer
			fmt.Println("before zip :", context, "len=", len(context))
		    compressor, err := zlib.NewWriterLevelDict(&buf, zlib.BestCompression, context)
		    if err != nil {
		        fmt.Println("压缩失败")
		        return nil
		    }
		    compressor.Write(context)
		    compressor.Close()
			fmt.Println("after zip :", buf.Bytes(), "len=", len(buf.Bytes()))
			return buf.Bytes()
	*/
		//写文件
		err := ioutil.WriteFile("./data/zip.data", context, 0644)
		if err != nil {
			fmt.Println("WriteFile failed!")
			return []byte("")
		}
		//压缩
		const dir = "./data/"
			//获取源文件列表
		f, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
		}
		fzip, _ := os.Create("img-50.zip")
		w := zip.NewWriter(fzip)
		for _, file := range f {
			fw, _ := w.Create(file.Name())
			filecontent, err := ioutil.ReadFile(dir + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			n, err := fw.Write(filecontent)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(n)
		}
		w.Close()
	//读文件
	
		dat, err := ioutil.ReadFile("./img-50.zip")
		if err != nil {
			fmt.Println("err:", err)
			return []byte("")
		}
		fmt.Println("dat:", dat, "len=", len(dat))
		return dat
	
	/*
		const File = "img-50.zip"
		const dir_2 = "./"
		//os.Mkdir(dir, 0777) //创建一个目录

		cf, err := zip.OpenReader(File) //读取zip文件
		if err != nil {
			fmt.Println(err)
		}
		
		for _, file := range cf.File {
			rc, err := file.Open()
			if err != nil {
				fmt.Println(err)
			}

			f, err := os.Create(dir_2 + file.Name)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()
			buf := make([]byte, 100)
			rc.Read(buf)
			n, err := io.Copy(f, rc)
			if err != nil {
				fmt.Println(err)
			}
			
			fmt.Println("rc.Read=",buf)
			fmt.Println(n)
		}
		cf.Close()
		//fmt.Println("rc.Read=", rc.Read())
		return []byte("")
		*/
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
