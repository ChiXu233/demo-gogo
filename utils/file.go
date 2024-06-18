package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func GetFiles(pathName string, recursive bool, f *[]string) error {
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			if recursive == true {
				err = GetFiles(pathName+"/"+fi.Name(), recursive, f)
				if err != nil {
					return err
				}
			}
		} else {
			*f = append(*f, pathName+"/"+fi.Name())
		}
	}
	return nil
}

func ReadExcel(excelPath string) (*[][]string, error) {
	var content [][]string
	xlsxRead, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}
	sheet := xlsxRead.GetActiveSheetIndex()
	name := xlsxRead.GetSheetName(sheet)
	// log.Println("sheet name is " + name)
	content = xlsxRead.GetRows(name)
	// log.Println(content)
	return &content, nil
}

func Zip(zipPath string, paths ...string) error {
	if filepath.Ext(zipPath) != ".zip" {
		zipPath += ".zip"
	}
	// 创建zip文件及其父目录。
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// 创建zip写入器。
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// 遍历文件或目录。
	for _, rootPath := range paths {
		// 如果路径是目录，则删除尾部的路径分隔符。
		rootPath = strings.TrimSuffix(rootPath, string(os.PathSeparator))

		// 遍历树中的所有文件或目录。
		err = filepath.Walk(rootPath, walkFunc(rootPath, zipWriter))
		if err != nil {
			return err
		}
	}
	return nil
}

func walkFunc(rootPath string, zipWriter *zip.Writer) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果文件是符号链接，则跳过。
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// 创建一个本地文件头。
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置压缩方法。
		header.Method = zip.Deflate

		// 将文件的相对路径设置为头部名称。
		header.Name, err = filepath.Rel(filepath.Dir(rootPath), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}

		// 为文件头创建写入器，并保存文件的内容。
		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(headerWriter, f)
		return err
	}
}

func Tar(source string) error {
	// tar --directory=files/source_data_export/20220324-121511 -cvf files/source_data_export/20220324-121511/0.tar 0/
	splitResult := strings.Split(source, "/")
	basePath := strings.Join(splitResult[0:len(splitResult)-1], "/")
	dirName := splitResult[len(splitResult)-1]

	cmd := exec.Command("tar", "--directory", basePath, "-cvf", source+".tar", dirName)

	//阻塞至等待命令执行完成
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		filePath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Mkdir 创建文件夹
func Mkdir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			err = os.MkdirAll(path, 0777)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GinFileWrite(srcFile *multipart.FileHeader, dstFile string) error {
	source, err := srcFile.Open()
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func Copy(srcFile, dstFile string) error {
	sourceFileStat, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcFile)
	}

	source, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer source.Close()

	err = os.MkdirAll(filepath.Dir(dstFile), 0777)
	if err != nil {
		return err
	}

	destination, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func Remove(filePath string) error {
	var err error
	if !Exists(filePath) {
		err = fmt.Errorf("路径不存在")
		return err
	}
	err = os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}

// BatchRemove 批量删除文件，删除失败继续执行
func BatchRemove(filePaths []string) error {
	var errFileExists []string
	var errRemove []string
	for _, filePath := range filePaths {
		if !Exists(filePath) {
			errFileExists = append(errFileExists, filePath)
			continue
		}
		err := os.RemoveAll(filePath)
		if err != nil {
			errRemove = append(errRemove, filePath)
		}
	}
	if len(errFileExists) != 0 || len(errRemove) != 0 {
		return fmt.Errorf("删除文件失败 路径不存在[%v] 删除失败[%v]", errFileExists, errRemove)
	}
	return nil
}

func WriteJson(data interface{}, path string) error {
	// jsonByte, err := json.Marshal(data)
	jsonByte, err := json.MarshalIndent(data, "", "	")
	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(jsonByte)
	if err != nil {
		return err
	}
	return nil
}

func ReadJson(path string, dest interface{}) error {
	//打开文件
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	//读取为[]bytes类型
	byteData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteData, &dest)
	if err != nil {
		return err
	}

	return nil
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FormatFileSize 将文件大小转换为合适的单位
func FormatFileSize(size int64) string {
	const (
		_          = iota             // ignore first value by assigning to blank identifier
		KB float64 = 1 << (10 * iota) // 1 << (10*1) = 1024
		MB                            // 1 << (10*2) = 1048576
		GB                            // 1 << (10*3) = 1073741824
		TB                            // 1 << (10*4) = 1099511627776
	)

	fileSize := float64(size)

	switch {
	case fileSize >= TB:
		return fmt.Sprintf("%.2f TB", fileSize/TB)
	case fileSize >= GB:
		return fmt.Sprintf("%.2f GB", fileSize/GB)
	case fileSize >= MB:
		return fmt.Sprintf("%.2f MB", fileSize/MB)
	case fileSize >= KB:
		return fmt.Sprintf("%.2f KB", fileSize/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
