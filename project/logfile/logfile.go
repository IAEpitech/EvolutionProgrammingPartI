package logfile

/*
#cgo CFLAGS: -I../libxlsxwriter/include
#cgo LDFLAGS: -L../libxlsxwriter/lib -lxlsxwriter
#include "xlsxwriter.h"
*/
import "C"
import "unsafe"

type FILE struct {
	filename  *C.char
	workbook  *C.lxw_workbook
	worksheet *C.lxw_worksheet
}

var file *FILE

func init() {
	file = new(FILE)
	name := "./log/log.xls"
	file.filename = C.strdup((*C.char)(unsafe.Pointer(&[]byte(name)[0])))
	file.workbook = C.workbook_new(file.filename)
	file.worksheet = C.workbook_add_worksheet(file.workbook, nil)
}

func Write_data(generation int, fitness float32) {
	C.worksheet_write_number(file.worksheet, (C.lxw_row_t)(generation), 0, (C.double)(generation), nil)
	C.worksheet_write_number(file.worksheet, (C.lxw_row_t)(generation), 1, (C.double)(fitness), nil)
}

func End() {
	C.workbook_close(file.workbook)
}
