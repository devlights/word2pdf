package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/devlights/gord"
	"github.com/devlights/gord/constants"
)

type target struct {
	filePath string
	absPath  string
	pdfPath  string
	verbose  bool
}

func (me *target) abs() string {
	if me.absPath == "" {
		v, _ := filepath.Abs(me.filePath)
		me.absPath = v
	}

	if me.verbose {
		slog.Info("abs", "path", me.absPath)
	}

	return me.absPath
}

func (me *target) convert() string {
	if me.absPath == "" {
		me.abs()
	}

	me.pdfPath = me.absPath[:strings.Index(me.absPath, filepath.Ext(me.absPath))] + ".pdf"

	if me.verbose {
		slog.Info("convert", "abs", me.absPath, "pdf", me.pdfPath)
	}

	return me.pdfPath
}

func main() {
	var (
		verbose bool
	)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: word2pdf.exe [-v] word-file-path")
		flag.PrintDefaults()
	}

	flag.BoolVar(&verbose, "v", false, "verbose log output")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	if err := run(&target{filePath: flag.Arg(0), verbose: verbose}); err != nil {
		slog.Error(err.Error())
	}
}

func run(p *target) error {
	if p.verbose {
		slog.Info("start")
		defer slog.Info("done")
	}

	quitFn := gord.MustInitGord()
	defer quitFn()

	g, gordReleaseFn := gord.MustNewGord()
	defer gordReleaseFn()

	if err := g.SetVisible(false); err != nil {
		return err
	}

	docs, err := g.Documents()
	if err != nil {
		return err
	}

	doc, docReleaseFn, err := docs.Open(p.abs())
	if err != nil {
		return err
	}
	defer docReleaseFn()

	if p.verbose {
		slog.Info("Export start")
	}

	if err := doc.ExportAsFixedFormat(p.convert(), constants.WdExportFormatPDF); err != nil {
		return err
	}

	if p.verbose {
		slog.Info("Export end")
	}

	return nil
}
