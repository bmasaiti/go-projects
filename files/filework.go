package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

func main() {
	path := "example.txt"
	fileExist, _ := fileExists(path)
	fmt.Println("file exists:", fileExist)
	fmt.Println("file created:", createFile(path))
	fmt.Println("file created if not exist:", createFileIfNotExists(path))
	fmt.Println("write file:",writeFile(path ,[]byte ("some text to write")))
	fmt.Println("Appending to file :", appendFile(path, []byte("We append stuff to the file")))
	
	data, _ := readFile(path)
	fmt.Println("file content:", string(data))

	copiedPath := "copy.text"
	fmt.Println("copying file:", copyFile(path,copiedPath))


	fmt.Println(	"\nlisted files :")
	//var listedFiles []fs.DirEntry
	listedFiles,_ := listDirFiles("./") //list files in current dir
	fmt.Println(listedFiles)
	// for_, listedFile := range listedFiles{
	// 	fmt.Println(listedFile.Name())
	// }

}

func fileExists(path string) (bool, error) {

	_, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, nil
	}
	return true, nil

}

func createFile(path string) error {
	file, err := os.Create(path) //truncactes to zero if the file has  content already
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

func createFileIfNotExists(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_EXCL|os.O_CREATE, 0666) //open to read and write open to create and extend plus the permissions

	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644) //read write permissions for the owner , overwrites contents of file
}

func appendFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644) //open file in write only to append
	if err != nil {
		return err
	}
	defer file.Close()

	_, er := file.Write(data)
	return er

}

func readFile(path string) ([]byte, error) { //reads in the whole file(very inefficient)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func deleteFile(path string) error {
	return os.Remove(path)
}

func fileStats(path string) (int64, time.Time, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return -1, time.Time{}, err
	}
	size := info.Size()
	modTime := info.ModTime()
	return size, modTime, nil
}

func copyFile (srcPath , destPath string) error {
	//io.Copy(dst Writer , src Reader)
	src,err := os.Open(srcPath)
	defer src.Close()

	if err!=nil {
		return err
	}
	dest, err := os.Create(destPath)
	defer dest.Close()
	if err!=nil {
		return err
	}
	_,err = io.Copy(dest,src)
	if err!=nil{
		return  err
	}
	return dest.Sync()

}


func listDirFiles(dirPath string )([]fs.DirEntry, error){
	entries,err := os.ReadDir(dirPath)
	if err!=nil {
		return nil, err
	}
	listedFiles := make([]fs.DirEntry, 0)
	for _, entry := range entries {
		if !entry.IsDir(){
			listedFiles = append(listedFiles, entry)
		}
	}
	return listedFiles, nil
}