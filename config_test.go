package balaur

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testTomlConfigPath string = "testfixtures/testconfig.toml"

func TestNewConfig(t *testing.T) {
	Convey("Test create new config", t, func() {
		config := NewConfig(testTomlConfigPath)
		_, ok := config.(*TomlConfig)
		So(ok, ShouldEqual, true)
	})
}

func TestTomlConfig(t *testing.T) {
	Convey("Test TOML config", t, func() {
		config := NewTomlConfig(testTomlConfigPath)
		Convey("Test Get()", func() {
			So(config.Get("parent", false), ShouldEqual, "/some/parent/*")
			So(config.Get("wrongkey", false), ShouldEqual, "")
		})

		Convey("Test GetArray()", func() {
			So(len(config.GetArray("arr", false)), ShouldEqual, 3)
			So(config.GetArray("wrongkey", false), ShouldEqual, nil)
		})

		Convey("Test GetChildren()", func() {
			So(len(config.GetChildren("route", false)), ShouldEqual, 2)
			So(config.GetChildren("wrongkey", false), ShouldEqual, nil)
		})
	})
}
