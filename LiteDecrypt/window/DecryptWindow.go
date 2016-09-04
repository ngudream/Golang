package window

import (
	"LiteDecrypt/util"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type DecryptWindow struct {
	*walk.MainWindow
	prevFilePath string
}

type FilePath struct {
	file string
	save string
}

const GridColumns int = 10
const TextEditColumnSpan int = 1
const PushButtonColumnSpan int = 9
const resetSecond time.Duration = 10

var progressLable *walk.Label
var copyResultLable *walk.Label
var copyResetTimer = time.NewTimer(0 * time.Second)
var copyFirstTime bool = true
var aboutMW = new(DecryptWindow)
var isExistAboutMW bool = false

var filePath FilePath //输入输出文件
var decryptError error

func (window *DecryptWindow) OpenDecryptWindow() {
	var zipEdit, saveZipEdit, unzipEdit, saveUnzipEdit *walk.TextEdit
	var zipBtn, unzipBtn, saveZipBtn, saveUnzipBtn, copyInBtn,
		startUnzipBtn, startZipBtn *walk.PushButton
	pathMW := new(DecryptWindow)
	pathMW.SetMaximizeBox(false) // 禁止最大化
	pathMW.SetMinimizeBox(true)  // 禁止最小化
	pathMW.SetFixedSize(false)   // 固定窗体大小
	err := MainWindow{
		AssignTo: &pathMW.MainWindow,
		Title:    "解压缩文件",
		MinSize:  Size{480, 500},
		Layout:   HBox{Spacing: 2},
		MenuItems: []MenuItem{
			Menu{
				Text: "帮助",
				Items: []MenuItem{
					Action{
						Text:        "关于LiteDecrypt",
						OnTriggered: pathMW.openHelperAbout,
					},
				},
			},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: GridColumns, Spacing: 10},
				Children: []Widget{
					VSplitter{
						ColumnSpan: TextEditColumnSpan,
						Children: []Widget{
							TextEdit{
								AssignTo:    &unzipEdit,
								Text:        filePath.file,
								MinSize:     Size{250, 0},
								ToolTipText: "请输入需要解压缩文件的路径",
							},
						},
					},
					VSplitter{
						ColumnSpan: PushButtonColumnSpan,
						Children: []Widget{
							PushButton{
								AssignTo: &unzipBtn,
								Text:     "选择解压文件",
								MinSize:  Size{130, 0},
								OnClicked: func() {
									path, err := pathMW.openFileManager()
									if err != nil {
										fmt.Println(err)
									}
									if len(path) > 0 {
										unzipEdit.SetText(path) //显示解压文件
									}
								},
							},
						},
					},
					VSplitter{
						ColumnSpan: TextEditColumnSpan,
						Children: []Widget{
							TextEdit{
								AssignTo:    &saveUnzipEdit,
								Text:        filePath.save,
								MinSize:     Size{250, 0},
								ToolTipText: "请输入要保存解压文件的路径",
							},
						},
					},
					VSplitter{
						ColumnSpan: PushButtonColumnSpan,
						Children: []Widget{
							PushButton{
								AssignTo: &saveUnzipBtn,
								Text:     "选择保存路径",
								MinSize:  Size{130, 0},
								OnClicked: func() {
									path, err := pathMW.openDirManager()
									if err != nil {
										fmt.Println(err)
									}
									if len(path) > 0 {
										saveUnzipEdit.SetText(path)
									}
								},
							},
						},
					},
					VSplitter{
						ColumnSpan: TextEditColumnSpan,
						Children: []Widget{
							TextEdit{
								AssignTo:    &zipEdit,
								Text:        filePath.file,
								MinSize:     Size{250, 0},
								ToolTipText: "请输入要压缩文件的路径",
							},
						},
					},
					VSplitter{
						ColumnSpan: PushButtonColumnSpan,
						Children: []Widget{
							PushButton{
								AssignTo: &zipBtn,
								Text:     "选择压缩文件",
								MinSize:  Size{130, 0},
								OnClicked: func() {
									path, err := pathMW.openFileManager()
									if err != nil {
										fmt.Println(err)
									}
									if len(path) > 0 {
										var result string = zipEdit.Text()
										if len(result) > 0 {
											zipEdit.SetText(result + ";" + path)
										} else {
											zipEdit.SetText(path)
										}
									}
								},
							},
						},
					},
					VSplitter{
						ColumnSpan: TextEditColumnSpan,
						Children: []Widget{
							TextEdit{
								AssignTo:    &saveZipEdit,
								Text:        filePath.save,
								MinSize:     Size{250, 0},
								ToolTipText: "请输入要保存压缩包的路径",
							},
						},
					},
					VSplitter{
						ColumnSpan: PushButtonColumnSpan,
						Children: []Widget{
							PushButton{
								AssignTo: &saveZipBtn,
								Text:     "选择保存路径",
								MinSize:  Size{130, 0},
								OnClicked: func() {
									path, err := pathMW.openDirManager()
									if err != nil {
										fmt.Println(err)
									}
									if len(path) > 0 {
										saveZipEdit.SetText(path)
									}
								},
							},
						},
					},
					VSplitter{
						ColumnSpan: GridColumns,
						Children: []Widget{
							Label{
								Text: "",
							},
						},
					},
					VSplitter{
						ColumnSpan: GridColumns,
						Children: []Widget{
							Label{
								AssignTo: &progressLable,
								Text:     "",
								//Font:     Font{Bold: true, PointSize: 13},
							},
						},
					},
					VSplitter{
						ColumnSpan: GridColumns,
						Children: []Widget{
							Label{
								AssignTo:           &copyResultLable,
								Row:                1,
								AlwaysConsumeSpace: true,
								Text:               "",
							},
						},
					},
					VSplitter{
						ColumnSpan: GridColumns,
						Children: []Widget{
							PushButton{
								AssignTo: &startUnzipBtn,
								Text:     "开始解压",
								MinSize:  Size{100, 30},
								OnClicked: func() {
									startUnzipBtn.SetEnabled(false)
									filePath.file = unzipEdit.Text()
									filePath.save = saveUnzipEdit.Text()
									finishChan := make(chan bool, 1) //用于接收结束标志
									go pathMW.startToUnZip(&filePath, finishChan)
									go pathMW.showResultMsgBox(true, startUnzipBtn, finishChan)
								},
							},
						},
					},
					VSplitter{
						ColumnSpan: GridColumns,
						Children: []Widget{
							PushButton{
								AssignTo: &startZipBtn,
								Text:     "开始压缩",
								MinSize:  Size{100, 30},
								OnClicked: func() {
									startZipBtn.SetEnabled(false)
									filePath.file = zipEdit.Text()
									filePath.save = saveZipEdit.Text()
									finishChan := make(chan bool, 1) //用于接收结束标志
									go pathMW.startToZip(&filePath, finishChan)
									go pathMW.showResultMsgBox(false, startZipBtn, finishChan)
								},
							},
						},
					},
					VSplitter{
						ColumnSpan: GridColumns,
						Children: []Widget{
							PushButton{
								AssignTo: &copyInBtn,
								Text:     "拷贝文件",
								MinSize:  Size{100, 10},
								OnClicked: func() {
									copyInBtn.SetEnabled(false)
									pathMW.copySelectedFile(unzipEdit.Text())
									copyInBtn.SetEnabled(true)
								},
							},
						},
					},
				},
			},
		},
	}.Create()
	if err != nil {
		fmt.Println(err)
	}
	pathMW.SetX(650)
	pathMW.SetY(300)
	pathMW.Run()
}

