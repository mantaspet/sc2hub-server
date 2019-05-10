package main

//func TestPing(t *testing.T) {
//	app := newTestApplication(t)
//	ts := newTestServer(t, app.router())
//	defer ts.Close()
//
//	code, _, body := ts.get(t, "/ping")
//
//	if code != http.StatusOK {
//		t.Errorf("want %d; got %d", http.StatusOK, code)
//	}
//
//	if string(body) != "OK" {
//		t.Errorf("want body to equal %q", "OK")
//	}
//}
