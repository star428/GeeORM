package sqllog

import (
	"os"
	"testing"
)

func TestSetLevel(t *testing.T) {
	SetLevel(InfoLevel)
	// info级别的日志两个都可以输出
	if infoLog.Writer() != os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set infolevel")
	}

	// error级别的日志只有errorLog可以输出
	SetLevel(ErrorLevel)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set errorlevel")
	}

	// 禁用日志
	SetLevel(Disabled)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() == os.Stdout {
		t.Fatal("failed to set disablelevel")
	}
}
