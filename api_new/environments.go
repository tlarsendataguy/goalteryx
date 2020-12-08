package api_new

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"unsafe"
)

type testEnvironment struct {
	sharedMemory *goPluginSharedMemory
	updateOnly   bool
	updateMode   string
	workflowDir  string
	locale       string
}

func (e *testEnvironment) UpdateOnly() bool {
	return e.updateOnly
}

func (e *testEnvironment) UpdateMode() string {
	return e.updateMode
}

func (e *testEnvironment) DesignerVersion() string {
	return `TestHarness`
}

func (e *testEnvironment) WorkflowDir() string {
	return e.workflowDir
}

func (e *testEnvironment) AlteryxInstallDir() string {
	return ``
}

func (e *testEnvironment) AlteryxLocale() string {
	return e.locale
}

func (e *testEnvironment) ToolId() int {
	return int(e.sharedMemory.toolId)
}

func (e *testEnvironment) UpdateToolConfig(newConfig string) {
	updateConfig(e.sharedMemory, newConfig)
}

type ayxEnvironment struct {
	sharedMemory *goPluginSharedMemory
}

func (e *ayxEnvironment) UpdateOnly() bool {
	initVar := getInitVarToEngine(e.sharedMemory, `UpdateOnly`)
	return initVar == `True`
}

func (e *ayxEnvironment) UpdateMode() string {
	return getInitVarToEngine(e.sharedMemory, `UpdateMode`)
}

func (e *ayxEnvironment) DesignerVersion() string {
	return getInitVarToEngine(e.sharedMemory, `Version`)
}

func (e *ayxEnvironment) WorkflowDir() string {
	return getInitVarToEngine(e.sharedMemory, `DefaultDir`)
}

func (e *ayxEnvironment) AlteryxInstallDir() string {
	return getInitVarToEngine(e.sharedMemory, `RuntimeDataPath`)
}

func (e *ayxEnvironment) AlteryxLocale() string {
	version := e.DesignerVersion()[:6]
	return getLocale(version)
}

func (e *ayxEnvironment) ToolId() int {
	return int(e.sharedMemory.toolId)
}

func (e *ayxEnvironment) UpdateToolConfig(newConfig string) {
	sendMessageToEngine(e.sharedMemory, UpdateOutputMetaInfoXml, newConfig)
	updateConfig(e.sharedMemory, newConfig)
}

type localeData struct {
	HelpLanguage string `xml:"GloablSettings>HelpLanguage"`
}

func updateConfig(sharedMemory *goPluginSharedMemory, newConfig string) {
	newConfigPtr := unsafe.Pointer(stringToUtf16Ptr(newConfig))
	newConfigLen := utf16PtrLen(newConfigPtr)
	sharedMemory.toolConfig = newConfigPtr
	sharedMemory.toolConfigLen = uint32(newConfigLen)
}

func getLocale(version string) string {
	settingsPath := filepath.Join(os.Getenv(`APPDATA`), `Alteryx`, `Engine`, version, `UserSettings.xml`)
	settingsBytes, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		return err.Error()
	}
	locale := localeData{}
	err = xml.Unmarshal(settingsBytes, &locale)
	if err != nil {
		return err.Error()
	}
	return locale.HelpLanguage
}
