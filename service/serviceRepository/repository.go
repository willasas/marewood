package serviceRepository

import (
	"FEDeployService/config"
	"FEDeployService/helper"
	"errors"
	"os"
	"os/exec"
	"strings"
)


//克隆仓库，userName，password可留空
func GitClone(gitUrl string, userName string, password string) (string, error) {

	var cmd *exec.Cmd

	if userName == "" || password == "" {
		cmd = exec.Command("git", "clone", gitUrl)
	} else {
		authUrl, err := helper.GitUrl2AuthUrl(gitUrl, userName, password)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("git", "clone", authUrl)
	}

	cmd.Dir = config.Cfg.RepositoryDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

//通过仓库url更新仓库
func GitPullByRepositoryUrl(url string) (string, error) {

	return runCmdOnRepositoryDir(url,"git", "pull")
}

func DeleteRepository(url string) error {

	repositoryName, err := helper.GetRepositoryNameByUrl(url)
	if err != nil {
		return err
	}

	repoDir := config.Cfg.RepositoryDir + "/" + repositoryName

	file, err := os.Stat(repoDir)

	if err != nil || !file.IsDir() {
		return  errors.New("找不到仓库目录,无法删除仓库")
	}

	removeErr := os.RemoveAll(repoDir)

	if removeErr != nil {
		return removeErr
	}
	return  nil
}

func GetBranchByRepositoryUrl(url string) ([]string, error) {

	out, err := runCmdOnRepositoryDir(url,"git", "branch","-r")

	if err != nil {
		return []string{}, err
	}
	branch := strings.Split(strings.Trim(strings.ReplaceAll(string(out)," ",""),"\n"), "\n")

	return branch[1:], nil
}

func GitCheckout(url string ,branch string) (string, error) {

	return runCmdOnRepositoryDir(url,"git", "checkout", branch)
}

//仓库URL， 构建命令 test、build、build:dev
func RunBuild(url string , buildCmd string) (string, error) {

	return runCmdOnRepositoryDir(url,"npm", "run", buildCmd)
}

func runCmdOnRepositoryDir(repositoryUrl string ,cmdName string, arg ...string) (string, error) {

	repositoryName, err := helper.GetRepositoryNameByUrl(repositoryUrl)
	if err != nil {
		return "", err
	}

	repositoryDir := config.Cfg.RepositoryDir + "/" + repositoryName

	file, err := os.Stat(repositoryDir)

	if err != nil || !file.IsDir() {
		return "", errors.New("找不到仓库目录=>" + repositoryDir)
	}

	cmd := exec.Command(cmdName , arg ...)
	cmd.Dir = repositoryDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}