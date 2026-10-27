package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	tempDir    string
	ytdlpPath  string
	ffmpegPath string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 建立專屬暫存目錄
	dir, err := os.MkdirTemp("", "govi2-*")
	if err != nil {
		fmt.Println("建立暫存目錄失敗:", err)
		return
	}

	a.tempDir = dir

	// 解壓縮 yt-dlp 和 ffmpeg
	for _, name := range []string{"yt-dlp.exe", "ffmpeg.exe"} {

		data, err := embeddedTools.ReadFile("bin/" + name)
		if err != nil {
			fmt.Println("讀取失敗:", err)
			return
		}

		path := filepath.Join(dir, name)

		if err := os.WriteFile(path, data, 0700); err != nil {
			fmt.Println("寫入失敗:", err)
			return
		}

		switch name {
		case "yt-dlp.exe":
			a.ytdlpPath = path
		case "ffmpeg.exe":
			a.ffmpegPath = path
		}
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.tempDir != "" {
		os.RemoveAll(a.tempDir)
	}
}

func (a *App) RunTool(name string, args []string) (string, error) {

	// 設定執行檔
	cmd := exec.Command(name, args...)

	// 指定工作目錄
	cmd.Dir = a.tempDir

	// 將內嵌工具目錄加入 PATH
	cmd.Env = append(
		os.Environ(),
		"PATH="+a.tempDir+string(os.PathListSeparator)+os.Getenv("PATH"),
	)

	// 執行
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf(
			"執行失敗: %w\n%s",
			err,
			string(output),
		)
	}

	return string(output), nil
}

func (a *App) Test1() string {

	if a.ytdlpPath == "" || a.ffmpegPath == "" {
		return "yt-dlp 或 ffmpeg 尚未初始化"
	}

	// 取得下載目錄
	home, err := os.UserHomeDir()
	if err != nil {
		return err.Error()
	}

	downloadDir := filepath.Join(home, "Downloads")

	res, err := a.RunTool(a.ytdlpPath, []string{
		"--verbose",
		"--ffmpeg-location", a.tempDir,
		"-P", downloadDir,
		"https://www.youtube.com/watch?v=SC-got3RjtY",
	})

	// 將輸出寫入檔案，方便 Finder 啟動時查看
	logPath := filepath.Join(home, "govi2-debug.log")

	logContent := fmt.Sprintf(
		"Error: %v\nOutput:\n%s\n",
		err,
		res,
	)

	if logErr := os.WriteFile(
		logPath,
		[]byte(logContent),
		0600,
	); logErr != nil {
		fmt.Println("寫入日誌失敗:", logErr)
	}

	if err != nil {
		return fmt.Sprintf("下載失敗: %v\n%s", err, res)
	}

	return "OK"
}

func (a *App) Test2(f any) string {
	fmt.Printf("%T", f)
	return "11"
}

func (a *App) SelectAndProcessFile() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "請選擇要上傳的檔案",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "圖片與文字檔 (*.png;*.jpg;*.txt)",
				//Pattern:     "*.png;*.jpg;*.txt", // 限制檔案類型，若不限制可不填
			},
		},
	}
	filePath, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return "", fmt.Errorf("開啟檔案視窗失敗: %w", err)
	}
	if filePath == "" {
		return "使用者取消了選擇", nil
	}
	fmt.Println("成功拿到絕對路徑：", filePath)
	return filePath, nil
}

func (a *App) SaveFileDialog() (string, error) {
	options := runtime.SaveDialogOptions{
		Title: "請選擇存放位置",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "media",
				Pattern:     "*.gif;", // 限制檔案類型，若不限制可不填
			},
			{
				DisplayName: "media",
				Pattern:     "*.mp4;", // 限制檔案類型，若不限制可不填
			},
		},
	}
	filePath, err := runtime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return "", fmt.Errorf("開啟檔案視窗失敗: %w", err)
	}
	if filePath == "" {
		return "使用者取消了選擇", nil
	}
	fmt.Println("成功拿到絕對路徑：", filePath)
	return filePath, nil
}

func (a *App) GetFileInfo(path string) map[string]any {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return nil
	}
	fileOpen, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer fileOpen.Close()
	buffer := make([]byte, 512)
	_, err = fileOpen.Read(buffer)
	mimeType := http.DetectContentType(buffer)
	res := map[string]any{
		"name":     fileInfo.Name(),
		"size":     fileInfo.Size(),
		"mode":     fileInfo.Mode(),
		"modTime":  fileInfo.ModTime(),
		"mimeType": mimeType,
		"path":     path,
	}
	return res
}

func (a *App) DownloadFileByPath(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *App) APIFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("APIFile")
	rawPath := r.URL.Query().Get("path")
	file, err := os.Open(rawPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		http.Error(w, "File error", 500)
		return
	}
	if !strings.EqualFold(
		r.Method, http.MethodGet,
	) {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set(
		"Content-Type",
		"video/mp4",
	)
	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		file,
	)
}

func (a *App) APIExecMedia(w http.ResponseWriter, r *http.Request) {
	fmt.Println("APIExecMedia")
	rawPath := r.URL.Query().Get("path")
	cmd := exec.CommandContext(
		r.Context(),
		a.ffmpegPath,
		"-i", rawPath,
		//"-vf", "setpts=0.5*PTS",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"-f", "mp4",
		"pipe:1",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	_, err = io.Copy(w, stdout)
	if err != nil {
		_ = cmd.Process.Kill()
	}
	_ = cmd.Wait()
}

func (a *App) ExportMedia(data map[string]any) bool {
	var vfFlag string = ""
	var vfString string = ""
	info, ok := data["info"].(map[string]any)
	if !ok {
		return false
	}
	timeRange, ok := data["timeRange"].([]any)
	if !ok {
		return false
	}
	exportConfig, ok := data["exportConfig"].(map[string]any)
	if !ok {
		return false
	}
	fps, _ := exportConfig["fps"].(float64)
	if fps != 0 {
		vfString += "fps=" + fmt.Sprintf("%g", fps)
	} else {
		vfString += ""
	}
	ratio, ok := exportConfig["ratio"].(float64)
	if ok {
		if len(vfString) != 0 {
			vfString += ","
		}
		vfString += "scale=-1:" + fmt.Sprintf("%g", ratio)
	} else {
		if len(vfString) != 0 {
			vfString += ","
		}
		vfString += "scale=-1:-1"
	}
	palette, _ := exportConfig["palette"].(bool)
	if palette {
		vfString += ":flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse"
	}
	if len(vfString) != 0 {
		vfFlag = "-vf"
	}
	fmt.Println(77777, vfString)
	savePath, savePathErr := a.SaveFileDialog()
	if savePathErr != nil {
		return false
	}
	_, err := a.RunTool(a.ffmpegPath, []string{
		"-i", info["path"].(string),
		"-ss", fmt.Sprintf("%g", timeRange[0]),
		"-to", fmt.Sprintf("%g", timeRange[1]),
		vfFlag, vfString,
		savePath,
	})
	if err != nil {
		return false
	}
	return true
}
