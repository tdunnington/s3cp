#s3cp
An scp-style command-line tool for uploading and downloading files to and from an AWS S3 bucket.

### License
Released under the Apache 2.0 license

### Current Status
s3cp is currently in beta verion 0.1
* Uploading to an S3 bucket works
* Errors are trapped and reported properly
* Exit codes work as expected (0 for success, non-0 for failure)
* *Downloads are not working yet*

# How to use s3cp

### Step 1 - Install s3cp
You can download the executable from the homepage, or you can get the sources and compile. The
executable is statically linked (as all Go programs are), so you can download the executable
and run it without installing Go.

If you intend to compile, you'll need to install the Go language
(for more information about installing Go on your platform, visit the [Go homepage](http://golang.org))

If you have Go installed, you can get s3cp the "Go way" like this:
  go install http://github.com/tdunnington/s3cp
  
### Step 2 - Setup your AWS S3 Bucket
TODO

### Step 3 - Create your credentials file/environment variables
TODO

### Step 4 - Use s3cp
TODO
