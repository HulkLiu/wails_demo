package internal

// Upload ...
func (a *App) Upload(filepath string, clis ...bool) (string, error) {
	a.Log.Infof("Upload filepath: %v", filepath)
	//fileExt := strings.ToLower(path.Ext(filepath))
	//fileSize := xfile.Size(filepath)
	//fileContent := xfile.Read(filepath)
	//fileMd5 := xcrypto.Md5(fileContent)
	//fileGitPath := fmt.Sprintf(config.GitFilePath, fileMd5[0:2], fileMd5, fileExt)
	// 请求上传文件
	//err := a.Git.Update(fileGitPath, fileContent)
	//if err != nil {
	//	a.Log.Infof("Upload Git err: %v", err)
	//	return "", err
	//}
	// 新增/更新数据
	//
	//fileInfo := &service.File{}
	//a.DB.Where("file_md5", fileMd5).First(fileInfo)
	//fileInfo.Name = path.Base(filepath)
	//fileInfo.Md5 = fileMd5
	//fileInfo.Size = xfile.SizeText(fileSize)
	//fileInfo.Path = fileGitPath
	//fileInfo.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	//a.Log.Infof("UploadFile fileInfo: %+v", fileInfo)
	//if err := a.DB.Save(fileInfo).Error; err != nil {
	//	a.Log.Errorf("UploadFile DB err: %v", err)
	//	return "", err
	//}
	//if len(clis) > 0 {
	//	////命令行调用
	//	//fileUrl := a.Git.GetFileUrl(fileInfo.Path)
	//	//a.Log.Infof("UploadFile success cli fileUrl: %s", fileUrl)
	//	//go a.SyncDatabase()
	//	//return fileUrl, nil
	//	return fileInfo.Path, nil
	//
	//}
	a.Log.Info("UploadFile success app")
	return "", nil
}
