/* s3cp
 *
 * Uploads or downloads a file from an S3 bucket, using scp conventions.
 *
 * The intent of this program is to provide an scp-like method of working with files in S3.
 * At the outset, I intend to use this for daily backups, from a shell script, but I also
 * see that it could be used for other purposes.
 *
 * Where possible, I'll follow scp conventions. AWS credentials will be stored in the default
 * location for the go AWS implementation - ~/.aws/credentials.
 *
 * USAGE:
 *
 *     s3cp [--help] [--quiet] [--debug] [--region regionname] source destination
 *
 * WHERE:
 *
 *     --help      : prints help
 *
 *     --quiet     : suppress output
 * 
 *     --debug         : debug mode, prints lots of debug info
 *
 *     --region        : the AWS region of the target bucket; defaults to us-east-1
 *
 *     "source" and "desination" can be either local or remote objects, and both
 *     are required.
 *
 *     For a remote object: s3:bucket:/folder.../file.name
 *
 *     For a local object : /folder.../file.name
 *
 *     Local objects can also be STDIN or STDOUT, allowing you to
 *     redirect STDIN to a destination file at S3, or pipe a file
 *     as it is being downloaded from S3.
 *
 * EXAMPLES:
 *
 *     s3cp s3:/mybucket/myfolder/backup.tar.gz /tmp
 *       - downloads backup.tar.gz from S3 and places it in /tmp folder
 *
 *     s3cp s3:/mybucket/myfolder/backup.tar.gz /tmp/foobar.tar.gz
 *       - downloads backup.tar.gz from S3 and places it in the file /tmp/foobar.tar.gz
 *
 *     s3cp /tmp/backup.tar.gz s3:/mybucket/myfolder
 *       - uploads backup.tar.gz to S3 in the bucket mybucket and folder myfolder
 *
 */

package main

import (
	"fmt"
	"os"
	"regexp"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"flag"
)

//
// globals
//
var isQuietMode bool = false
var isDebugMode bool = false
var region string = "us-east-1"
var s3pathre = "^s3:([^:]+):(.+)$"

// Prints the help information to stdout
func printHelp() {
	// print help using the flag package default
	flag.Usage()
	// add the two trailing arguments for source and dest
	fmt.Fprintf(os.Stderr, "  source:  The source of the copy, either a local file path or an s3 path like s3:bucket:/path\n")
	fmt.Fprintf(os.Stderr, "  destination:  The destination of the copy, in the same format as source (above)\n")
	fmt.Fprintf(os.Stderr, "\nBoth source and destination are required, and one must be an s3 path, another must be a local path\n\n")
}

// Print to commandline if not in quiet mode
func printc(s string) {
	if !isQuietMode {
		fmt.Printf("%s", s)
	}
}

// Print debug information if in debug mode
func debug(s string) {
	if isDebugMode {
		fmt.Printf("DEBUG: %s", s)
	}
}

// Download a file from an S3 location, to a local file
func copyFromS3(source, destination string) error {
	// no-op
	return fmt.Errorf("Download is not currently implemented")
}

// Upload a file to an S3 location, from a local file
func copyToS3(s, d string) error {
	re := regexp.MustCompile(s3pathre)
	parts := re.FindStringSubmatch(d)
	bucket := parts[1]
	destpath := parts[2]

	if len(bucket) == 0 || len(destpath) == 0 {
		return fmt.Errorf("The destination path '%s' is invalid, must be in the form s3:bucket:/path/to/destination\n", d)
	}
	
	reader, err := os.Open(s)
	if err != nil {
		return fmt.Errorf("Failed to open file '%s', error was: %s\n", s, err)
	}

	uploader := s3manager.NewUploader(nil)
	if uploader == nil {
		return fmt.Errorf("Internal Error: Filure creating NewUploader @copyToS3()\n")
	}

	defer reader.Close()
	result, err2 := uploader.Upload(&s3manager.UploadInput{
		Body:     reader,
		Bucket:   aws.String(bucket),
		Key:      aws.String(destpath),
	})

	if err2 != nil {
		return fmt.Errorf("Failed to upload source file '%s' to destination '%s'\nError from S3 was: %s\n", s, d, err2)
	}

	debug(fmt.Sprintf("Post-upload file destionation URL:='%s'\n",result.Location))

	return nil
}

func parseCmdline() (source string, dest string) {
	var isHelpRequested bool = false

	flag.BoolVar(&isHelpRequested, "help", false, "(optional) Prints this help message")
	flag.BoolVar(&isDebugMode, "debug", false, "(optional) Used for debugging; outputs lots of debug info")
	flag.BoolVar(&isQuietMode, "quiet", false, "(optional) Suppresses output")
	flag.StringVar(&region, "region", "us-east-1", "(optional) The AWS region holding the target bucket; defaults to 'us-east-1'")
	flag.Parse()

	debug(fmt.Sprintf("got args '%s'\n", os.Args))

	// tried to use flag.ErrHelp here, but it didn't work as expected...seems that ErrHelp is
	// always being set to the "help requested" value. I wonder if the flag package assumes
	// all flags are required or something
	if flag.NArg() != 2 || isHelpRequested {
		debug(fmt.Sprintf("NArg = %d, ErrHelp = %s\n", flag.NArg(), flag.ErrHelp))
		printHelp()
		os.Exit(1)
	}

	return flag.Arg(0), flag.Arg(1) 
}

func main() {
	source, destination := parseCmdline()

	aws.DefaultConfig.Region = region

	// one of source or destination must be an s3 path
	isMatchSource, _ := regexp.MatchString(s3pathre,source)
	isMatchDestination, _ := regexp.MatchString(s3pathre,destination)
	debug(fmt.Sprintf("Source match = %t, dest match = %t\n", isMatchSource, isMatchDestination))

	var err error
	if isMatchSource {
		debug("Calling copyFromS3")
		err = copyFromS3(source,destination)
	} else if isMatchDestination {
		debug("Calling copyToS3")
		err = copyToS3(source,destination)
	} else {
		printHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("%s\n",err.Error())
		os.Exit(1)
	}

	printc(fmt.Sprintf("%s -> %s : transfer complete\n",source,destination))
	os.Exit(0)
}
