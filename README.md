# resize-image

image resize utility for windows command prompt

## usage

```
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
```

