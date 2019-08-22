// 20 december 2015

package ui

// #include "pkgui.h"
import "C"

// MsgBoxError ..TODO
func MsgBoxError(w *Window, title string, description string) {
	ctitle := C.CString(title)
	defer freestr(ctitle)
	cdescription := C.CString(description)
	defer freestr(cdescription)
	C.uiMsgBoxError(w.w, ctitle, cdescription)
}

// OpenFile ..
func OpenFile(w *Window) string {
	cname := C.uiOpenFile(w.w)
	if cname == nil {
		return ""
	}
	defer C.uiFreeText(cname)
	return C.GoString(cname)
}

//SaveFile ..
func SaveFile(w *Window) string {
	cname := C.uiSaveFile(w.w)
	if cname == nil {
		return ""
	}
	defer C.uiFreeText(cname)
	return C.GoString(cname)
}

//MsgBox ..
func MsgBox(w *Window, title string, description string) {
	ctitle := C.CString(title)
	defer freestr(ctitle)
	cdescription := C.CString(description)
	defer freestr(cdescription)
	C.uiMsgBox(w.w, ctitle, cdescription)
}
