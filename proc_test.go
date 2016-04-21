package proc

import (
	"github.com/gogather/com/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProcess(t *testing.T) {
	Convey("Test Process Tree sections", t, func() {
		proc := GetProc(1)
		log.Greenln(proc)

		So(proc.GetName(), ShouldEqual, "init")
	})
}
