package video

import "testing"

func TestCodecs(t *testing.T) {
	if val := Codecs[""]; val != CodecUnknown {
		t.Fatal("default codec should be CodecUnknown")
	}

	if val := Codecs["avc"]; val != CodecAvc1 {
		t.Fatal("codec should be CodecAVC")
	}

	if val := Codecs["av1"]; val != CodecAv01 {
		t.Fatal("codec should be CodecAV1")
	}

	if val := Codecs["evc"]; val != CodecEvc1 {
		t.Fatal("codec should be CodecEVC")
	}

	if val := Codecs["vvcC"]; val != CodecVvc1 {
		t.Fatal("codec should be CodecVVC")
	}

	if val := Codecs["magicyuv"]; val != CodecMagicYUV {
		t.Fatal("codec should be CodecMagicYUV")
	}

	for _, fourcc := range []string{"m8rg", "m8ra", "m8rb", "m8y0", "m8y2", "m8y4", "m8ya", "m8g0"} {
		if val := Codecs[fourcc]; val != CodecMagicYUV {
			t.Fatalf("FourCC %q should map to CodecMagicYUV, got %q", fourcc, val)
		}
	}
}
