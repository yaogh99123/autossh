package app

import (
	"autossh/src/i18n"
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

type Download struct {
	isDir bool
	cfg   *Config
	maxConcurrency int
}

// 简化的下载功能
func showDownload(configFile string) {
	var err error
	cfg, err := loadConfig(configFile)
	if err != nil {
		utils.Errorln(err)
		return
	}

	download := Download{
		cfg: cfg,
		maxConcurrency: 5, // 默认最大并发数
	}
	if err := download.parse(); err != nil {
		utils.Errorln(err)
		return
	}

	// 解析参数
	os.Args = flag.Args()
	flag.BoolVar(&download.isDir, "r", false, "文件夹")
	flag.IntVar(&download.maxConcurrency, "j", 5, "并发数")
	flag.Parse()

	var args = flag.Args()

	if len(args) < 2 {
		utils.Errorln(i18n.T("down_usage"))
		utils.Errorln(i18n.T("down_example1"))
		utils.Errorln(i18n.T("down_example2"))
		utils.Errorln(i18n.T("down_example3"))
		return
	}

	// 解析远程源
	remoteSource := args[0]
	localTarget := args[1]

	// 解析远程源
	parts := strings.Split(remoteSource, ":")
	if len(parts) != 2 {
		utils.Errorln(i18n.T("down_err_format"))
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
		utils.Errorln(i18n.T("err_server_not_found", serverName))
		return
	}

	// 建立SFTP连接
	sftpClient, err := serverIndex.server.GetSftpClient()
	if err != nil {
		utils.Errorln(i18n.T("err_conn_fail", err))
		return
	}
	defer func() {
		_ = sftpClient.Close()
	}()

	// 创建IO客户端
	srcIOClient := &SftpIOClient{SftpClient: sftpClient}
	dstIOClient := &LocalIOClient{}

	// 检查远程文件是否存在
	remoteFileInfo, err := srcIOClient.Stat(remotePath)
	if err != nil {
		utils.Errorln(i18n.T("down_err_remote_not_found", remotePath))
		return
	}

	// 如果是目录但未指定-r参数
	if remoteFileInfo.IsDir() && !download.isDir {
		utils.Errorln(i18n.T("down_err_is_dir"))
		return
	}

	// 执行下载
	if err := download.downloadFile(srcIOClient, dstIOClient, remotePath, localTarget); err != nil {
		utils.Errorln(i18n.T("down_err_download", err))
		return
	}

	fmt.Println(i18n.T("down_success"))
}

// 下载文件或目录
func (d *Download) downloadFile(srcIO IOClient, dstIO IOClient, srcPath string, dstPath string) error {
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
		return d.downloadDirectory(srcIO, dstIO, srcPath, dstPath)
	} else {
		// 处理单个文件
		return d.downloadSingleFile(srcIO, dstIO, srcPath, dstPath)
	}
}

// 下载单个文件
func (d *Download) downloadSingleFile(srcIO IOClient, dstIO IOClient, srcPath string, dstPath string) error {
	// 打开源文件
	srcFile, err := srcIO.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	// 确定目标路径
	finalDstPath := d.determineDestinationPath(dstIO, srcPath, dstPath)

	// 确保目标目录存在
	targetDir := filepath.Dir(finalDstPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// 创建目标文件
	dstFile, err := dstIO.Create(finalDstPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()

	// 复制文件内容
	_, err = d.copyFileContent(srcFile, dstFile, path.Base(srcPath))
	return err
}

// 下载目录
func (d *Download) downloadDirectory(srcIO IOClient, dstIO IOClient, srcPath string, dstPath string) error {
	// 读取源目录内容
	entries, err := srcIO.ReadDir(srcPath)
	if err != nil {
		return err
	}

	// 创建目标子目录
	dirName := path.Base(srcPath)
	targetDir := filepath.Join(dstPath, dirName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
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

	// 并行下载文件
	if len(files) > 0 {
		d.downloadFilesParallel(srcIO, dstIO, srcPath, targetDir, files)
	}

	// 递归下载目录（串行，避免过多并发连接）
	for _, dir := range dirs {
		srcEntryPath := path.Join(srcPath, dir.Name())
		if err := d.downloadDirectory(srcIO, dstIO, srcEntryPath, targetDir); err != nil {
			fmt.Printf(i18n.T("down_err_dir_fail", srcEntryPath, err))
		}
	}

	return nil
}

// 并行下载文件
func (d *Download) downloadFilesParallel(srcIO IOClient, dstIO IOClient, srcPath string, targetDir string, files []os.FileInfo) {
	// 创建信号量控制并发数
	semaphore := make(chan struct{}, d.maxConcurrency)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(file os.FileInfo) {
			defer wg.Done()
			
			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			srcEntryPath := path.Join(srcPath, file.Name())
			if err := d.downloadSingleFile(srcIO, dstIO, srcEntryPath, targetDir); err != nil {
				fmt.Printf(i18n.T("down_err_file_fail", srcEntryPath, err))
			}
		}(file)
	}

	wg.Wait()
}

// 确定最终的目标路径
func (d *Download) determineDestinationPath(dstIO IOClient, srcPath string, dstPath string) string {
	// 检查目标路径是否存在
	dstInfo, err := dstIO.Stat(dstPath)
	if err != nil {
		// 目标路径不存在，直接使用
		return dstPath
	}

	if dstInfo.IsDir() {
		// 目标是目录，在目录下创建同名文件
		return filepath.Join(dstPath, path.Base(srcPath))
	}

	// 目标是文件，直接覆盖
	return dstPath
}

// 复制文件内容（带进度显示）
func (d *Download) copyFileContent(srcFile FileLike, dstFile FileLike, filename string) (int64, error) {
	// 获取源文件大小
	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return 0, err
	}

	fileSize := srcFileInfo.Size()
	bytesCopied := int64(0)
	buffer := make([]byte, 64*1024) // 64KB buffer

	fmt.Printf("正在下载: %s (%.2f MB)\n", filename, float64(fileSize)/(1024*1024))

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
	fmt.Printf("\r下载完成: %s (100%% - %d bytes)\n", filename, bytesCopied)
	return bytesCopied, nil
}

// 解析参数（占位符，保持接口一致）
func (d *Download) parse() error {
	return nil
}