/**
*解压指定文件
 */
func (mw *DecryptWindow) startToUnZip(filePath *FilePath, ch chan<- bool) {
	var success bool = false
	defer func() {
		ch <- success
	}()
	reader, err := zip.OpenReader(filePath.file)
	if err != nil {
		return
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return
		}
		defer rc.Close()
		filename := filePath.save + file.Name
		newname, _ := utf8ToGBK(filename)
		if file.FileInfo().IsDir() { //如果是目录，则获取目录路径，主要是清除非法字符
			newpath := getDir(newname)
			if !util.IsFileExist(newpath) { //文件夹不存在就创建
				err = os.MkdirAll(newpath, 0755)
				fmt.Println("make a dir:" + newpath)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		if !file.FileInfo().IsDir() {
			w, err := os.Create(newname)
			fmt.Println("create a file:" + newname)
			if err != nil {
				return
			}
			defer w.Close()
			_, err = io.Copy(w, rc)
			if err != nil {
				return
			}
			w.Close()
		}
		rc.Close()
	}
	success = true
}

func utf8ToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewDecoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}

/**
*获取目录
 */
func getDir(path string) string {
	var index = strings.LastIndex(path, "/")
	var check bool = true
	if index == -1 {
		check = false
		index = strings.LastIndex(path, "\\")
		if index == -1 {
			return path
		}
	}
	var result = util.SubString(path, 0, index)
	for {
		var sep string = "/"
		if !check {
			sep = "\\"
		}
		if strings.HasSuffix(result, sep) {
			index--
			result = util.SubString(result, 0, index)
			fmt.Println("-> result:" + result)
		} else {
			break
		}
	}
	fmt.Println("path:" + path + "-> result:" + result + " index: " + strconv.Itoa(index))
	return result
}

/**
*压缩指定文件
 */
func (mw *DecryptWindow) startToZip(filePath *FilePath, ch chan<- bool) {
	var success bool = false
	defer func() {
		ch <- success
	}()
	list := strings.Split(filePath.file, ";")
	var length = len(list)
	var files [2]*os.File
	fmt.Println("length =" + strconv.Itoa(length))
	for i := 0; i < length; i++ {
		f, err := os.Open(list[i])
		if err != nil {
			fmt.Println(err)
		} else {
			files[i] = f
		}
		defer f.Close()
	}

	d, _ := os.Create(filePath.save)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {

		err := compress(file, "测试压缩", w)
		if err != nil {
			return
		}
	}
	success = true
}

