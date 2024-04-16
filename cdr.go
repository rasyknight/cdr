package cdr

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/rasyknight/logg"
)

type CDR struct {
	CDRPrefix     string
	FileAliveTime int
	FileRecords   int
	f             *os.File
	id            int
	Records       int
	event         chan int
	filename      string
	NullFlag      int
}

func (c *CDR) CheckFile() {
	if c.Records >= c.FileRecords {
		c.f.Close()
		os.Rename(c.filename+".w", c.filename+".r")
		cdrFile := c.GetCDRFileName()
		var err error
		c.f, err = os.OpenFile(cdrFile+".w", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if nil != err {
			logg.Error("Open File %s.w failed, reason[%s.w].", cdrFile, err.Error())
		}
		logg.Debug("Close cdr file[%s.w], use new file [%s.w]", c.filename, cdrFile)
		c.filename = cdrFile
		c.id++
		c.Records = 0
	}
}
func (c *CDR) Wln(val ...any) {
	str := fmt.Sprintln(val...)
	c.CheckFile()
	c.Records++
	c.f.WriteString(str)
}

func (c *CDR) W(message string, args ...interface{}) {
	str := fmt.Sprintf(message, args...) + "\n"
	c.CheckFile()
	c.Records++
	c.f.WriteString(str)
}

func OpenCDR(cdrPrefix string, FileAliveTime, fileRecords int) CDR {
	cdr := CDR{
		CDRPrefix:     cdrPrefix,
		FileAliveTime: FileAliveTime,
		FileRecords:   fileRecords,
		id:            0,
	}
	cdrFile := cdr.GetCDRFileName()
	var err error
	cdr.filename = cdrFile
	cdr.f, err = os.OpenFile(cdrFile+".w", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		logg.Error("Open File %s.w failed, reason[%s].", cdrFile, err.Error())
	}
	logg.Debug("Open CDR file [%s.w]", cdrFile)
	cdr.id++
	return cdr
}

func (c *CDR) Close() {
	c.f.Close()
	if c.Records == 0 && 0 == c.NullFlag {
		os.Remove(c.filename + ".w")
	} else {
		os.Rename(c.filename+".w", c.filename+".r")
	}
}
func (c *CDR) RedoCDR() {
	c.f.Close()
	if c.Records == 0 && 0 == c.NullFlag {
		os.Remove(c.filename + ".w")
	} else {
		os.Rename(c.filename+".w", c.filename+".r")
	}

	cdrFile := c.GetCDRFileName()

	var err error
	c.f, err = os.OpenFile(cdrFile+".w", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		logg.Error("Open File %s failed, reason[%s].", cdrFile, err.Error())
	}
	logg.Debug("Close cdr file[%s.w], use new file [%s.w]", c.filename, cdrFile)
	c.filename = cdrFile
	c.id++
	c.Records = 0
}
func (c *CDR) CdrJob() {
	save := time.NewTicker(time.Duration(c.FileAliveTime) * time.Second)
	for {
		select {
		case <-save.C:
			c.RedoCDR()
		}
	}
}
func (c *CDR) GetCDRFileName() string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	timelayout := "20060102150405"
	cdrFile := fmt.Sprintf("%s/cdr/%s_%s_%06d", homeDir, c.CDRPrefix, time.Now().Format(timelayout), c.id)
	return cdrFile

}
