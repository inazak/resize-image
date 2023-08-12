package main

import (
  "flag"
	"fmt"
  "os"
  "path/filepath"
  "strings"
	"github.com/disintegration/imaging"
)

var usage =`
Usage:

resize-image [OPTIONS] FILE_PATTERN

  FILE_PATTERN is shell file name syntax like '*.jpg'.

  OPTIONS:
  -w | -width=1024             resize image to the specified width,
                               aspect ratio is preserved. 
  -f | -force                  forced overwriting with resized image.
  -p | -postfix=".resized"     sub extension of converted file.
  -r | rtype=0-3               resampling filter, default is 0.
       0 ... Lanczos - high-quality, sharp results
       1 ... CatmullRom - sharp cubic filter, faster than Lanczos
       2 ... Linear - bilinear resampling, faster than cubic filters
       3 ... NearestNeighbor - fastest, no antialiasing
`

var optionsWidth   int
var optionsForce   bool
var optionsFilter  int
var optionsPostfix string

var resampleFilterMap = map[int]imaging.ResampleFilter {
  0: imaging.Lanczos,
  1: imaging.CatmullRom,
  2: imaging.Linear,
  3: imaging.NearestNeighbor,
}


func main() {

  flag.IntVar(&optionsWidth,  "width", 1024, "width")
  flag.IntVar(&optionsWidth,  "w",     1024, "width")
  flag.BoolVar(&optionsForce, "force", false, "forced overwriting with resized image")
  flag.BoolVar(&optionsForce, "f",     false, "forced overwriting with resized image")
  flag.IntVar(&optionsFilter, "rtype", 0, "resampling filter type(0-3)")
  flag.IntVar(&optionsFilter, "r",     0, "resampling filter type(0-3)")
  flag.StringVar(&optionsPostfix, "postfix", ".resized", "sub extension of converted file")
  flag.StringVar(&optionsPostfix, "p",       ".resized", "sub extension of converted file")
  flag.Parse()

  if len(flag.Args()) != 1 {
    fmt.Printf("%s", usage)
    os.Exit(1)
  }

  filename := flag.Args()[0]
  width    := optionsWidth
  filter   := optionsFilter
  postfix  := optionsPostfix

  if (filter < 0) || (3 < filter) {
    fmt.Printf("%s", usage)
    os.Exit(1)
  }

  list, err := getFilepathList(filename)
  if err != nil {
    fmt.Printf("%s", err.Error())
    os.Exit(2)
  }

  filecount := len(list)
  errcount  := 0

  for i, file := range list {

    fmt.Printf("%d/%d working ... %s\n", i+1, filecount, file)

    src, err := imaging.Open(file)
    if err != nil {
      errcount += 1
      fmt.Printf("[ERROR] opening %s\n", file)
      continue
    }

    dst := imaging.Resize(src, 0, width, resampleFilterMap[filter])

    dir, name, ext := splitFilepath(file)
    newfile := filepath.Join(dir, name + postfix + ext)

    err = imaging.Save(dst, newfile)
    if err != nil {
      errcount += 1
      fmt.Printf("[ERROR] saving %s\n", newfile)
      continue
    }

    if optionsForce {
      err = os.Rename(newfile, file)
      if err != nil {
        errcount += 1
        fmt.Printf("[ERROR] rename from %s to %s\n", newfile, file)
        continue
      }
    }

  }

  fmt.Printf("done ... error count is %d\n", errcount)
  os.Exit(0)
}


func getFilepathList(filter string) (list []string, err error) {

  entries, err := filepath.Glob(filter)
  if err != nil {
    return nil, err
  }

  for _, entry := range entries {
    isdir, err := isDir(entry)
    if err != nil {
      return nil, err
    }
    if ! isdir {
      list = append(list, entry)
    }
  }

  return list, nil
}


func isDir(p string) (bool, error) {
  f, err := os.Stat(p)
  if err != nil {
    return false, err
  }
  return f.Mode().IsDir(), nil
}


func splitFilepath(fpath string) (dir, name, ext string) {

  dir   = filepath.Dir(fpath)
  base := filepath.Base(fpath)

  dotindex := strings.LastIndex(base, ".")

  // file has no extention
  if dotindex == -1 {
    name = base
    ext  = ""
    return
  }

  name = base[0:dotindex]
  ext  = base[dotindex:len(base)]
  return
}

