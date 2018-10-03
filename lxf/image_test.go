package lxf_test

import (
	"testing"
	"time"
)

func TestListImages(t *testing.T) {
	lt := newLXFTest(t)
	imgs := lt.listImages("")
	if len(imgs) != 1 {
		t.Errorf("expected one local image but there are %v", len(imgs))
	}
}

func TestGetMissingImageRemote(t *testing.T) {
	lt := newLXFTest(t)
	img := lt.getImage("lala.com/foobar")
	if img != nil {
		t.Errorf("expected missing image to be nil but is %v", img)
	}
}
func TestGetMissingImage(t *testing.T) {
	lt := newLXFTest(t)
	img := lt.getImage("critest.asag.io/foobar/lala")
	if img != nil {
		t.Errorf("expected missing image to be nil but is %v", img)
	}
}

func TestGetImage(t *testing.T) {
	lt := newLXFTest(t)
	hash := lt.pullImage("critest.asag.io/cri-tools/test-image-1")
	img := lt.getImage("critest.asag.io/cri-tools/test-image-1")
	if img == nil {
		t.Errorf("expected to find image but it's nil")
	}
	if hash != img.Hash {
		t.Errorf("expected image to have same hash but they are different")
	}
}

func TestGetImageByHash(t *testing.T) {
	lt := newLXFTest(t)
	hash := lt.pullImage("critest.asag.io/cri-tools/test-image-1")
	img := lt.getImage(hash)
	if img == nil {
		t.Errorf("expected to find image by hash but it's nil")
	}
}

func TestPullImage(t *testing.T) {
	lt := newLXFTest(t)
	lt.pullImage("critest.asag.io/cri-tools/test-image-1")
	imgs := lt.listImages("")
	if len(imgs) != 2 {
		t.Errorf("expected two local images but there are %v", len(imgs))
	}
}

func TestFilterImage(t *testing.T) {
	lt := newLXFTest(t)
	lt.pullImage("critest.asag.io/cri-tools/test-image-1")
	imgs := lt.listImages("")
	imgs = lt.listImages(imgs[1].Hash)
	if len(imgs) != 1 {
		t.Errorf("expected to get only one image but got %v", len(imgs))
	}
}

func TestRemoveImage(t *testing.T) {
	lt := newLXFTest(t)
	lt.pullImage("critest.asag.io/cri-tools/test-image-1")
	lt.removeImage("critest.asag.io/cri-tools/test-image-1")

	imgs := lt.listImages("")
	if len(imgs) != 1 {
		t.Errorf("expected one local image but there are %v", len(imgs))
	}
}

func TestRemoveImageByHASH(t *testing.T) {
	lt := newLXFTest(t)
	hash := lt.pullImage("critest.asag.io/cri-tools/test-image-1")
	lt.removeImage(hash)

	imgs := lt.listImages("")
	if len(imgs) != 1 {
		t.Errorf("expected one local image but there are %v", len(imgs))
	}
}

func TestPullBigImageCache(t *testing.T) {
	lt := newLXFTest(t)
	lt.pullImage("images/ubuntu/18.04")
	start := time.Now()
	lt.pullImage("images/ubuntu/18.04")
	duration := time.Since(start)
	if duration > time.Second {
		t.Errorf("took to long (%v) to download big image which is already on target server",
			duration)
	}
}
