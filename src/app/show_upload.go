package app

import (
	"autossh/src/utils"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Upload struct {
	isDir          bool
	cfg            *Config
	maxConcurrency int
}

// 简化的上传功能
func showUpload(configFile string) {
	var err error
	cfg, err := loadConfig(configFile)
	if err != nil {
		utils.Errorln(err)
		return
	}

	upload := Upload{
		cfg:            cfg,
		maxConcurrency: 5, // 默认最大并发数
	}
	if err := upload.parse(); err != nil {
		utils.Errorln(err)
		return
	}

	// 解析参数
	os.Args = flag.Args()
	flag.BoolVar(&upload.isDir, "r", false, "文件夹")
	flag.IntVar(&upload.maxConcurrency, "j", 5, "并发数")
	flag.Parse()

	var args = flag.Args()

	if len(args) < 2 {
		utils.Errorln("用法: autossh upload/up [-r] [-j 并发数] <本地文件/目录> <服务器名/序号:远程路径>")
		utils.Errorln("示例: autossh upload file.txt server1:/home/user/")
		utils.Errorln("示例: autossh up file.txt 01:/home/user/  (使用序号)")
		utils.Errorln("示例: autossh upload -r -j 10 ./localdir 01:/home/user/  (10个并发)")
		return
	}

	// 解析本地路径
	localPath := args[0]
	remoteTarget := args[1]

	// 检查本地文件是否存在
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		utils.Errorln("本地文件不存在: " + localPath)
		return
	}

	// 解析远程目标
	parts := strings.Split(remoteTarget, ":")
	if len(parts) != 2 {
		utils.Errorln("远程目标格式错误，应为: 服务器名/序号:路径")
		return
	}

	serverName := parts[0]
	remotePath := parts[1]

	// 查找服务器（支持序号和名称）
	var serverIndex *ServerIndex
	var exists bool

	// 首先尝试按名称查找
	if idx, found := cfg.serverIndex[serverName]; found {
		serverIndex = &idx
		exists = true
	}

	// 如果按名称没找到，尝试按序号查找
	if !exists {
		// 检查是否是数字序号
		if serverNum, err := strconv.Atoi(serverName); err == nil {
			// 序号从1开始，数组从0开始
			if serverNum > 0 && serverNum <= len(cfg.Servers) {
				serverIndex = &ServerIndex{
					server:      cfg.Servers[serverNum-1],
					serverIndex: serverNum - 1,
				}
				exists = true
			}
		}
	}

	if !exists {
		utils.Errorln("服务器 " + serverName + " 不存在")
		return
	}

	// 建立SFTP连接
	sftpClient, err := serverIndex.server.GetSftpClient()
	if err != nil {
		utils.Errorln("连接服务器失败: " + err.Error())
		return
	}
	defer func() {
		_ = sftpClient.Close()
	}()

	// 创建IO客户端
	srcIOClient := &LocalIOClient{}
	dstIOClient := &SftpIOClient{SftpClient: sftpClient}

	// 获取本地文件信息
	localFileInfo, err := os.Stat(localPath)
	if err != nil {
		utils.Errorln("获取本地文件信息失败: " + err.Error())
		return
	}

	// 如果是目录但未指定-r参数
	if localFileInfo.IsDir() && !upload.isDir {
		utils.Errorln("本地路径是目录，请使用 -r 参数")
		return
	}

	// 执行上传
	if err := upload.uploadFile(srcIOClient, dstIOClient, localPath, remotePath); err != nil {
		utils.Errorln("上传失败: " + err.Error())
		return
	}

	fmt.Println("上传完成!")
}

// 上传文件或目录
func (u *Upload) uploadFile(srcIO IOClient, dstIO IOClient, srcPath string, dstPath string) error {
	// 打开源文件
	srcFile, err := srcIO.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// 获取源文件信息
	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if srcFileInfo.IsDir() {
		// 处理目录
		return u.uploadDirectory(srcIO, dstIO, srcPath, dstPath)
	} else {
		// 处理单个文件
		return u.uploadSingleFile(srcIO, dstIO, srcPath, dstPath)
	}
}

