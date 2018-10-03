package lxf

import (
	"os"
	"testing"
)

type nameconv struct {
	In string
	R  string
	T  string
}

func TestParseImage(t *testing.T) {
	lxs, err := New("", os.Getenv("HOME")+"/.config/lxc/config.yml")
	if err != nil {
		t.Fatalf("could not create lx facade, %v", err)
	}

	cases := []nameconv{
		nameconv{
			In: "critest.asag.io/foobar:latest",
			R:  "critest.asag.io", T: "foobar"},
		nameconv{
			In: "critest.asag.io/hey/foo/bar",
			R:  "critest.asag.io", T: "hey/foo/bar"},
		nameconv{
			In: "foobar:latest",
			R:  "local", T: "foobar"},
		nameconv{
			In: "images/ubuntu",
			R:  "images", T: "ubuntu"},
		nameconv{
			In: "critest.asag.io/0aa4659583d2",
			R:  "critest.asag.io", T: "0aa4659583d2"},
		nameconv{
			In: "2dd442a946be6cee3d3a3cfb619a936fae0ff10600511596261d045960416f54",
			R:  "local", T: "2dd442a946be6cee3d3a3cfb619a936fae0ff10600511596261d045960416f54"},
	}

	for _, tc := range cases {
		t.Run(tc.In, func(t *testing.T) {
			id, err := lxs.parseImage(tc.In)
			if err != nil {
				t.Errorf("parse image fialed with error %v", err)
			}
			if id.Remote != tc.R {
				t.Errorf("for %v expected remote to be %v but is %v", tc.In, tc.R, id.Remote)
			}
			if id.Alias != tc.T {
				t.Errorf("for %v expected tag to be %v but is %v", tc.In, tc.T, id.Alias)
			}
		})
	}
}
