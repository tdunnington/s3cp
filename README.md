#s3cp
An scp-style command-line tool for uploading and downloading files to and from an AWS S3 bucket.

### TOC
* [License](README.md#license)
* [Project Status](README.md#current-status)
* [How to use s3cp](README.md#how-to-use-s3cp)
** [Install s3cp](README.md#step-1-install-s3cp)

### License
Released under the Apache 2.0 license

### Current Status
s3cp is currently in beta verion 0.2
* Uploading to an S3 bucket
* Downloading from an S3 bucket
* Errors are trapped and reported properly
* Exit codes work as expected (0 for success, non-0 for failure)

# How to use s3cp

### Step 1 - Install s3cp
You can download the executable from the homepage, or you can get the sources and compile. The
executable is statically linked (as all Go programs are), so you can download the executable
and run it without installing Go.

If you intend to compile, you'll need to install the Go language
(for more information about installing Go on your platform, visit the [Go homepage](http://golang.org))

If you have Go installed, you can get s3cp the "Go way" like this:
```
go get -u http://github.com/aws/aws-sdk-go/...
go get http://github.com/tdunnington/s3cp
```

### Step 2 - Setup your AWS S3 Bucket
Setting up an S3 bucket to work with s3cp is really about security. The policy editors have become much
easier to use over the last year.

Please note that these are my recommended settings. You are ultimately responsible for the security of your bucket!

**Setting up with a user account**

This is the only method you can use if your server is outside of AWS. If your server is inside AWS (hosted via
EC2), then the "Setting up with EC2 IAM roles" is a better choice for you, but if you did not create an IAM
role when you created the server, you will either have to use the user account method, or re-launch your server
with an IAM role.

In this method, you will create an IAM user and an IAM group in AWS. The group will contain the policy permissions
that will allow you to upload and download to the target bucket. The advantage of this method is that you can easily
swap out credentials, but the downside is that the credentials have to be stored on the system, where they are
vulnerable.

Create a new group in IAM (let's call it "s3cp_test"). Do not create a policy for the group when you create it; 
we'll add the policy in the next step:

...TODO add images...

Now add an inline policy to the group:

...TODO add images...

Use the "Custom Policy" option and paste this sample policy into the editor. Note that you have to change the 
bucket name:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "removed for security",
            "Effect": "Allow",
            "Action": [
                "s3:DeleteObject",
                "s3:DeleteObjectVersion",
                "s3:GetBucketPolicy",
                "s3:GetObject",
                "s3:GetObjectAcl",
                "s3:GetObjectVersion",
                "s3:GetObjectVersionAcl",
                "s3:ListBucket",
                "s3:ListBucketMultipartUploads",
                "s3:ListBucketVersions",
                "s3:ListMultipartUploadParts",
                "s3:PutObject",
                "s3:PutObjectAcl",
                "s3:PutObjectVersionAcl"
            ],
            "Resource": [
                "arn:aws:s3:::<PUT BUCKET NAME HERE>/*"
            ]
        }
    ]
}
```

**Note**: I am not certain that these are the minimal permissions necessary. I have not found documentation that 
makes it clear what permissions are required to add an item to a bucket using the [s3manager](http://github.com/aws/aws-sdk-go/service/s3/s3manager) that is used by s3cp.

Now, create a user for s3cp, and save the AWS security keys. These screenshots show an actual AWS key pair;
rest assured that key pair has been deleted and is not valid.

...TODO add screenshots...

Hold on to the key pairs; you will need them for the next step. Go on to Step 3 now.

**Setting up with EC2 IAM roles**
This method allows you to use AWS resources from an EC2 instance without having to store the key pair on the
server anywhere; AWS will know the server has permissions based on the server role.

In order for this method to work, the IAM role must be applied to the server when it is created.

Create a role for using s3cp, and add a custom inline policy, exactly like the policy shown in above:

...TODO add screenshots...

Then launch your server, and make sure to select the new role when you launch:

...TODO add screenshots...

### Step 3 - Create your credentials file/environment variables
If you are using a security role, you can skip this step and start using s3cp on your server.

In the previous step, when you created the IAM group, you saved the security keys. You can inform s3cp of these
keys in one of two ways:

**Using a credentials file**

Using this method, you can setup separate credentials for each user account.

Create a file ~/.aws/credentials, with the following contents (you can use a profile if you set the environment
variable `AWS_PROFILE`):
```
[default]
aws_access_key_id = <put your access key id here>
aws_secret_access_key = <put the secret key here>
```

If you are going to run a cronjob or other script from various user accounts, make sure that the user in question has a credentials file. So you may have to create multiple credentials files, or use the command line approach below.

I also recommend:

```
chmod 700 ~/.aws
chmod 400 ~/.aws/credentials
```

For more information about credentials, check out the [Getting Started Credentials](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials) guide.

**Using environment variables**

Using this method, you create two environment variables:

```
export AWS_ACCESS_KEY_ID = <put your access key id here>
export AWS_SECRET_ACCESS_KEY = <put the secret key here>
```

Of course, these variables have to be set before you execute s3cp.

### Step 4 - Use s3cp

Using s3cp was meant to be similar to scp. Here's the big difference:

In scp, you would do something like this to upload a file:
```
scp /some/local/file.abc user@server:/destination/path
```

In s3cp, the "remote" format is a bit different:
```
s3cp /some/local/file.abc s3:bucket:/destination/path/file.abc
```

You must provide a full path and filename when working with S3; you cannot take the "shortcut"
that you can with cp or scp.

Otherwise, you have some commandline options:
```
  -debug: (optional) Used for debugging; outputs lots of debug info; default is false
  -help: (optional) Prints this help message; default is false
  -quiet: (optional) Suppresses output, useful for scripting; default is false
  -region: (optional) The AWS region holding the target bucket; defaults to 'us-east-1'
```

**Note about region**: The region name will display in the AWS console differently than the region id 
s3cp expects. To find your region id, login to the AWS console, open your bucket, and you'll see the
region id in the URL, like so: `https://console.aws.amazon.com/s3/home?region=us-east-1#&bucket=...`

### About this project
I started this project because I wanted to backup my websites to s3, but getting the file up to S3 isn't
trivial. Yes there were some publicly available options, including python scripts, and I suppose if I 
had looked long enough I might have found somethig like s3cp. But I didn't/couldn't and just decided to
do my own.

Go was a good choice of language because it compiles to a static binary, which should make it easier for
people to install and use. Go should also be cross-platform, so I'm hopeful that this will work on Windows
and OSX as well as it works on Linux.

As I was working on this, I realized I could do a lot of other tools, and so I intend to extend this to an
s3tools project at some point. In additiona to s3cp, I would like to do s3rm, s3mkdir, s3ls 
and more. With the complete toolkit, you should be able to script some cool cmdline jobs.

### Upcoming Features
For s3cp, still left to do are:
* Support for downloading
* Redirect STDIN and STDOUT, or piping
* On-the-fly compression
* Not sure how the `go get` dependencies work for my package, need to check that out