// 上传单个文件
func (u *Upload) uploadSingleFile(srcIO IOClient, dstIO IOClient, srcPath string, dstPath string) error {
	// 打开源文件
	srcFile, err := srcIO.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// 确定目标路径
	finalDstPath := u.determineDestinationPath(dstIO, srcPath, dstPath)

	// 创建目标文件
	dstFile, err := dstIO.Create(finalDstPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	// 复制文件内容
	_, err = u.copyFileContent(srcFile, dstFile, path.Base(srcPath))
	return err
}

// 上传目录
func (u *Upload) uploadDirectory(srcIO IOClient, dstIO IOClient, srcPath string, dstPath string) error {
	// 确保目标目录存在
	if err := u.ensureDirectoryExists(dstIO, dstPath); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := srcIO.ReadDir(srcPath)
	if err != nil {
		return err
	}

	// 创建目标子目录
	dirName := filepath.Base(srcPath)
	targetDir := path.Join(dstPath, dirName)
	if err := u.ensureDirectoryExists(dstIO, targetDir); err != nil {
		return err
	}

	// 分离目录和文件
	var dirs []os.FileInfo
	var files []os.FileInfo

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	// 并行上传文件
	if len(files) > 0 {
		u.uploadFilesParallel(srcIO, dstIO, srcPath, targetDir, files)
	}

	// 递归上传目录（串行，避免过多并发连接）
	for _, dir := range dirs {
		srcEntryPath := filepath.Join(srcPath, dir.Name())
		if err := u.uploadDirectory(srcIO, dstIO, srcEntryPath, targetDir); err != nil {
			fmt.Printf("上传目录 %s 失败: %v\n", srcEntryPath, err)
		}
	}

	return nil
}

// 并行上传文件
func (u *Upload) uploadFilesParallel(srcIO IOClient, dstIO IOClient, srcPath string, targetDir string, files []os.FileInfo) {
	// 创建信号量控制并发数
	semaphore := make(chan struct{}, u.maxConcurrency)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			srcEntryPath := filepath.Join(srcPath, file.Name())
			if err := u.uploadSingleFile(srcIO, dstIO, srcEntryPath, targetDir); err != nil {
				fmt.Printf("上传文件 %s 失败: %v\n", srcEntryPath, err)
			}
		}(file)
	}

	wg.Wait()
}

// 确定最终的目标路径
func (u *Upload) determineDestinationPath(dstIO IOClient, srcPath string, dstPath string) string {
	// 检查目标路径是否存在
	dstInfo, err := dstIO.Stat(dstPath)
	if err != nil {
		// 目标路径不存在，直接使用
		return dstPath
	}

	if dstInfo.IsDir() {
		// 目标是目录，在目录下创建同名文件
		return path.Join(dstPath, filepath.Base(srcPath))
	}

	// 目标是文件，直接覆盖
	return dstPath
}

// 确保目录存在
func (u *Upload) ensureDirectoryExists(dstIO IOClient, dirPath string) error {
	_, err := dstIO.Stat(dirPath)
	if err != nil {
		// 目录不存在，尝试创建
		return dstIO.Mkdir(dirPath)
	}
	return nil
}

// 复制文件内容（带进度显示）
func (u *Upload) copyFileContent(srcFile FileLike, dstFile FileLike, filename string) (int64, error) {
	// 获取源文件大小
	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return 0, err
	}

	fileSize := srcFileInfo.Size()
	bytesCopied := int64(0)
	buffer := make([]byte, 64*1024) // 64KB buffer

	fmt.Printf("正在上传: %s (%.2f MB)\n", filename, float64(fileSize)/(1024*1024))

	// 启动进度显示协程
	done := make(chan bool)
	startTime := time.Now()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if fileSize > 0 {
					percentage := float64(bytesCopied) / float64(fileSize) * 100
					elapsed := time.Since(startTime).Seconds()
					if elapsed > 0 {
						speed := float64(bytesCopied) / elapsed / 1024 / 1024 // MB/s
						fmt.Printf("\r进度: %.1f%% (%d/%d bytes) %.2f MB/s",
							percentage, bytesCopied, fileSize, speed)
					} else {
						fmt.Printf("\r进度: %.1f%% (%d/%d bytes)",
							percentage, bytesCopied, fileSize)
					}
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	for {
		n, err := srcFile.Read(buffer)
		if n > 0 {
			written, writeErr := dstFile.Write(buffer[:n])
			if writeErr != nil {
				done <- true
				return bytesCopied, writeErr
			}
			bytesCopied += int64(written)
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			done <- true
			return bytesCopied, err
		}
	}

	done <- true
	fmt.Printf("\r上传完成: %s (100%% - %d bytes)\n", filename, bytesCopied)
	return bytesCopied, nil
}

// 解析参数（占位符，保持接口一致）
func (u *Upload) parse() error {
	return nil
}