/**
* 执行压缩
 */
func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

/**
*获取当前时间，返回拼接的字符串，用于解密文件名
 */
func getCurrentTime() string {
	timeStr := time.Now().Format("2006-09-04 11:04:05")
	timeStr = strings.Replace(timeStr, "-", "", -1)
	timeStr = strings.Replace(timeStr, " ", "-", -1)
	timeStr = strings.Replace(timeStr, ":", "", -1)
	return timeStr
}

/**
*弹出选文件的框
 */
func (mw *DecryptWindow) openFileManager() (filePath string, ror error) {
	dlg := new(walk.FileDialog)
	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Text Files(*.*)"
	dlg.Title = "选择文件"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return filePath, err
	} else if !ok {
		return filePath, nil
	}
	filePath = dlg.FilePath
	return filePath, nil
}

/**
*弹出选文件夹的框
 */
func (mw *DecryptWindow) openDirManager() (dirPath string, ror error) {
	dlg := new(walk.FileDialog)
	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Text Files(*.zip)"
	dlg.Title = "选择保存路径"

	if ok, err := dlg.ShowBrowseFolder(mw); err != nil {
		return dirPath, err
	} else if !ok {
		return dirPath, nil
	}
	dirPath = dlg.FilePath
	return dirPath, nil
}

/**
*打开帮助对话框
 */
func (mw *DecryptWindow) openHelperAbout() {
	if !isExistAboutMW || aboutMW.IsDisposed() {
		isExistAboutMW = true
		err := MainWindow{
			AssignTo: &aboutMW.MainWindow,
			Title:    "关于LiteDecrypt",
			MinSize:  Size{400, 200},
			Layout:   HBox{Spacing: 2},
			Children: []Widget{
				Composite{
					Layout: Grid{Columns: 1, Spacing: 5},
					Children: []Widget{
						VSplitter{
							ColumnSpan: 1,
							MaxSize:    Size{100, 20},
							Children: []Widget{
								Label{
									Text: "欢迎使用LiteDecrypt，它是一个简易的解压缩工具。",
								},
							},
						},
						VSplitter{
							ColumnSpan: 1,
							MaxSize:    Size{100, 20},
							Children: []Widget{
								Label{
									Text: "版本：0.1 beta",
								},
							},
						},
						VSplitter{
							ColumnSpan: 1,
							MaxSize:    Size{100, 20},
							Children: []Widget{
								Label{
									Text: "技术支持：ngudreamlee@gmail.com",
								},
							},
						},
						VSplitter{
							ColumnSpan: 1,
							MaxSize:    Size{100, 20},
							Children: []Widget{
								Label{
									Text: "Copyright © ngudream",
								},
							},
						},
					},
				},
			},
		}.Create()
		if err != nil {
			walk.MsgBox(mw, "提示信息", "打开关于界面失败", walk.MsgBoxIconError)
			return
		}
		aboutMW.SetX(600)
		aboutMW.SetY(420)
		aboutMW.Run()
	}
}

/**
*拷贝指定文件到指定目录
 */
func (mw *DecryptWindow) copySelectedFile(srcPath string) {
	if copyFirstTime {
		copyFirstTime = false
		copyResetTimer.Stop()
		<-copyResetTimer.C
	}
	destPath, err := mw.openDirManager()
	if err != nil {
		copyResultLable.SetText("拷贝失败！路径：" + destPath)
		goto finish
	}
	_, err, destPath = util.CopyFile(srcPath, destPath)
	if err != nil {
		copyResultLable.SetText("拷贝失败！路径：" + destPath)
		goto finish
	}
	copyResultLable.SetText("拷贝成功！路径：" + destPath)
finish:
	{
		copyResetTimer.Reset(resetSecond * time.Second)
		go func() {
			//定时器，指定时间后提示消息消失
			for {
				select {
				case <-copyResetTimer.C:
					{
						copyResultLable.SetText("")
					}
				}
			}
		}()
	}
}

func (mw *DecryptWindow) showResultMsgBox(unzip bool, decryptBtn *walk.PushButton, ch <-chan bool) {
	success := <-ch
	if success {
		var text string = "解压成功"
		if !unzip{
			text = "压缩成功"
		}
		progressLable.SetText(progressLable.Text() + "    " +　text)
	} else {
		var text string = "解压失败"
		if !unzip{
			text = "压缩失败"
		}
		var errStr string
		if decryptError != nil {
			errStr = decryptError.Error()
		}
		progressLable.SetText(progressLable.Text() + "    " + text + errStr)
		bg, err := walk.NewSolidColorBrush(walk.RGB(255, 0, 0)) //红色，现在不生效？？
		if err != nil {
			return
		}
		progressLable.SetBackground(bg)
	}
	decryptBtn.SetEnabled(true)
}
