package logfile


/*
#cgo CFLAGS: -I../libxlsxwriter/include
#cgo LDFLAGS: -L../libxlsxwriter/lib -lxlsxwriter
#include "xlsxwriter.h"
*/
import "C"
import "unsafe"

type File struct {
	filename  *C.char
	workbook  *C.lxw_workbook
	worksheet *C.lxw_worksheet
}

func New(filename string) *File {
	filename =  "./log/" + filename + ".xls"
	file := &File{filename: C.strdup((*C.char)(unsafe.Pointer(&[]byte(filename)[0])))}
	file.workbook = C.workbook_new(file.filename)
	file.worksheet = C.workbook_add_worksheet(file.workbook, nil)
	return file
}

func (file *File) Write_data(generation int, fitness float32) {
	C.worksheet_write_number(file.worksheet, (C.lxw_row_t)(generation), 0, (C.double)(generation), nil)
	C.worksheet_write_number(file.worksheet, (C.lxw_row_t)(generation), 1, (C.double)(fitness), nil)
}

func (file *File) Close() {
	C.workbook_close(file.workbook)
}
